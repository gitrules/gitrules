package motionapi

import (
	"context"
	"fmt"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

func ScoreMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) git.Change[form.Map, motionproto.Motions] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := ScoreMotions_StageOnly(ctx, cloned, args...)
	return proto.CommitIfChanged(ctx, cloned.Public, chg)
}

func ScoreMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) git.Change[form.Map, motionproto.Motions] {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	for i, motion := range motions {
		if motion.Archived || motion.Closed {
			continue
		}
		p := motionproto.GetMotionPolicy(ctx, motion)
		// NOTE: motion structure may change during scoring (if Score calls motion methods)
		score, notices := p.Score(
			ctx,
			cloned,
			motion,
			args...,
		)
		AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), motions[i].ID, notices)

		// reload motion, update score and save
		m := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, motions[i].ID)
		m.Score = score
		motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, motions[i].ID, m)
	}

	motions.Sort()

	return git.NewChange(
		fmt.Sprintf("Score all %d motions", len(motions)),
		"motion_score_all",
		form.Map{},
		motions,
		form.Forms{},
	)
}
