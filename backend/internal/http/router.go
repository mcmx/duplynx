package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mcmx/duplynx/internal/http/handlers"
	appmiddleware "github.com/mcmx/duplynx/internal/http/middleware"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/templ"
	"github.com/mcmx/duplynx/internal/tenancy"
)

// Dependencies encapsulates services required by HTTP handlers.
type Dependencies struct {
	TenancyRepo *tenancy.Repository
	ScanService scans.Service
}

// NewRouter wires baseline routes and middleware; handlers attach in feature phases.
func NewRouter(deps Dependencies) *chi.Mux {
	r := chi.NewRouter()
	r.Use(appmiddleware.Instrumentation)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	if deps.TenancyRepo != nil {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			tenants, err := deps.TenancyRepo.ListTenants(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			body := templ.LaunchPage(tenants)
			markup := templ.RenderLayout("DupLynx", "", "", body)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(markup))
		})

		tenantsHandler := handlers.TenantsHandler{Repo: deps.TenancyRepo}
		machinesHandler := handlers.MachinesHandler{Repo: deps.TenancyRepo}
		scanListHandler := handlers.ScanListHandler{Service: deps.ScanService}
		scanBoardHandler := handlers.ScanBoardHandler{Service: deps.ScanService}

		r.Get("/tenants", tenantsHandler.ServeHTTP)
		r.Get("/tenants/{tenantSlug}/machines", machinesHandler.ServeHTTP)
		r.Get("/tenants/{tenantSlug}/scans", scanListHandler.ServeHTTP)
		r.Get("/scans/{scanID}", scanBoardHandler.ServeHTTP)
	}

	return r
}
