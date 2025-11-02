package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mcmx/duplynx/ent"
)

// OpenSQLite opens (and creates if necessary) a SQLite-backed Ent client using the provided DSN.
func OpenSQLite(ctx context.Context, dsn string) (*ent.Client, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("sqlite dsn must not be empty")
	}

	if err := ensureDirectory(dsn); err != nil {
		return nil, err
	}

	client, err := ent.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	return client, nil
}

// Migrate runs Ent schema migrations for the supplied client.
func Migrate(ctx context.Context, client *ent.Client) error {
	if client == nil {
		return errors.New("ent client is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	return client.Schema.Create(ctx)
}

// Close releases database resources.
func Close(client *ent.Client) error {
	if client == nil {
		return nil
	}
	return client.Close()
}

func ensureDirectory(dsn string) error {
	// DSN is of format file:<path>?...
	withoutScheme := strings.TrimPrefix(dsn, "file:")
	path := strings.SplitN(withoutScheme, "?", 2)[0]
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
