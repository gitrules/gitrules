package motionapi

import (
	"context"

	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
)

func Pipeline(
	ctx context.Context,
	addr gov.OwnerAddress,

) {

	cloned := gov.CloneOwner(ctx, addr)
	Pipeline_StageOnly(ctx, cloned)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_pipeline", "motion pipeline")
	cloned.PublicClone().Push(ctx)
}

func Pipeline_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,

) {

	// update and aggregate motion policies
	for i := 0; i < 2; i++ {
		base.Infof("PIPELINE: updating motions")
		UpdateMotions_StageOnly(ctx, cloned)
		base.Infof("PIPELINE: aggregating motions")
		AggregateMotions_StageOnly(ctx, cloned)
	}

	// rescore motions to capture updated tallies
	base.Infof("PIPELINE: scoring motions")
	ScoreMotions_StageOnly(ctx, cloned)

	// clearance
	base.Infof("PIPELINE: clear motions")
	ClearMotions_StageOnly(ctx, cloned)

	// archive closed motions
	base.Infof("PIPELINE: archive motions")
	ArchiveMotions_StageOnly(ctx, cloned)

}
