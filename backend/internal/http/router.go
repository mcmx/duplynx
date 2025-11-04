package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/internal/http/handlers"
	appmiddleware "github.com/mcmx/duplynx/internal/http/middleware"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/templ"
	templerrors "github.com/mcmx/duplynx/internal/templ/errors"
	"github.com/mcmx/duplynx/internal/tenancy"
)

// Dependencies encapsulates services required by HTTP handlers.
type Dependencies struct {
	TenancyRepo       *tenancy.Repository
	ScanRepo          *scans.Repository
	ActionsRepo       *actions.Repository
	ActionsDispatcher *actions.Dispatcher
	StaticFS          http.FileSystem
}

// NewRouter wires baseline routes and middleware; handlers attach in feature phases.
func NewRouter(deps Dependencies) *chi.Mux {
	r := chi.NewRouter()
	r.Use(appmiddleware.Instrumentation)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	staticFS := deps.StaticFS
	if staticFS == nil {
		staticFS = http.Dir("web/static")
	}
	r.Handle("/static/*", handlers.StaticHandler{Root: staticFS})

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
		scopeRenderer := templerrors.Renderer{}
		scopeMiddleware := tenancy.RequireTenantScope(deps.TenancyRepo, scopeRenderer)

		r.Get("/tenants", tenantsHandler.ServeHTTP)
		r.With(scopeMiddleware).Get("/tenants/{tenantSlug}/machines", machinesHandler.ServeHTTP)

		if deps.ScanRepo != nil {
			service := scans.Service{Repo: deps.ScanRepo}
			scanListHandler := handlers.ScanListHandler{Service: service}
			scanBoardHandler := handlers.ScanBoardHandler{Service: service}

			r.With(scopeMiddleware).Get("/tenants/{tenantSlug}/scans", scanListHandler.ServeHTTP)
			r.With(scopeMiddleware).Get("/scans/{scanID}", scanBoardHandler.ServeHTTP)
		}

		if deps.ActionsDispatcher != nil && deps.ActionsRepo != nil {
			keeperHandler := handlers.KeeperHandler{Dispatcher: deps.ActionsDispatcher}
			actionHandler := handlers.ActionHandler{Dispatcher: deps.ActionsDispatcher}

			r.With(scopeMiddleware).Post("/duplicate-groups/{groupId}/keeper", keeperHandler.ServeHTTP)
			r.With(scopeMiddleware).Post("/duplicate-groups/{groupId}/actions", actionHandler.ServeHTTP)
			r.With(scopeMiddleware).Post("/duplicate-groups/{groupId}/htmx", handlers.ActionHTMXHandler)
		}
	}

	return r
}
