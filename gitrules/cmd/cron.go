package cmd

import (
	"time"

	govgh "github.com/gitrules/gitrules/github/org/lib"
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/proto/cron"
	"github.com/spf13/cobra"
)

var (
	cronCmd = &cobra.Command{
		Use:   "cron",
		Short: "Cron tasks",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cronOrgCmd = &cobra.Command{
		Use:   "org",
		Short: "cron performs time-dependent update operations to the governance system for orgs",
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
						time.Duration(cronOrgGithubFreqSeconds)*time.Second,
						time.Duration(cronOrgCommunityFreqSeconds)*time.Second,
						syncFetchPar,
					)
					return result
				},
			)
		},
	}

	cronIndividualCmd = &cobra.Command{
		Use:   "individual",
		Short: "cron for individuals",
		Long: `
This command is intended as a target for a cronjob which runs every couple of minutes.
`,
		Run: func(cmd *cobra.Command, args []string) {
			// api.Invoke1(
			// 	func() any {
			// 		LoadConfig()
			// 		repo := govgh.ParseRepo(ctx, githubProject)
			// 		govgh.SetTokenSource(ctx, repo, govgh.MakeStaticTokenSource(ctx, githubToken))
			// 		ghc := govgh.GetGithubClient(ctx, repo)
			// 		result := cron.Cron(
			// 			ctx,
			// 			repo,
			// 			ghc,
			// 			setup.Organizer,
			// 			time.Duration(cronOrgGithubFreqSeconds)*time.Second,
			// 			time.Duration(cronOrgCommunityFreqSeconds)*time.Second,
			// 			syncFetchPar,
			// 		)
			// 		return result
			// 	},
			// )
		},
	}
)

var (
	cronOrgGithubFreqSeconds    int
	cronOrgCommunityFreqSeconds int
)

func init() {
	cronCmd.AddCommand(cronOrgCmd)
	cronOrgCmd.Flags().StringVar(&githubProject, "project", "", "GitHub project owner/repo")
	cronOrgCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	cronOrgCmd.Flags().IntVar(&cronOrgGithubFreqSeconds, "github_freq", govgh.DefaultGithubFreq, "frequency of GitHub import, in seconds")
	cronOrgCmd.Flags().IntVar(&cronOrgCommunityFreqSeconds, "community_freq", govgh.DefaultCommunityFreq, "frequency of community tallies, in seconds")
	cronOrgCmd.Flags().IntVar(&syncFetchPar, "fetch_par", govgh.DefaultFetchParallelism, "parallelism while clonging member repos for vote collection")

	cronOrgCmd.MarkFlagRequired("project")
	cronOrgCmd.MarkFlagRequired("token")
	cronOrgCmd.MarkFlagRequired("github_freq")
	cronOrgCmd.MarkFlagRequired("community_freq")
	cronOrgCmd.MarkFlagRequired("fetch_par")

	cronCmd.AddCommand(cronIndividualCmd)
}
