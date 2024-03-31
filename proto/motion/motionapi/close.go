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

func CloseMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
	decision motionproto.Decision,
	args ...any,

) (motionproto.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := CloseMotion_StageOnly(ctx, cloned, id, decision, args...)
	proto.Commitf(ctx, cloned.Public, "motion_close", "Close motion %v", id)
	return report, notices
}

func CloseMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id motionproto.MotionID,
	decision motionproto.Decision,
	args ...any,

) (motionproto.Report, notice.Notices) {

	t := cloned.Public.Tree()
	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, motionproto.ErrMotionAlreadyClosed)

	// apply policy
	pcy := motionproto.GetPolicy(ctx, motion.Policy)
	report, notices := pcy.Close(
		ctx,
		cloned,
		motion,
		decision,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit closure
	motion.Closed = true
	motion.ClosedAt = time.Now()
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "motion_close",
		Args:   trace.M{"id": id},
		Result: trace.M{"motion": motion},
	})

	return report, notices
}
