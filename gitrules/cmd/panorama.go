package cmd

import (
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/proto/panorama"
	"github.com/spf13/cobra"
)

var (
	panoramaCmd = &cobra.Command{
		Use:   "panorama",
		Short: "Panoramic view of user and motions, capturing pending votes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					return panorama.Panorama(ctx, setup.Gov, setup.Member)
				},
			)
		},
	}
)
