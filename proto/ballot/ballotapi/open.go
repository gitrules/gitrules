package ballotapi

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/history/trace"
	"github.com/gitrules/gitrules/proto/member"
	"github.com/gitrules/gitrules/proto/motion"
	"github.com/gitrules/gitrules/proto/purpose"
)

func Open(
	ctx context.Context,
	strat ballotproto.PolicyName,
	addr gov.OwnerAddress,
	id ballotproto.BallotID,
	owner account.AccountID,
	purpose purpose.Purpose,
	motionPolicy motion.PolicyName,
	title string,
	description string,
	choices []string,
	participants member.Group,

) git.Change[form.Map, ballotproto.BallotAddress] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := Open_StageOnly(
		ctx,
		strat,
		cloned,
		id,
		owner,
		purpose,
		motionPolicy,
		title,
		description,
		choices,
		participants,
	)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Open_StageOnly(
	ctx context.Context,
	policyName ballotproto.PolicyName,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	owner account.AccountID,
	purpose purpose.Purpose,
	motionPolicy motion.PolicyName,
	title string,
	description string,
	choices []string,
	participants member.Group,

) git.Change[form.Map, ballotproto.BallotAddress] {

	// check no open ballots by the same name
	if ballotproto.BallotKV.Contains(ctx, ballotproto.BallotNS, cloned.PublicClone().Tree(), id) {
		must.Errorf(ctx, "ballot already exists: %v", id.AdNS().GitPath())
	}

	// verify group exists
	if !member.IsGroup_Local(ctx, cloned.PublicClone(), participants) {
		must.Errorf(ctx, "participant group %v does not exist", participants)
	}

	// create escrow account
	account.Create_StageOnly(
		ctx, cloned.PublicClone(),
		ballotproto.BallotEscrowAccountID(id),
		account.NobodyAccountID,
		fmt.Sprintf("opening ballot %v", id),
	)

	// create ballot
	ballotproto.BallotKV.Set(ctx, ballotproto.BallotNS, cloned.PublicClone().Tree(), id, struct{}{})

	// write ad
	ad := ballotproto.Ad{
		Gov:          cloned.GovAddress(),
		ID:           id,
		Owner:        owner,
		Purpose:      purpose,
		MotionPolicy: motionPolicy,
		//
		Title:       title,
		Description: description,
		//
		Choices:      choices,
		Policy:       policyName,
		Participants: participants,
		//
		Frozen:    false,
		Closed:    false,
		Cancelled: false,
		//
		ParentCommit: git.Head(ctx, cloned.Public.Repo()),
	}
	git.ToFileStage(ctx, cloned.Public.Tree(), id.AdNS(), ad)

	// initialize tally
	policy := ballotio.LookupPolicy(ctx, policyName)
	tally := policy.Open(ctx, cloned, &ad)
	git.ToFileStage(ctx, cloned.Public.Tree(), id.TallyNS(), tally)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_open",
		Args:   trace.M{"id": id},
		Result: trace.M{"ad": ad},
	})

	return git.NewChange(
		fmt.Sprintf("Create ballot of type %v", policyName),
		"ballot_open",
		form.Map{
			"policy":       policyName,
			"id":           id,
			"participants": participants,
		},
		ballotproto.BallotAddress{Gov: cloned.GovAddress(), Name: id},
		nil,
	)
}
