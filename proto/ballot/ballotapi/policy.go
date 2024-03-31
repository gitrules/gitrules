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

func LoadPolicyState[PS form.Form](
	ctx context.Context,
	addr gov.Address,
	id ballotproto.BallotID,

) PS {

	cloned := gov.Clone(ctx, addr)
	return LoadPolicyState_Local[PS](ctx, cloned, id)
}

func LoadPolicyState_Local[PS form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,

) PS {

	t := cloned.Tree()
	return git.FromFile[PS](ctx, t, id.PolicyNS())
}

func SavePolicyState[PS form.Form](
	ctx context.Context,
	addr gov.Address,
	id ballotproto.BallotID,
	policyState PS,

) {

	cloned := gov.Clone(ctx, addr)
	SavePolicyState_StageOnly[PS](ctx, cloned, id, policyState)
	proto.Commitf(ctx, cloned, "ballot_save_policy_state", "update ballot policy state")
	cloned.Push(ctx)
}

func SavePolicyState_StageOnly[PS form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,
	policyState PS,

) {

	t := cloned.Tree()
	ad := ballotio.LoadAd_Local(ctx, t, id)
	must.Assertf(ctx, !ad.Closed, "ballot already closed")
	git.ToFileStage[PS](ctx, t, id.PolicyNS(), policyState)
}
