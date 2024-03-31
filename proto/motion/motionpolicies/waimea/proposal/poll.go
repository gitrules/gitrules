package proposal

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
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

const ProposalApprovalPollPolicyName ballotproto.PolicyName = "waimea-proposal-approval"

type ApprovalPollState struct {
	MotionID              motionproto.MotionID `json:"motion_id"`
	InverseCostMultiplier float64              `json:"inverse_cost_multiplier"`
	Bounty                float64              `json:"bounty"`
}

type ScoreKernel struct{}

func (sk ScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	el ballotproto.AcceptedElections,

) sv.ScoredVotes {

	state := ballotapi.LoadPolicyState_Local[ApprovalPollState](ctx, cloned, ad.ID)
	qvSK := sv.MakeQVScoreKernel(ctx, state.InverseCostMultiplier)
	return qvSK.Score(ctx, cloned, ad, el)
}

func (sk ScoreKernel) CalcJS(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,

) *ballotproto.Margin {

	state := ballotapi.LoadPolicyState_Local[ApprovalPollState](ctx, cloned, ad.ID)
	qvSK := sv.MakeQVScoreKernel(ctx, state.InverseCostMultiplier)
	margin := qvSK.CalcJS(ctx, cloned, ad, tally)
	margin.Reward = &ballotproto.MarginCalculator{
		Label:       "Reward",
		Description: "Potential reward to voter, in the event their vote is aligned with the outcome",
		FnJS:        fmt.Sprintf(rewardJSFmt, form.SprintJSON(tally)),
	}
	return margin
}

/*
	For testing:
		let tally = {
			"scores_by_user": {
				"p1": {"rank": { "score": 1, "strength": 1}},
				"p2": {"rank": { "score": 2, "strength": 4}},
				"n1": {"rank": { "score": -1, "strength": -1}},
				"n2": {"rank": { "score": -2, "strength": -4}},
			},
		};
*/

const (
	rewardJSFmt = `
	function(voteUser, voteChoice, voteImpact) {

		if (voteImpact === 0) {
			return 0;
		}
		let voterSign = Math.sign(voteImpact);

		let tally = %s;
		var scoresByUser = tally.scores_by_user;
		if (scoresByUser === undefined) {
			return voteImpact*voteImpact;
		}

		let voterWinnerShare = Math.abs(voteImpact);
		var totalWinnerShares = Math.abs(voteImpact);
		var totalLoserCost = 0.0;
		for (const user in scoresByUser) {
			if (user === voteUser) {
				continue;
			}
			let scoresByChoice = scoresByUser[user];
			var ss = scoresByChoice[voteChoice];
			var ss = scoresByChoice[voteChoice];
			if (ss != undefined) {
				if (Math.sign(ss.score) === voterSign) {
					totalWinnerShares += Math.abs(ss.score);
				} else {
					totalLoserCost += Math.abs(ss.strength);
				}
			}
		}

		if (totalWinnerShares > 0) {
			return voteImpact*voteImpact+totalLoserCost*(voterWinnerShare/totalWinnerShares);
		} else {
			return voteImpact*voteImpact;
		}
	}
	`
)
