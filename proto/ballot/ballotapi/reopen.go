package ballotapi

import (
	"context"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
)

func Reopen(
	ctx context.Context,
	addr gov.OwnerAddress,
	id ballotproto.BallotID,

) git.Change[form.Map, form.None] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := Reopen_StageOnly(ctx, cloned, id)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Reopen_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,

) git.Change[form.Map, form.None] {

	t := cloned.Public.Tree()

	// verify ad and policy are present
	ad, policy := ballotio.LoadAdPolicy_Local(ctx, t, id)
	must.Assertf(ctx, ad.Closed, "ballot is not closed")
	must.Assertf(ctx, !ad.Cancelled, "ballot was cancelled")

	tally := loadTally_Local(ctx, t, id)
	chg := policy.Reopen(ctx, cloned, &ad, &tally)

	// remove prior outcome
	_, err := git.TreeRemove(ctx, t, id.OutcomeNS())
	must.NoError(ctx, err)

	// write state
	ad.Closed = false
	ad.Cancelled = false
	git.ToFileStage(ctx, t, id.AdNS(), ad)

	return chg
}
