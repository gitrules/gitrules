package zero

import (
	"testing"

	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/proto/motion/motionapi"
	"github.com/gitrules/gitrules/proto/motion/motionpolicies/zero"
	"github.com/gitrules/gitrules/proto/motion/motionproto"
	"github.com/gitrules/gitrules/runtime"
	"github.com/gitrules/gitrules/test"
)

func TestOpenCancel(t *testing.T) {
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

	// list
	ms := motionapi.ListMotions(ctx, cty.Gov())
	if len(ms) != 1 {
		t.Errorf("expecting 1 motion, got %v", len(ms))
	}

	// give credits to user
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 13.0), "test")

	// score
	motionapi.ScoreMotions(ctx, cty.Organizer())

	// cancel
	motionapi.CancelMotion(ctx, cty.Organizer(), id)

	// verify state changed
	m := motionapi.ShowMotion(ctx, cty.Gov(), id)
	if !m.Motion.Closed || !m.Motion.Cancelled {
		t.Errorf("expecting closed and cancelled")
	}

	// testutil.Hang()
}
