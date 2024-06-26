package git

import (
	"context"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Proxy interface {
	// CloneOne fetches just the branch specified in addr then checks out the branch in addr, creating it if not present.
	CloneOne(ctx context.Context, addr Address) Cloned
	// CloneAll fetches all branches then checks out the branch in addr, creating it if not present.
	CloneAll(ctx context.Context, addr Address) Cloned
}

type Cloned interface {
	// Push all branches to the origin.
	Push(context.Context)
	// Pull the branches indicated by the clone call that created this clone.
	Pull(context.Context)
	Repo() *Repository
	Tree() *Tree
}

func InitInMemory(ctx context.Context) *Repository {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	must.NoError(ctx, err)
	return repo
}
