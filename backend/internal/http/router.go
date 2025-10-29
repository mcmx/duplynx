package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	appmiddleware "github.com/mcmx/duplynx/internal/http/middleware"
)

// NewRouter wires baseline routes and middleware; handlers attach in feature phases.
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(appmiddleware.Instrumentation)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return r
}
