package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gitrules/gitrules"
	gh "github.com/gitrules/gitrules/github/lib"
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	_ "github.com/gitrules/gitrules/runtime"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gitrules",
		Short: "gitrules is a command-line client for the gitrules community governance protocol",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

var ctx = gh.WithTokenSource(git.WithTTL(git.WithAuth(context.Background(), nil), nil), nil)

var (
	configPath     string
	verbose        bool
	cpuProfilePath string
	memProfilePath string
)

func init() {
	cobra.OnInitialize(initAfterFlags)

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.gitrules/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "run in developer mode with verbose logging")
	rootCmd.PersistentFlags().StringVarP(&cpuProfilePath, "cpu", "p", "", "cpu profile path")
	rootCmd.PersistentFlags().StringVarP(&memProfilePath, "mem", "m", "", "memory profile path")

	rootCmd.AddCommand(initIDCmd)
	rootCmd.AddCommand(initGovCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(memberCmd)
	rootCmd.AddCommand(ballotCmd)
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(bureauCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(cronCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(githubCmd)
	rootCmd.AddCommand(motionCmd)
	rootCmd.AddCommand(etcCmd)
	rootCmd.AddCommand(panoramaCmd)
}

func initAfterFlags() {
	if verbose {
		base.LogVerbosely()
	} else {
		base.LogQuietly()
	}
	base.Infof("gitrules version: %v, os: %v, arch: %v", gitrules.Short(), runtime.GOOS, runtime.GOARCH)
	api.SetCPUProfilePath(cpuProfilePath)
	api.SetMemProfilePath(memProfilePath)
}

func LoadConfig() {
	if configPath == "" {
		// find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			base.Fatalf("looking for home dir (%v)", err)
		}
		base.AssertNoErr(err)

		// search for config in ~/.gitrules/config.json
		configPath = filepath.Join(home, api.LocalAgentPath, "config.json")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		base.Fatalf("reading config file (%v)", err)
	}

	config, err := form.DecodeBytes[api.Config](ctx, data)
	if err != nil {
		base.Fatalf("decoding config file (%v)", err)
	}

	if config.CacheDir != "" {
		ctx = git.WithCache(ctx, config.CacheDir)
	}

	setup = config.Setup(ctx)
}

var (
	setup api.Setup
)

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return 0
}

func ExecuteWithConfig(cfgPath string) int {
	configPath = cfgPath
	return Execute()
}
