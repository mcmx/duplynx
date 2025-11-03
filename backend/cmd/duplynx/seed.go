package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mcmx/duplynx/internal/config"
	"github.com/mcmx/duplynx/internal/data"
	"github.com/mcmx/duplynx/internal/observability"
)

func newSeedCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Rebuild the DupLynx demo database with canonical fixtures",
		RunE:  runSeed,
	}
	return cmd
}

func runSeed(cmd *cobra.Command, _ []string) (err error) {
	ctx := cmd.Context()

	cfg, ok := config.FromContext(ctx)
	if !ok {
		cfg = runtimeCfg
	}

	actor := resolveActor()
	metadata := map[string]any{
		"db_file":    cfg.DBFile,
		"assets_dir": cfg.AssetsDir,
	}
	writer := observability.NewEventWriter(nil)
	start := time.Now()

	writer.Write(observability.Event{
		Action:   "seed_start",
		Actor:    actor,
		Outcome:  "starting",
		Metadata: metadata,
	})

	defer func() {
		outcome := "success"
		if err != nil {
			outcome = "failure"
		}
		writer.Write(observability.Event{
			Action:   "seed_stop",
			Actor:    actor,
			Outcome:  outcome,
			Duration: time.Since(start),
			Metadata: metadata,
			Error:    err,
		})
	}()

	client, err := data.OpenSQLite(ctx, cfg.SQLiteDSN())
	if err != nil {
		return err
	}
	defer func() {
		closeErr := data.Close(client)
		if err == nil {
			err = closeErr
		}
	}()

	report, err := data.SeedDemoDataset(ctx, client, data.CanonicalDemoDataset())
	if err != nil {
		return err
	}

	metadata["tenants"] = report.Tenants
	metadata["machines"] = report.Machines
	metadata["scans"] = report.Scans
	metadata["duplicate_groups"] = report.DuplicateGroups

	fmt.Fprintf(cmd.OutOrStdout(),
		"Seeded %d tenants, %d machines, %d scans, %d duplicate groups, and %d file instances\n",
		report.Tenants,
		report.Machines,
		report.Scans,
		report.DuplicateGroups,
		report.FileInstances,
	)

	return nil
}
