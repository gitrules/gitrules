package cmd

import (
	"encoding/json"
	"io"
	"os"

	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/proto/etc"
	"github.com/spf13/cobra"
)

var (
	etcCmd = &cobra.Command{
		Use:   "etc",
		Short: "Manage system settings",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	etcGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get system settings",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					settings := etc.GetSettings(ctx, setup.Gov)
					return settings
				},
			)
		},
	}

	etcSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set system settings",
		Long:  `System settings must be given as JSON on the standard input.`,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()

					jsonData, err := io.ReadAll(os.Stdin)
					must.NoError(ctx, err)

					var settings etc.Settings
					err = json.Unmarshal(jsonData, &settings)
					must.NoError(ctx, err)

					etc.SetSettings(ctx, setup.Gov, settings)
				},
			)
		},
	}
)

func init() {
	etcCmd.AddCommand(etcGetCmd)
	etcCmd.AddCommand(etcSetCmd)
}
