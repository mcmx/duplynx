package tenancy

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

const HeaderTenantSlug = "X-Duplynx-Tenant"

type scopeKey struct{}

// Scope captures the tenant context for a request lifecycle.
type Scope struct {
	TenantSlug string
}

// ScopeViolation describes a mismatch between the active scope and a requested resource.
type ScopeViolation struct {
	CurrentTenant string
	RequestedSlug string
	Reason        string
}

// ScopeViolationRenderer renders tenant scope violation responses.
type ScopeViolationRenderer interface {
	RenderTenantScopeViolation(http.ResponseWriter, *http.Request, ScopeViolation)
}

// ScopeViolationRendererFunc adapts a function to ScopeViolationRenderer.
type ScopeViolationRendererFunc func(http.ResponseWriter, *http.Request, ScopeViolation)

// RenderTenantScopeViolation implements ScopeViolationRenderer.
func (f ScopeViolationRendererFunc) RenderTenantScopeViolation(w http.ResponseWriter, r *http.Request, violation ScopeViolation) {
	f(w, r, violation)
}

type scopeError struct {
	status    int
	message   string
	violation *ScopeViolation
}

func (e *scopeError) Error() string {
	return e.message
}

// RequireTenantScope ensures requests include a tenant slug and attaches it to the context.
func RequireTenantScope(repo *Repository, renderer ScopeViolationRenderer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scope, err := resolveScope(r, repo)
			if err != nil {
				if renderer != nil && err.violation != nil {
					renderer.RenderTenantScopeViolation(w, r, *err.violation)
					return
				}
				http.Error(w, err.message, err.status)
				return
			}

			w.Header().Set(HeaderTenantSlug, scope.TenantSlug)
			ctx := context.WithValue(r.Context(), scopeKey{}, scope)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ScopeFromContext retrieves the tenant scope from context.
func ScopeFromContext(ctx context.Context) (Scope, bool) {
	scope, ok := ctx.Value(scopeKey{}).(Scope)
	return scope, ok
}

func resolveScope(r *http.Request, repo *Repository) (Scope, *scopeError) {
	headerSlug := strings.TrimSpace(r.Header.Get(HeaderTenantSlug))
	pathSlug := strings.TrimSpace(chi.URLParam(r, "tenantSlug"))

	tenantSlug := headerSlug
	if tenantSlug == "" {
		tenantSlug = pathSlug
	}

	if tenantSlug == "" {
		return Scope{}, &scopeError{
			status:  http.StatusBadRequest,
			message: "tenant scope required",
		}
	}

	if repo != nil {
		if _, ok := repo.Tenant(tenantSlug); !ok {
			return Scope{}, &scopeError{
				status:  http.StatusNotFound,
				message: "tenant not found",
				violation: &ScopeViolation{
					CurrentTenant: headerSlug,
					RequestedSlug: tenantSlug,
					Reason:        "tenant not registered",
				},
			}
		}
	}

	if pathSlug != "" && pathSlug != tenantSlug {
		return Scope{}, &scopeError{
			status:  http.StatusNotFound,
			message: "tenant scope violation",
			violation: &ScopeViolation{
				CurrentTenant: tenantSlug,
				RequestedSlug: pathSlug,
				Reason:        "path tenant does not match active scope",
			},
		}
	}

	return Scope{TenantSlug: tenantSlug}, nil
}
