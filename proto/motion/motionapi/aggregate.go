package motionapi

import (
	"context"

	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

func AggregateMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) {

	cloned := gov.CloneOwner(ctx, addr)
	AggregateMotions_StageOnly(ctx, cloned, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_aggregate", "Aggregate motions")
}

func AggregateMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)

	policyMotions := map[motion.PolicyName]motionproto.Motions{}

	for _, motion := range motions {
		if motion.Archived || motion.Closed {
			continue
		}
		policyMotions[motion.Policy] = append(policyMotions[motion.Policy], motion)
	}

	for policyName, policyMotions := range policyMotions {
		policy := motionproto.GetMotionPolicyByName(ctx, policyName)
		policyMotions.Sort()
		policy.Aggregate(ctx, cloned, policyMotions)
	}
}
