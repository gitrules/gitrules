package motionapi

import (
	"context"

	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
	"github.com/gitrules/gitrules/proto/notice"
)

func ClearMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) ([]motionproto.Report, []notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := ClearMotions_StageOnly(ctx, cloned, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_clear", "Clear motions")
	return report, notices
}

func ClearMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) ([]motionproto.Report, []notice.Notices) {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	reportList := []motionproto.Report{}
	noticesList := []notice.Notices{}
	for i, motion := range motions {
		if motion.Archived {
			continue
		}
		p := motionproto.GetMotionPolicy(ctx, motion)
		report, notices := p.Clear(
			ctx,
			cloned,
			motion,
			args...,
		)
		reportList = append(reportList, report)
		noticesList = append(noticesList, notices)
		AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), motions[i].ID, notices)
	}

	motions.Sort()

	return reportList, noticesList
}
