package app

import (
	"context"
	"net/http"
	"time"
)

// ServerOptions configure the HTTP server used by DupLynx CLI commands.
type ServerOptions struct {
	Addr            string
	Handler         http.Handler
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// NewHTTPServer constructs an *http.Server with sensible defaults for the DupLynx dashboard.
func NewHTTPServer(opts ServerOptions) *http.Server {
	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = 15 * time.Second
	}
	if opts.WriteTimeout == 0 {
		opts.WriteTimeout = 30 * time.Second
	}
	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = 60 * time.Second
	}

	return &http.Server{
		Addr:         opts.Addr,
		Handler:      opts.Handler,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}
}

// Shutdown gracefully stops the HTTP server within the configured timeout.
func Shutdown(ctx context.Context, srv *http.Server, timeout time.Duration) error {
	if srv == nil {
		return nil
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
