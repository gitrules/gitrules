package regime

import (
	"context"

	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/proto/history/metric"
	"github.com/gitrules/gitrules/proto/history/trace"
	"github.com/gitrules/gitrules/proto/notice"
)

func Dry(ctx context.Context) context.Context {
	ctx = metric.Mute(ctx)
	ctx = trace.Mute(ctx)
	ctx = notice.Mute(ctx)
	ctx = git.MuteStaging(ctx)
	return ctx
}
