// Package pmp implements the Plural Management Protocol.
package pmp_1

import (
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/motion"
	"github.com/gitrules/gitrules/proto/motion/motionpolicies/pmp_0"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

const (
	ConcernBallotChoice  = pmp_0.ConcernBallotChoice
	ProposalBallotChoice = pmp_0.ProposalBallotChoice

	ConcernPolicyName  motion.PolicyName = "pmp-concern-policy-v1"
	ProposalPolicyName motion.PolicyName = "pmp-proposal-v1"

	ConcernPolicyGithubLabel  = "gitrules:pmp-v1"
	ProposalPolicyGithubLabel = ConcernPolicyGithubLabel

	ClaimsRefType = pmp_0.ClaimsRefType
)

func ConcernAccountID(id motionproto.MotionID) account.AccountID {
	return pmp_0.ConcernAccountID(id)
}

func ProposalAccountID(id motionproto.MotionID) account.AccountID {
	return pmp_0.ProposalAccountID(id)
}

func ConcernPollBallotName(id motionproto.MotionID) ballotproto.BallotID {
	return pmp_0.ConcernPollBallotName(id)
}

func ProposalApprovalPollName(id motionproto.MotionID) ballotproto.BallotID {
	return pmp_0.ProposalApprovalPollName(id)
}

func ProposalBountyAccountID(motionID motionproto.MotionID) account.AccountID {
	return pmp_0.ProposalBountyAccountID(motionID)
}

func ProposalRewardAccountID(motionID motionproto.MotionID) account.AccountID {
	return pmp_0.ProposalRewardAccountID(motionID)
}
