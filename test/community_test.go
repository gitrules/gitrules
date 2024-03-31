package test

import (
	"testing"

	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/runtime"
)

func TestTestCommunity(t *testing.T) {
	// base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	NewTestCommunity(t, ctx, 3)
	// testutil.Hang()
}
