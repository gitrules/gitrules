package common

import (
	"context"

	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/lib/provider"
	ghprovider "github.com/gitrules/gitrules/lib/provider/github"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/google/go-github/v66/github"
)

type PublicPrivateRepos struct {
	Public     Repo
	Private    Repo
	PublicURLs  *provider.Repository
	PrivateURLs *provider.Repository
}

func CreatePublicPrivateRepos(
	ctx context.Context,
	ghc *github.Client,
	project Repo,
	govPrefix Repo,

) (PublicPrivateRepos, id.OwnerAddress) {

	// create governance public and private repos
	v := ghprovider.NewGithubVendorWithClient(ctx, ghc)

	govPublic := Repo{
		Owner: govPrefix.Owner,
		Name:  govPrefix.Name + GitRulesPublicSuffix,
	}
	base.Infof("creating GitHub repository %v", govPublic)
	govPublicURLs, err := v.CreateRepo(ctx, govPublic.Name, govPublic.Owner, false)
	must.NoError(ctx, err)

	govPrivate := Repo{
		Owner: govPrefix.Owner,
		Name:  govPrefix.Name + GitRulesPrivateSuffix,
	}
	base.Infof("creating GitHub repository %v", govPrivate)
	govPrivateURLs, err := v.CreateRepo(ctx, govPrivate.Name, govPrivate.Owner, true)
	must.NoError(ctx, err)

	return PublicPrivateRepos{
			Public:     govPublic,
			Private:    govPrivate,
			PublicURLs:  govPublicURLs,
			PrivateURLs: govPrivateURLs,
		}, id.OwnerAddress{
			Public: id.PublicAddress{
				Repo:   git.URL(govPublicURLs.HTTPSURL),
				Branch: git.MainBranch,
			},
			Private: id.PrivateAddress{
				Repo:   git.URL(govPrivateURLs.HTTPSURL),
				Branch: git.MainBranch,
			},
		}
}
