package cmd

import (
	"github.com/gitrules/gitrules"
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version and build information",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					return gitrules.GetVersionInfo()
				},
			)
		},
	}
)
