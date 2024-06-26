package proposal

import (
	"context"
	"fmt"
	"math"

	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/ballot/ballotapi"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/history/metric"
	"github.com/gitrules/gitrules/proto/member"
	"github.com/gitrules/gitrules/proto/motion/motionapi"
	"github.com/gitrules/gitrules/proto/motion/motionpolicies/waimea"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

func loadResolvedConcerns(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,

) (resolved motionproto.Motions, projectedBounties []float64, totalProjectedBounty float64) {

	eligible := computeEligibleConcerns(ctx, cloned.PublicClone(), prop)
	for _, ref := range eligible {
		con := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
		conState := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID)
		//
		resolved = append(resolved, con)
		projectedBounties = append(projectedBounties, conState.ProjectedPriorityBounty())
	}

	totalProjectedBounty = 0.0
	for _, pb := range projectedBounties {
		totalProjectedBounty += pb
	}

	return resolved, projectedBounties, totalProjectedBounty
}

func computeEligibleConcerns(ctx context.Context, cloned gov.Cloned, prop motionproto.Motion) motionproto.Refs {
	eligible := motionproto.Refs{}
	for _, ref := range prop.RefTo {
		if waimea.AreEligible(ctx, cloned, ref.To, prop.ID, ref.Type) {
			eligible = append(eligible, ref)
		}
	}
	eligible.Sort()
	return eligible
}

func closeResolvedConcerns(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	cons motionproto.Motions,

) float64 {

	for _, con := range cons {
		// close resolved concerns, and transfer concern escrows to proposal-owned bounty account
		motionapi.CloseMotion_StageOnly(
			ctx,
			cloned,
			con.ID,
			motionproto.Accept,
			waimea.ProposalBountyAccountID(prop.ID), // account to send bounty to
			prop,                                    // proposal that resolves the issue
		)
	}

	return account.Get_Local(
		ctx,
		cloned.PublicClone(),
		waimea.ProposalBountyAccountID(prop.ID),
	).Assets.Balance(account.PluralAsset).Quantity
}

func loadApprovalPoll(
	ctx context.Context,
	cloned gov.Cloned,
	prop motionproto.Motion,

) ballotproto.AdTallyMargin {

	pollName := waimea.ProposalApprovalPollName(prop.ID)
	return ballotapi.Show_Local(ctx, cloned, pollName)
}

func disberseReviewersRewards(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	accepted bool,

) (Rewards, metric.Receipts, float64) {

	rewards := Rewards{}
	approvalPoll := loadApprovalPoll(ctx, cloned.PublicClone(), prop)

	var isWinner func(score float64) bool
	if accepted {
		isWinner = func(score float64) bool { return score > 0 }
	} else {
		isWinner = func(score float64) bool { return score < 0 }
	}

	// compute winner shares
	winnerShares := map[member.User]float64{}
	winnerCost := map[member.User]float64{}
	winnerTotalShares := 0.0
	winnerTotalCost, loserTotalCost := 0.0, 0.0
	for user, choices := range approvalPoll.Tally.ScoresByUser {
		ss := choices[waimea.ProposalBallotChoice]
		if isWinner(ss.Score) {
			winnerShares[user] = math.Abs(ss.Score)
			winnerTotalShares += math.Abs(ss.Score)
			winnerCost[user] = math.Abs(ss.Strength)
			winnerTotalCost += math.Abs(ss.Strength)
		} else {
			loserTotalCost += math.Abs(ss.Strength)
		}
	}

	// disberse reviewer rewards
	for user, choices := range approvalPoll.Tally.ScoresByUser {
		ss := choices[waimea.ProposalBallotChoice]
		if isWinner(ss.Score) {
			q := loserTotalCost*winnerShares[user]/winnerTotalShares + winnerCost[user]
			fmt.Printf("payout to %v is %v\n", user, q)
			payout := account.H(account.PluralAsset, q)
			rewards = append(rewards,
				Reward{
					To:     user,
					Amount: payout,
				},
			)
			// transfer reward
			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				waimea.ProposalRewardAccountID(prop.ID),
				member.UserAccountID(user),
				payout,
				fmt.Sprintf("reviewer reward for proposal %v", prop.ID),
			)
		}
	}

	// send remainder to matching fund
	receipts := metric.Receipts{}
	rewardAccount := account.Get_Local(ctx, cloned.PublicClone(), waimea.ProposalRewardAccountID(prop.ID))
	remainder := rewardAccount.Balance(account.PluralAsset).Quantity
	donation := account.H(account.PluralAsset, 0.0)
	if remainder > 0 {
		donation = account.H(
			account.PluralAsset,
			remainder,
		)
		account.Transfer_StageOnly(
			ctx,
			cloned.PublicClone(),
			waimea.ProposalRewardAccountID(prop.ID),
			waimea.PennyAccountID,
			donation,
			fmt.Sprintf("donation to penny jar for proposal %v", prop.ID),
		)
		receipts = append(
			receipts,
			metric.OneReceipt(
				waimea.PennyAccountID.MetricAccountID(),
				metric.ReceiptTypeDonation,
				donation.MetricHolding(),
			)...,
		)
	}

	rewards.Sort()
	return rewards, receipts, donation.Quantity
}
