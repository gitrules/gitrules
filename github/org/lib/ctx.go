package lib

import (
	"context"

	"github.com/gitrules/gitrules/lib/git"
)

func InitCtx(ctx context.Context) context.Context {
	return WithTokenSource(git.WithTTL(git.WithAuth(ctx, nil), nil), nil)
}
