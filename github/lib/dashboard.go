package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/gitrules/gitrules"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/lib/ns"
	"github.com/gitrules/gitrules/materials"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/gitrules/gitrules/proto/metrics"
	"github.com/google/go-github/v66/github"
)

// the dashboard is published and updated on the first issue that is labelled "gitrules:dashboard"

func PublishDashboard(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	cloned gov.Cloned,
) {

	assetsAddr := git.Address{
		Repo:   cloned.Address().Repo,
		Branch: cloned.Address().Branch + ".web-assets",
	}

	assetsRepo, err := ParseGithubRepoURL(string(assetsAddr.Repo))
	must.NoError(ctx, err)

	assets := metrics.AssembleReport_Local(
		ctx,
		cloned,
		func(assetRepoPath string) (url string) {
			return uploadedAssetURL(assetsRepo, string(assetsAddr.Branch), assetRepoPath)
		},
		metrics.TimeDailyLowerBound,
		metrics.Today().AddDate(0, 0, 1),
	)

	uploadAssets(ctx, assetsAddr, assets.Assets)

	header := fmt.Sprintf(
		"## <a href=%q><img src=%q alt=\"This project is governed with GitRules.\" width=\"65\" /></a> %s\n"+
			"On `%s` by GitRules `%s`\n\n",
		materials.GitRulesWebsiteURL,
		materials.GitRulesAvatarURL,
		"GitRules community dashboard",
		time.Now().Format(time.RFC850),
		gitrules.GetVersionInfo().Version,
	)
	updateDashboard(ctx, ghc, repo, "GitRules community dashboard", header+assets.ReportMD)
}

func uploadedAssetURL(repo Repo, branch string, gitPath string) string {
	return fmt.Sprintf(
		"https://raw.githubusercontent.com/%s/%s/%s/%s",
		repo.Owner,
		repo.Name,
		branch,
		gitPath,
	)
}

func uploadAssets(
	ctx context.Context,
	addr git.Address, // addr must be a GitHub repo URL
	assets map[string][]byte, // git path (in assets repo branch) -> content; git path must have no leading slashes

) {

	cloned := git.CloneOne(ctx, addr)
	for path, content := range assets {
		git.BytesToFileStage(ctx, cloned.Tree(), ns.ParseFromGitPath(path), content)
	}
	git.Commit(ctx, cloned.Tree(), "upload assets")
	cloned.Push(ctx)
}

func updateDashboard(
	ctx context.Context,
	ghc *github.Client,
	repo Repo,
	title string,
	body string,

) {

	labels := []string{DashboardIssueLabel}

	// check if there is an existing dashboard issue
	opt := &github.IssueListByRepoOptions{
		State:  "open",
		Labels: labels,
	}
	issues, _, err := ghc.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
	must.NoError(ctx, err)

	// create a dashboard issue if there is none
	req := &github.IssueRequest{
		Title:  github.String(title),
		Body:   github.String(body),
		Labels: &labels,
	}
	if len(issues) == 0 {
		_, _, err := ghc.Issues.Create(ctx, repo.Owner, repo.Name, req)
		must.NoError(ctx, err)
	} else {
		_, _, err := ghc.Issues.Edit(ctx, repo.Owner, repo.Name, issues[0].GetNumber(), req)
		must.NoError(ctx, err)
	}
}
