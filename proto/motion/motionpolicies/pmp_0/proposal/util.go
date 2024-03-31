package proposal

import (
	"context"

	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/ns"
	"github.com/gitrules/gitrules/proto/motion/motionpolicies/pmp_0"
)

func SaveState_StageOnly(ctx context.Context, t *git.Tree, policyNS ns.NS, state *pmp_0.ProposalState) {
	git.ToFileStage[*pmp_0.ProposalState](ctx, t, policyNS.Append(pmp_0.StateFilebase), state)
}

func LoadState_Local(ctx context.Context, t *git.Tree, policyNS ns.NS) *pmp_0.ProposalState {
	state := git.FromFile[pmp_0.ProposalState](ctx, t, policyNS.Append(pmp_0.StateFilebase))
	return &state
}
