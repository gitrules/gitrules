package id

import (
	"testing"

	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/runtime"
)

func TestInit(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	testID := NewTestID(ctx, t, git.MainBranch, true)
	Init(ctx, testID.OwnerAddress())
	if err := must.Try(func() { Init(ctx, testID.OwnerAddress()) }); err == nil {
		t.Fatal("second init must fail")
	}
}
