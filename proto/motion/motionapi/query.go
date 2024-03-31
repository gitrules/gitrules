package motionapi

import (
	"context"

	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
)

func IsMotion(ctx context.Context, addr gov.Address, id motionproto.MotionID) bool {
	return IsMotion_Local(ctx, gov.Clone(ctx, addr).Tree(), id)
}

func IsMotion_Local(ctx context.Context, t *git.Tree, id motionproto.MotionID) bool {
	err := must.Try(func() { motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id) })
	return err == nil
}
