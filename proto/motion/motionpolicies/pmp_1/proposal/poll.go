package proposal

import (
	"context"

	"github.com/gitrules/gitrules/proto/ballot/ballotapi"
	"github.com/gitrules/gitrules/proto/ballot/ballotio"
	"github.com/gitrules/gitrules/proto/ballot/ballotpolicies/sv"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

func init() {
	ctx := context.Background()
	ballotio.Install(
		ctx,
		ProposalApprovalPollPolicyName,
		sv.SV{
			Kernel: ScoreKernel{},
		},
	)
}

const ProposalApprovalPollPolicyName ballotproto.PolicyName = "pmp-proposal-approval-v1"

type ScoreKernel struct{}

type ScoreKernelState struct {
	MotionID              motionproto.MotionID `json:"motion_id"`
	InverseCostMultiplier float64              `json:"inverse_cost_multiplier"`
	Bounty                float64              `json:"bounty"`
}

func (sk ScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	el ballotproto.AcceptedElections,

) sv.ScoredVotes {

	state := ballotapi.LoadPolicyState_Local[ScoreKernelState](ctx, cloned, ad.ID)
	qvSK := sv.MakeQVScoreKernel(ctx, state.InverseCostMultiplier)
	return qvSK.Score(ctx, cloned, ad, el)
}

func (sk ScoreKernel) CalcJS(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,

) *ballotproto.Margin {

	state := ballotapi.LoadPolicyState_Local[ScoreKernelState](ctx, cloned, ad.ID)
	qvSK := sv.MakeQVScoreKernel(ctx, state.InverseCostMultiplier)
	margin := qvSK.CalcJS(ctx, cloned, ad, tally)
	margin.Reward = &ballotproto.MarginCalculator{
		Label:       "Reward",
		Description: "Potential reward for the voter, assuming the vote is aligned with the PR outcome",
		FnJS:        rewardJSFmt,
	}
	return margin
}

const (
	rewardJSFmt = `
	function(voteUser, voteChoice, voteImpact) {
		return 2*Math.abs(voteImpact);
	}
	`
)
