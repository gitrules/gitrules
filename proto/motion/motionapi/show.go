package motionapi

import (
	"context"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

func ShowMotion(
	ctx context.Context,
	addr gov.Address,
	id motionproto.MotionID,
	args ...any,

) motionproto.MotionView {

	return ShowMotion_Local(ctx, gov.Clone(ctx, addr), id, args...)
}

func ShowMotion_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,
	args ...any,
) motionproto.MotionView {

	t := cloned.Tree()
	m := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)

	mv, err := must.Try1[motionproto.MotionView](
		func() motionproto.MotionView {
			p := motionproto.GetPolicy(ctx, m.Policy)
			pv, pb := p.Show(ctx, cloned, m, args...)

			return motionproto.MotionView{
				Motion:  m,
				Ballots: pb,
				Policy:  pv,
			}
		},
	)
	if err != nil {
		return motionproto.MotionView{
			Motion: m,
		}
	}
	return mv
}
