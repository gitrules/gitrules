package id

import (
	"context"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
)

func FetchPublicCredentials(ctx context.Context, addr PublicAddress) PublicCredentials {
	return GetPublicCredentials(ctx, git.CloneOne(ctx, git.Address(addr)).Tree())
}

func GetPublicCredentials(ctx context.Context, t *git.Tree) PublicCredentials {
	cred := form.FromFile[PublicCredentials](ctx, t.Filesystem, PublicCredentialsNS)
	must.Assertf(ctx, cred.IsValid(), "credentials are not valid")
	return cred
}
