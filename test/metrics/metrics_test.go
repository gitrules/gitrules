package metrics

import (
	"fmt"
	"os"
	"testing"

	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/proto/metrics"
	"github.com/gitrules/gitrules/runtime"
	"github.com/gitrules/gitrules/test"
	pmp "github.com/gitrules/gitrules/test/motion/pmp_0"
)

func TestDashboardPMP(t *testing.T) {
	t.SkipNow()
	os.Setenv("PATH", "/opt/homebrew/bin")

	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	pmp.SetupTest(t, ctx, cty)

	urlCalc := func(assetRepoPath string) (url string) {
		return assetRepoPath
	}
	report := metrics.AssembleReport(ctx, cty.Gov(), urlCalc, metrics.TimeDailyLowerBound, metrics.Today().AddDate(0, 0, 1))
	fmt.Println(report.ReportMD)

	if report.Series.AllTime.DailyNumConcernVotes.Total() != 2 {
		t.Errorf("expecting 2, got %v", report.Series.AllTime.DailyNumConcernVotes.Total())
	}
}
