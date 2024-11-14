package lib

import (
	"context"

	liborg "github.com/gitrules/gitrules/github/org/lib"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/google/go-github/v66/github"
)

func PublishDashboard(
	ctx context.Context,
	repo liborg.Repo,
	ghc *github.Client,
	cloned id.Cloned,
) {
}
