package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

const envPrefix = "DUPLYNX"

// RuntimeConfig captures the shared configuration required by DupLynx CLI commands.
type RuntimeConfig struct {
	DBFile    string
	AssetsDir string
	Addr      string
	LogLevel  string
}

// DefaultRuntimeConfig returns the baseline configuration before flags or environment overrides.
func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		DBFile:    "var/duplynx.db",
		AssetsDir: "backend/web/dist",
		Addr:      "0.0.0.0:8080",
		LogLevel:  "info",
	}
}

// BindRuntimeFlags wires persistent CLI flags to the provided configuration struct.
func BindRuntimeFlags(flagSet *pflag.FlagSet, cfg *RuntimeConfig) {
	if cfg == nil || flagSet == nil {
		return
	}

	flagSet.StringVar(&cfg.DBFile, "db-file", cfg.DBFile, "Path to the SQLite database file")
	flagSet.StringVar(&cfg.AssetsDir, "assets-dir", cfg.AssetsDir, "Directory containing precompiled static assets")
	flagSet.StringVar(&cfg.Addr, "addr", cfg.Addr, "Address for the HTTP server to bind")
	flagSet.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Logging level for CLI output (debug, info, warn, error)")
}

// ApplyEnvOverrides updates configuration values with environment variables unless flags already set them.
func (cfg *RuntimeConfig) ApplyEnvOverrides(flagSet *pflag.FlagSet) {
	if cfg == nil {
		return
	}

	apply := func(flagName, envSuffix string, target *string) {
		if target == nil {
			return
		}
		if flagSet != nil && flagSet.Changed(flagName) {
			return
		}
		if val, ok := os.LookupEnv(envKey(envSuffix)); ok && strings.TrimSpace(val) != "" {
			*target = strings.TrimSpace(val)
		}
	}

	apply("db-file", "DB_FILE", &cfg.DBFile)
	apply("assets-dir", "ASSETS_DIR", &cfg.AssetsDir)
	apply("addr", "ADDR", &cfg.Addr)
	apply("log-level", "LOG_LEVEL", &cfg.LogLevel)
}

// SQLiteDSN constructs the SQLite DSN with sensible defaults.
func (cfg RuntimeConfig) SQLiteDSN() string {
	path := cfg.DBFile
	if strings.TrimSpace(path) == "" {
		path = DefaultRuntimeConfig().DBFile
	}
	path = filepath.ToSlash(path)
	return fmt.Sprintf("file:%s?_busy_timeout=5000&_foreign_keys=1", path)
}

type contextKey struct{}

// WithRuntimeConfig embeds the configuration into a context so subcommands can retrieve it.
func WithRuntimeConfig(ctx context.Context, cfg RuntimeConfig) context.Context {
	return context.WithValue(ctx, contextKey{}, cfg)
}

// FromContext extracts the runtime configuration from context, if present.
func FromContext(ctx context.Context) (RuntimeConfig, bool) {
	if ctx == nil {
		return RuntimeConfig{}, false
	}
	value := ctx.Value(contextKey{})
	if cfg, ok := value.(RuntimeConfig); ok {
		return cfg, true
	}
	return RuntimeConfig{}, false
}

func envKey(suffix string) string {
	return fmt.Sprintf("%s_%s", envPrefix, suffix)
}
