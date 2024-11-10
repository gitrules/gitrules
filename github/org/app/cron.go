package app

import (
	"context"
	"strings"
	"time"

	"github.com/gitrules/gitrules/github/org/lib"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/google/go-github/v66/github"
	"github.com/palantir/go-githubapp/githubapp"
)

const (
	AppInstallationTokenLifetime = time.Hour
	DeployTokenLifetime          = AppInstallationTokenLifetime / 2
	CronFrequency                = AppInstallationTokenLifetime / 3
)

func updateCron(ctx context.Context, cc githubapp.ClientCreator) {
	for {
		runStamp := time.Now()
		err := must.Try(
			func() {
				updateInstallations(ctx, cc)
			},
		)
		if err != nil {
			base.Errorf("updating installations (%v)", err)
		}

		nextRun := runStamp.Add(CronFrequency)
		time.Sleep(nextRun.Sub(time.Now()))
	}
}

func updateInstallations(ctx context.Context, cc githubapp.ClientCreator) {

	ac, err := cc.NewAppClient()
	must.NoError(ctx, err)

	var all []*github.Installation
	opt := &github.ListOptions{PerPage: 100}

	for {
		installations, resp, err := ac.Apps.ListInstallations(ctx, opt)
		must.NoError(ctx, err)
		all = append(all, installations...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, installation := range all {
		base.Infof("updating installation id=%d account=%s", installation.GetID(), installation.Account.GetLogin())
		err := must.Try(func() {
			updateInstallation(ctx, cc, installation)
		})
		if err != nil {
			base.Errorf("updating installation (%v)", err)
		}
	}
}

func updateInstallation(ctx context.Context, cc githubapp.ClientCreator, installation *github.Installation) {

	if strings.ToLower(installation.GetRepositorySelection()) != "selected" {
		base.Infof("ignoring installations without explicit repo selection")
		return
	}

	ic, err := cc.NewInstallationClient(installation.GetID())
	must.NoError(ctx, err)

	// list repos where app installed
	var all []*github.Repository
	opt := &github.ListOptions{PerPage: 100}

	for {
		repos, resp, err := ic.Apps.ListRepos(ctx, opt)
		must.NoError(ctx, err)
		all = append(all, repos.Repositories...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	base.Infof("planning to update repositories: %v", form.SprintJSON(repoFullNames(all)))

	// for each repo
	for _, repo := range all {
		err := must.Try(
			func() {
				updateRepo(ctx, cc, ic, installation, repo)
			},
		)
		if err != nil {
			base.Errorf("updating installation at repository %v (%v)", repo.GetFullName(), err)
		}
	}
}

func repoFullNames(repos []*github.Repository) []string {
	names := make([]string, len(repos))
	for i, repo := range repos {
		names[i] = repo.GetFullName()
	}
	return names
}

func updateRepo(
	ctx context.Context,
	cc githubapp.ClientCreator,
	ic *github.Client,
	installation *github.Installation,
	repo *github.Repository,
) {

	curSecret, _, err := ic.Actions.GetEnvSecret(ctx, int(repo.GetID()), lib.DeployEnvName, lib.DeployEnvOrganizerToken)
	must.NoError(ctx, err)

	now := time.Now()

	// if the token has been updated recently, do nothing
	if now.Sub(curSecret.UpdatedAt.UTC()) < DeployTokenLifetime {
		return
	}

	// create a new app installation token
	ac, err := cc.NewAppClient()
	must.NoError(ctx, err)

	token, _, err := ac.Apps.CreateInstallationToken(ctx, installation.GetID(), &github.InstallationTokenOptions{})
	must.NoError(ctx, err)

	// update access token
	pubKey, _, err := ic.Actions.GetEnvPublicKey(ctx, int(repo.GetID()), lib.DeployEnvName)
	must.NoError(ctx, err)

	encryptedToken := lib.EncryptSecret(ctx, pubKey, token.GetToken())

	newSecret := &github.EncryptedSecret{
		Name:           lib.DeployEnvOrganizerToken,
		KeyID:          pubKey.GetKeyID(),
		EncryptedValue: encryptedToken,
	}

	_, err = ic.Actions.CreateOrUpdateEnvSecret(ctx, int(repo.GetID()), lib.DeployEnvName, newSecret)
	must.NoError(ctx, err)

	base.Infof("updated organizer token for %v", repo.GetFullName())
}
