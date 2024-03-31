package motionapi

import (
	"context"
	"time"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/history/trace"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
	"github.com/gitrules/gitrules/proto/notice"
)

func CancelMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
	args ...any,

) (motionproto.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := CancelMotion_StageOnly(ctx, cloned, id, args...)
	proto.Commitf(ctx, cloned.Public, "motion_cancel", "Cancel motion %v", id)
	return report, notices
}

func CancelMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id motionproto.MotionID,
	args ...any,

) (motionproto.Report, notice.Notices) {

	t := cloned.Public.Tree()
	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)
	must.Assertf(ctx, !motion.Closed, "motion %v already closed", motion.ID)
	must.Assertf(ctx, !motion.Cancelled, "motion %v already cancelled", motion.ID)

	// apply policy
	pcy := motionproto.GetPolicy(ctx, motion.Policy)
	report, notices := pcy.Cancel(
		ctx,
		cloned,
		motion,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit cancellation
	motion.Closed = true
	motion.Cancelled = true
	motion.ClosedAt = time.Now()
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "motion_cancel",
		Args:   trace.M{"id": id},
		Result: trace.M{"motion": motion},
	})

	return report, notices
}
