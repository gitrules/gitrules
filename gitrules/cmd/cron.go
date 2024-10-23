package cmd

import (
	"time"

	govgh "github.com/gitrules/gitrules/github/lib"
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/proto/cron"
	"github.com/spf13/cobra"
)

var (
	cronCmd = &cobra.Command{
		Use:   "cron",
		Short: "cron performs time-dependent update operations to the governance system",
		Long: `
This command is intended as a target for a cronjob which runs every couple of minutes.
It will ensure that:
- Governance is synchronized with the issues and pull requests of a GitHub project at a configurable frequency, and
- Votes from community members are incorporated in governance ballots at a configurable frequency.
`,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					repo := govgh.ParseRepo(ctx, githubProject)
					govgh.SetTokenSource(ctx, repo, govgh.MakeStaticTokenSource(ctx, githubToken))
					ghc := govgh.GetGithubClient(ctx, repo)
					result := cron.Cron(
						ctx,
						repo,
						ghc,
						setup.Organizer,
						time.Duration(cronGithubFreqSeconds)*time.Second,
						time.Duration(cronCommunityFreqSeconds)*time.Second,
						syncFetchPar,
					)
					return result
				},
			)
		},
	}
)

var (
	cronGithubFreqSeconds    int
	cronCommunityFreqSeconds int
)

func init() {
	cronCmd.Flags().StringVar(&githubProject, "project", "", "GitHub project owner/repo")
	cronCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	cronCmd.Flags().IntVar(&cronGithubFreqSeconds, "github_freq", govgh.DefaultGithubFreq, "frequency of GitHub import, in seconds")
	cronCmd.Flags().IntVar(&cronCommunityFreqSeconds, "community_freq", govgh.DefaultCommunityFreq, "frequency of community tallies, in seconds")
	cronCmd.Flags().IntVar(&syncFetchPar, "fetch_par", govgh.DefaultFetchParallelism, "parallelism while clonging member repos for vote collection")

	cronCmd.MarkFlagRequired("project")
	cronCmd.MarkFlagRequired("token")
	cronCmd.MarkFlagRequired("github_freq")
	cronCmd.MarkFlagRequired("community_freq")
	cronCmd.MarkFlagRequired("fetch_par")
}
