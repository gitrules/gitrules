package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ghlib "github.com/gitrules/gitrules/github/org/lib"
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/google/go-github/v66/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type InstallationHandler struct {
	githubapp.ClientCreator
	DeployRelease string `json:"deploy_release"`
}

func (h *InstallationHandler) Handles() []string {
	return []string{"installation"}
}

func (h *InstallationHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {

	ctx = ghlib.InitCtx(ctx)

	var event github.InstallationEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "parsing installation event payload")
	}

	base.Infof("installation event: type %v, delivery %v, payload %v", eventType, deliveryID, form.SprintJSON(event))

	action := event.GetAction()
	if action != "created" {
		base.Infof("ignoring non-create events")
		return nil
	}

	installation := event.GetInstallation()

	if targetType := strings.ToLower(installation.GetTargetType()); targetType != "organization" {
		base.Errorf("gitrules for orgs cannot deploy to individual accounts")
		return fmt.Errorf("installing gitrules for orgs on an individual account")
	}

	if repoSelection := installation.GetRepositorySelection(); strings.ToLower(repoSelection) != "selected" {
		base.Errorf("gitrules for orgs can be installed only on explicitly selected repos")
		return fmt.Errorf("gitrules for orgs can be installed only on explicitly selected repos")
	}

	installationID := installation.GetID()
	base.Infof("installation ID %v", installationID)

	client, err := h.NewAppClient()
	if err != nil {
		base.Errorf("acquiring app client (%v)", err)
		return err
	}

	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, &github.InstallationTokenOptions{})
	if err != nil {
		base.Errorf("acquiring installation token (%v)", err)
		return err
	}
	base.Infof("acquired token %v", form.SprintJSON(token))

	// for each repo in the installation
	for _, repo := range event.Repositories {
		cfg, err := must.Try1[api.Config](
			func() api.Config {
				r := ghlib.ParseRepo(ctx, repo.GetFullName())
				return ghlib.Deploy(ctx, token.GetToken(), r, r, h.DeployRelease)
			},
		)
		if err != nil {
			base.Errorf("deploying (%v)", err)
			continue
		}
		base.Infof("deployed: %v", form.SprintJSON(cfg))
	}

	return nil
}
