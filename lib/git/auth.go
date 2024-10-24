package git

import (
	"context"
	"sync"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func MakePasswordAuth(ctx context.Context, user string, pass string) transport.AuthMethod {
	return &http.BasicAuth{Username: user, Password: pass}
}

func MakeTokenAuth(ctx context.Context, token string) transport.AuthMethod {
	return &http.BasicAuth{Username: "123", Password: token} // "123" can be anything but empty
}

func MakeSSHFileAuth(ctx context.Context, user string, privKeyFile string) transport.AuthMethod {
	pubKey, err := ssh.NewPublicKeysFromFile(user, privKeyFile, "")
	must.NoError(ctx, err)
	return pubKey
}

// auth manager in context

type contextKeyAuthManager struct{}

func WithAuth(ctx context.Context, am *AuthManager) context.Context {
	if am == nil {
		am = NewAuthManager()
	}
	return context.WithValue(ctx, contextKeyAuthManager{}, am)
}

func SetAuth(ctx context.Context, forRepo URL, a transport.AuthMethod) {
	ctx.Value(contextKeyAuthManager{}).(*AuthManager).SetAuth(forRepo, a)
}

func GetAuth(ctx context.Context, forRepo URL) transport.AuthMethod {
	if am, ok := ctx.Value(contextKeyAuthManager{}).(*AuthManager); ok {
		return am.GetAuth(forRepo)
	}
	return nil
}

// AuthManager provides authentication methods given a repo URL.
type AuthManager struct {
	lk  sync.Mutex
	url map[URL]transport.AuthMethod
}

func NewAuthManager() *AuthManager {
	return &AuthManager{url: map[URL]transport.AuthMethod{}}
}

func (x *AuthManager) SetAuth(forRepo URL, a transport.AuthMethod) {
	x.lk.Lock()
	defer x.lk.Unlock()
	x.url[forRepo] = a
}

func (x *AuthManager) GetAuth(forRepo URL) transport.AuthMethod {
	x.lk.Lock()
	defer x.lk.Unlock()
	return x.url[forRepo]
}
