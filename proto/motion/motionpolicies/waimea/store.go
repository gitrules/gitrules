package waimea

import (
	"context"

	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionapi"
)

func LoadConcernClassState_Local(ctx context.Context, cloned gov.OwnerCloned) *ConcernPolicyState {
	return motionapi.LoadClassState_Local[*ConcernPolicyState](ctx, cloned, ConcernPolicyName)
}

func SaveConcernClassState_StageOnly(ctx context.Context, cloned gov.OwnerCloned, ps *ConcernPolicyState) {
	motionapi.SaveClassState_StageOnly[*ConcernPolicyState](ctx, cloned, ConcernPolicyName, ps)
}
