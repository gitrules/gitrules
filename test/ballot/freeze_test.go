package ballot

import (
	"fmt"
	"testing"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/ballot/ballotapi"
	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/member"
	"github.com/gitrules/gitrules/proto/purpose"
	"github.com/gitrules/gitrules/runtime"
	"github.com/gitrules/gitrules/test"
)

func TestVoteFreezeVote(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ballotproto.ParseBallotID("a/b/c")
	choices := []string{"x", "y", "z"}

	// open
	strat := ballotio.QVPolicyName
	openChg := ballotapi.Open(
		ctx,
		strat,
		cty.Organizer(),
		ballotName,
		account.NobodyAccountID,
		purpose.Unspecified,
		"",
		"ballot title",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", form.SprintJSON(openChg))

	// give voter credits
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 1.0), "test")

	// vote
	elections := ballotproto.Elections{
		ballotproto.NewElection(choices[0], 1.0),
	}
	voteChg := ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)
	fmt.Println("vote: ", form.SprintJSON(voteChg))

	// freeze ballot
	freezeChg := ballotapi.Freeze(ctx, cty.Organizer(), ballotName)
	fmt.Println("freeze: ", form.SprintJSON(freezeChg))

	// try voting while frozen
	if must.Try(
		func() { ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections) },
	) == nil {
		t.Fatalf("voting on a frozen ballot should have failed")
	}

	// unfreeze ballot
	unfreezeChg := ballotapi.Unfreeze(ctx, cty.Organizer(), ballotName)
	fmt.Println("unfreeze: ", form.SprintJSON(unfreezeChg))

	// tally
	tallyChg := ballotapi.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	closeChg := ballotapi.Close(ctx, cty.Organizer(), ballotName, account.BurnAccountID)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// testutil.Hang()
}

// TestVoteFreezeTally tests that votes made during a freeze are consumed and discarded by tallying.
func TestVoteFreezeTally(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ballotproto.ParseBallotID("a/b/c")
	choices := []string{"x", "y", "z"}

	// open
	strat := ballotio.QVPolicyName
	openChg := ballotapi.Open(ctx, strat, cty.Organizer(), ballotName, account.NobodyAccountID, purpose.Unspecified, "", "ballot title", "ballot description", choices, member.Everybody)
	fmt.Println("open: ", form.SprintJSON(openChg))

	// give voter credits
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 1.0), "test")

	// vote
	elections := ballotproto.Elections{
		ballotproto.NewElection(choices[0], 1.0),
	}
	voteChg := ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)
	fmt.Println("vote: ", form.SprintJSON(voteChg))

	// freeze
	freezeChg := ballotapi.Freeze(ctx, cty.Organizer(), ballotName)
	fmt.Println("freeze: ", form.SprintJSON(freezeChg))

	// verify state changed
	ast := ballotapi.Show(ctx, gov.Address(cty.Organizer().Public), ballotName)
	if !ast.Ad.Frozen {
		t.Errorf("expecting frozen")
	}

	// tally
	tallyChg := ballotapi.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 0.0 {
		t.Errorf("expecting %v, got %v", 0.0, tallyChg.Result.Scores[choices[0]])
	}

	// testutil.Hang()
}
