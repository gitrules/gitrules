package individual

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gitrules/gitrules"
	lib_individual "github.com/gitrules/gitrules/github/individual/lib"
	lib_org "github.com/gitrules/gitrules/github/org/lib"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/lib/ns"
	"github.com/gitrules/gitrules/proto"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/google/go-github/v66/github"
)

var CronNS = ns.NS{"cron", "cron.json"}

func Cron(
	ctx context.Context,
	repo lib_org.Repo,
	ghc *github.Client,
	addr id.OwnerAddress,
	//
	githubFreq time.Duration, // frequency of syncing with github
) form.Map {

	cloned := id.CloneOwner(ctx, addr)
	tree := cloned.Public.Tree()

	// use a separate branch for cron logs
	cronAddr := git.Address(addr.Public)
	cronAddr.Branch = cronAddr.Branch + ".cron"
	cronCloned := git.CloneOne(ctx, cronAddr)
	cronTree := cronCloned.Tree()

	// read cron state
	state, err := git.TryFromFile[CronState](ctx, cronTree, CronNS)
	must.Assertf(ctx, err == nil || err == os.ErrNotExist, "opening cron state (%v)", err)

	now := time.Now()
	shouldSyncGithub := now.Sub(state.LastGithubSync) > githubFreq

	report := form.Map{}

	if shouldSyncGithub {
		// nop
		state.LastGithubSync = time.Now()
	}

	// update community dashboard on github
	base.Infof("CRON: publishing individual dashboard")
	lib_individual.PublishDashboard(ctx, repo, ghc, cloned.PublicClone())

	// prepare commit message
	report["cron"] = state
	ver := gitrules.GetVersionInfo()
	latestChange := LatestChange{
		Stamp:           now,
		GitRulesVersion: ver,
	}

	git.ToFileStage[LatestChange](ctx, cloned.PublicClone().Tree(), LatestChangeMetaNS, latestChange)

	cronChg := git.NewChange[form.Map, LatestChange](
		fmt.Sprintf("GitRules for Individuals %s cron job.", ver.Version),
		"cron",
		nil,
		// We used to include the report in the commit message. However this causes a problem on GitHub.
		// The report includes the bodies of the issues that were processed.
		// It turns out GitHub scans the commit message for "resolves issue" text and automatically closes issues based on those.
		// This triggers spurious closures.
		latestChange,
		nil,
	)

	// push state, if changed
	status, err := tree.Status()
	must.NoError(ctx, err)
	if !status.IsClean() {
		proto.Commit(ctx, cloned.Public.Tree(), cronChg)
		cloned.Public.Push(ctx)
	}

	// always push cron state
	git.ToFileStage(ctx, cronTree, CronNS, state)
	proto.Commit(ctx, cronTree, cronChg)
	cronCloned.Push(ctx)

	return report
}

var LatestChangeMetaNS = ns.ParseFromGitPath("latest_change.json")

type LatestChange struct {
	Stamp           time.Time            `json:"change_stamp"`
	GitRulesVersion gitrules.VersionInfo `json:"gitrules_version"`
}

type CronState struct {
	LastGithubSync time.Time `json:"last_github_sync"`
}
