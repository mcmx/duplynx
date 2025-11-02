package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/mcmx/duplynx/internal/config"
)

var runtimeCfg = config.DefaultRuntimeConfig()

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "duplynx",
		Short: "DupLynx developer CLI",
		Long:  "DupLynx developer CLI for running demo servers and data maintenance commands.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			runtimeCfg.ApplyEnvOverrides(cmd.Flags())
			cfg := runtimeCfg
			ctx := config.WithRuntimeConfig(cmd.Context(), cfg)
			cmd.SetContext(ctx)
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.SetVersionTemplate("DupLynx CLI\n")
	cmd.PersistentFlags().SortFlags = false

	config.BindRuntimeFlags(cmd.PersistentFlags(), &runtimeCfg)

	cmd.AddCommand(newServeCommand())

	cmd.SetContext(context.Background())
	return cmd
}
