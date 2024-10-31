package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gitrules/gitrules/github/lib"
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/google/go-github/v66/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

const GitRulesDeployRelease = "vX.X.X" //XXX

type InstallationHandler struct {
	githubapp.ClientCreator
}

func (h *InstallationHandler) Handles() []string {
	return []string{"installation"}
}

func (h *InstallationHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {

	var event github.InstallationEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "parsing installation event payload")
	}

	base.Infof("installation event: type %v, delivery %v, payload %v", eventType, deliveryID, form.SprintJSON(event))

	action := event.GetAction()
	if action != "created" {
		return nil
	}

	if org := event.GetOrg(); org == nil {
		return fmt.Errorf("installing GitRules org app on a GitHub individual account")
	}

	installation := event.GetInstallation()
	installationID := installation.GetID()

	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		base.Errorf("acquiring installation client (%v)", err)
		return err
	}

	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, &github.InstallationTokenOptions{})
	if err != nil {
		base.Errorf("acquiring installation token (%v)", err)
		return err
	}

	// for each repo in the installation
	for _, repo := range event.Repositories {
		cfg, err := must.Try1[api.Config](func() api.Config {
			r := lib.FromGithubRepo(repo)
			return lib.Deploy(ctx, token.GetToken(), r, r, GitRulesDeployRelease)
		})
		if err != nil {
			base.Errorf("deploying (%v)", err)
			continue
		}
		base.Infof("deployed: %v", form.SprintJSON(cfg))
	}

	return nil
}
