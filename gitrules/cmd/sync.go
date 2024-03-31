package cmd

import (
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/proto/sync"
	"github.com/spf13/cobra"
)

var (
	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync governance with the community",
		Long: `
Sync is the heartbeat that advances the state of the governance forward.
Sync fetches all outstanding votes from community users and incorporates them in ballot tallies.
Sync also fetches and processes all outstanding service requests from community users.`,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					chg := sync.Sync(ctx, setup.Organizer, syncFetchPar)
					return chg.Result
				},
			)
		},
	}
)

var (
	syncFetchPar int
)

func init() {
	syncCmd.Flags().IntVar(&syncFetchPar, "fetch_par", 5, "parallelism while clonging member repos for vote collection")
}
