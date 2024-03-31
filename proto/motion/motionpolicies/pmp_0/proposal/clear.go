package proposal

import (
	"context"
	"fmt"
	"math"

	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/ballot/ballotapi"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/member"
	"github.com/gitrules/gitrules/proto/motion/motionapi"
	"github.com/gitrules/gitrules/proto/motion/motionpolicies/pmp_0"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

func loadResolvedConcerns(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,

) motionproto.Motions {

	eligible := computeEligibleConcerns(ctx, cloned.PublicClone(), prop)
	resolved := motionproto.Motions{}
	for _, ref := range eligible {
		con := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
		resolved = append(resolved, con)
	}
	return resolved
}

func computeEligibleConcerns(ctx context.Context, cloned gov.Cloned, prop motionproto.Motion) motionproto.Refs {
	eligible := motionproto.Refs{}
	for _, ref := range prop.RefTo {
		if pmp_0.IsConcernProposalEligible(ctx, cloned, ref.To, prop.ID, ref.Type) {
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

) account.Holding {

	for _, con := range cons {
		// close resolved concerns, and transfer concern escrows to proposal-owned bounty account
		motionapi.CloseMotion_StageOnly(
			ctx,
			cloned,
			con.ID,
			motionproto.Accept,
			pmp_0.ProposalBountyAccountID(prop.ID), // account to send bounty to
			prop,                                   // proposal that resolves the issue
		)
	}

	return account.Get_Local(
		ctx,
		cloned.PublicClone(),
		pmp_0.ProposalBountyAccountID(prop.ID),
	).Assets.Balance(account.PluralAsset)
}

func loadPropApprovalPollTally(
	ctx context.Context,
	cloned gov.Cloned,
	prop motionproto.Motion,

) ballotproto.AdTallyMargin {

	pollName := pmp_0.ProposalApprovalPollName(prop.ID)
	return ballotapi.Show_Local(ctx, cloned, pollName)
}

func disberseRewards(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,

) Rewards {

	rewards := Rewards{}
	adt := loadPropApprovalPollTally(ctx, cloned.PublicClone(), prop)

	// get reward account balance
	// totalWinnings := account.Get_Local(
	// 	ctx,
	// 	cloned.PublicClone(),
	// 	pmp.ProposalRewardAccountID(prop.ID),
	// ).Assets.Balance(account.PluralAsset).Quantity

	// compute reward distribution
	rewardFund := 0.0                      // total credits spent on negative votes
	totalCut := 0.0                        // sum of all positive votes
	winnerCut := map[member.User]float64{} // positive votes per user
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp_0.ProposalBallotChoice]
		if ss.Score <= 0.0 {
			// compute total credits spent on negative votes
			rewardFund += math.Abs(ss.Strength)
		} else {
			totalCut += ss.Score
			winnerCut[user] = ss.Score
		}
	}

	// payout winnings
	for user, choices := range adt.Tally.ScoresByUser {
		ss := choices[pmp_0.ProposalBallotChoice]
		if ss.Score > 0.0 {
			payout := account.H(
				account.PluralAsset,
				math.Abs(ss.Strength)+rewardFund*winnerCut[user]/totalCut,
			)
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
				pmp_0.ProposalRewardAccountID(prop.ID),
				member.UserAccountID(user),
				payout,
				fmt.Sprintf("reward for proposal %v", prop.ID),
			)
		}
	}

	rewards.Sort()
	return rewards
}
