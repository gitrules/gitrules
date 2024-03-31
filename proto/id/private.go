package id

import (
	"context"

	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
)

func FetchOwnerCredentials(ctx context.Context, addr OwnerAddress) PrivateCredentials {
	return GetOwnerCredentials(ctx, CloneOwner(ctx, addr))
}

func GetOwnerCredentials(ctx context.Context, owner OwnerCloned) PrivateCredentials {
	return GetPrivateCredentials(ctx, owner.Private.Tree())
}

func GetPrivateCredentials(ctx context.Context, privateTree *git.Tree) PrivateCredentials {
	cred := form.FromFile[PrivateCredentials](ctx, privateTree.Filesystem, PrivateCredentialsNS)
	must.Assertf(ctx, cred.IsValid(), "credentials are not valid")
	return cred
}
