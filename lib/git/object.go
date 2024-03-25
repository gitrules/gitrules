package git

import (
	"context"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetCommit(ctx context.Context, r *Repository, h plumbing.Hash) *object.Commit {
	c, err := object.GetCommit(r.Storer, h)
	must.NoError(ctx, err)
	return c
}

func GetTree(ctx context.Context, r *Repository, th plumbing.Hash) *object.Tree {
	tree, err := object.GetTree(r.Storer, th)
	must.NoError(ctx, err)
	return tree
}

func GetCommitTree(ctx context.Context, r *Repository, commitHash plumbing.Hash) *object.Tree {
	return GetTree(ctx, r, GetCommit(ctx, r, commitHash).TreeHash)
}

func GetBranchTree(ctx context.Context, r *Repository, branch Branch) *object.Tree {
	return GetTree(ctx, r, ResolveBranch(ctx, r, branch).TreeHash)
}
