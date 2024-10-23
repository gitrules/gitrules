package tools

import (
	"context"
	"sync"

	govgh "github.com/gitrules/gitrules/github"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

func ClearComments(
	ctx context.Context,
	token string,
	repo govgh.Repo,
	issueNo int64,

) {

	// create authenticated GitHub client
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	ghc := github.NewClient(tc)

	const maxGoroutines = 5
	throttle := make(chan struct{}, maxGoroutines)

	for {
		opts := &github.IssueListCommentsOptions{}
		comments, _, err := ghc.Issues.ListComments(ctx, repo.Owner, repo.Name, int(issueNo), opts)
		must.NoError(ctx, err)
		if len(comments) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, comment := range comments {
			throttle <- struct{}{}
			wg.Add(1)
			go func(comment *github.IssueComment) {
				defer wg.Done()
				defer func() { <-throttle }()
				base.Infof("Deleting comment %v from issue %v", comment.GetID(), issueNo)
				_, err := ghc.Issues.DeleteComment(ctx, repo.Owner, repo.Name, comment.GetID())
				must.NoError(ctx, err)
			}(comment)
		}
		wg.Wait()
	}
}
