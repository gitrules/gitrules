package lib

import (
	"context"

	"github.com/gitrules/gitrules/lib/git"
)

func InitCommandCtx(ctx context.Context) context.Context {
	return WithTokenSource(git.WithTTL(git.WithAuth(ctx, nil), nil), nil)
}
