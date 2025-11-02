package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/mcmx/duplynx/internal/app"
	"github.com/mcmx/duplynx/internal/config"
	"github.com/mcmx/duplynx/internal/data"
	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/observability"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the DupLynx demo web server",
		RunE:  runServe,
	}

	return cmd
}

func runServe(cmd *cobra.Command, _ []string) (err error) {
	ctx := cmd.Context()
	cfg, ok := config.FromContext(ctx)
	if !ok {
		cfg = runtimeCfg
	}

	actor := resolveActor()
	metadata := map[string]any{
		"addr":       cfg.Addr,
		"db_file":    cfg.DBFile,
		"assets_dir": cfg.AssetsDir,
		"pid":        os.Getpid(),
		"go_version": runtime.Version(),
	}

	writer := observability.NewEventWriter(nil)
	start := time.Now()
	writer.Write(observability.Event{
		Action:   "serve_start",
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
			Action:   "serve_stop",
			Actor:    actor,
			Outcome:  outcome,
			Duration: time.Since(start),
			Metadata: metadata,
			Error:    err,
		})
	}()

	if err = app.ValidateAssetDirectory(cfg.AssetsDir); err != nil {
		return err
	}

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

	if err = data.Migrate(ctx, client); err != nil {
		return err
	}

	tenancyRepo := tenancy.NewRepositoryFromClient(client, &tenancy.AuditLogger{})
	scanRepo := scans.NewRepositoryFromClient(client)

	server := app.NewHTTPServer(app.ServerOptions{
		Addr: cfg.Addr,
		Handler: apphttp.NewRouter(apphttp.Dependencies{
			TenancyRepo: tenancyRepo,
			ScanRepo:    scanRepo,
			StaticFS:    http.Dir(cfg.AssetsDir),
		}),
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case serveErr := <-errCh:
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			err = serveErr
		}
	case <-ctx.Done():
		shutdownErr := app.Shutdown(context.Background(), server, 15*time.Second)
		serveErr := <-errCh
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			err = serveErr
		} else if shutdownErr != nil && !errors.Is(shutdownErr, context.Canceled) {
			err = shutdownErr
		}
	}

	return err
}

func resolveActor() string {
	if v := os.Getenv("CI"); v != "" {
		return "ci"
	}
	if v := os.Getenv("GITHUB_ACTOR"); v != "" {
		return v
	}
	if v := os.Getenv("USER"); v != "" {
		return v
	}
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	if host, err := os.Hostname(); err == nil && host != "" {
		return host
	}
	return "unknown"
}
