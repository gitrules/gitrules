package motionapi

import (
	"context"

	"github.com/gitrules/gitrules/proto/ballot/ballotapi"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
	"github.com/gitrules/gitrules/proto/regime"
)

func TrackMotion(
	ctx context.Context,
	addr gov.Address,
	voterAddr id.OwnerAddress,
	mid motionproto.MotionID,

) motionproto.MotionView {

	voterOwner := id.CloneOwner(ctx, voterAddr)
	return TrackMotion_Local(ctx, gov.Clone(ctx, addr), voterAddr, voterOwner, mid)
}

func TrackMotion_Local(
	ctx context.Context,
	cloned gov.Cloned,
	voterAddr id.OwnerAddress,
	voterOwner id.OwnerCloned,
	mid motionproto.MotionID,

) motionproto.MotionView {

	ctx = regime.Dry(ctx)

	mv := ShowMotion_Local(ctx, cloned, mid)
	if len(mv.Ballots) > 0 {
		vs := ballotapi.Track_StageOnly(
			ctx,
			voterAddr,
			voterOwner,
			cloned,
			mv.Ballots[0].BallotID,
		)
		mv.Voter = &vs
	}
	return mv
}

func TrackMotionBatch(
	ctx context.Context,
	addr gov.Address,
	voterAddr id.OwnerAddress,

) motionproto.MotionViews {

	voterOwner := id.CloneOwner(ctx, voterAddr)
	return TrackMotionBatch_Local(ctx, gov.Clone(ctx, addr), voterAddr, voterOwner)
}

func TrackMotionBatch_Local(
	ctx context.Context,
	cloned gov.Cloned,
	voterAddr id.OwnerAddress,
	voterOwner id.OwnerCloned,

) motionproto.MotionViews {

	ctx = regime.Dry(ctx)

	mids := motionproto.MotionKV.ListKeys(ctx, motionproto.MotionNS, cloned.Tree())
	mvs := make(motionproto.MotionViews, len(mids))
	for i, mid := range mids {
		mvs[i] = TrackMotion_Local(ctx, cloned, voterAddr, voterOwner, mid)
	}
	mvs.Sort()
	return mvs
}
