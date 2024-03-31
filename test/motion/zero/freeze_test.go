package zero

import (
	"testing"

	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/proto/motion/motionapi"
	"github.com/gitrules/gitrules/proto/motion/motionpolicies/zero"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
	"github.com/gitrules/gitrules/runtime"
	"github.com/gitrules/gitrules/test"
)

func TestFreeze(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	id := motionproto.MotionID("123")

	// open
	motionapi.OpenMotion(
		ctx,
		cty.Organizer(),
		id,
		motionproto.MotionConcernType,
		zero.ZeroPolicyName,
		cty.MemberUser(0),
		"concern #1",
		"description #1",
		"https://1",
		nil)

	// freeze
	motionapi.FreezeMotion(ctx, cty.Organizer(), id)

	// score
	motionapi.ScoreMotions(ctx, cty.Organizer())

	// verify state changed
	m := motionapi.ShowMotion(ctx, cty.Gov(), id)
	if !m.Motion.Frozen {
		t.Errorf("expecting frozen")
	}

	// unfreeze
	motionapi.UnfreezeMotion(ctx, cty.Organizer(), id)

	m = motionapi.ShowMotion(ctx, cty.Gov(), id)
	if m.Motion.Frozen {
		t.Errorf("expecting not frozen")
	}

	// testutil.Hang()
}
