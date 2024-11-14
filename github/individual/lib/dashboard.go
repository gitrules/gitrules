package lib

import (
	"context"

	"github.com/gitrules/gitrules/github/common"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/google/go-github/v66/github"
)

func PublishDashboard(
	ctx context.Context,
	repo common.Repo,
	ghc *github.Client,
	cloned id.Cloned,
) {
}
