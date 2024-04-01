package waimea

import (
	"github.com/gitrules/gitrules/proto/ballot/ballotproto"
	"github.com/gitrules/gitrules/proto/motion"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

const (
	ConcernBallotChoice  = "rank"
	ProposalBallotChoice = "rank"

	ConcernPolicyName  motion.PolicyName = "waimea-concern"
	ProposalPolicyName motion.PolicyName = "waimea-proposal"

	ConcernPolicyGithubLabel  = "gitrules:waimea"
	ProposalPolicyGithubLabel = ConcernPolicyGithubLabel

	ClaimsRefType = "claims"
)

func ConcernPollBallotName(id motionproto.MotionID) ballotproto.BallotID {
	return ballotproto.BallotID("waimea/motion/priority_poll/" + id.String())
}

func ProposalApprovalPollName(id motionproto.MotionID) ballotproto.BallotID {
	return ballotproto.BallotID("waimea/motion/approval_poll/" + id.String())
}
