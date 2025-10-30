package unit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/mcmx/duplynx/internal/tenancy"
)

func TestRequireTenantScopeInjectsContext(t *testing.T) {
	repo := tenancy.NewRepository(tenancy.SampleTenants(), nil)
	middleware := tenancy.RequireTenantScope(repo, nil)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scope, ok := tenancy.ScopeFromContext(r.Context())
		if !ok {
			t.Fatalf("expected tenant scope in context")
		}
		if scope.TenantSlug != "sample-tenant-a" {
			t.Fatalf("unexpected tenant slug: %s", scope.TenantSlug)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-a/machines", nil)
	req.Header.Set("X-Duplynx-Tenant", "sample-tenant-a")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("tenantSlug", "sample-tenant-a")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestRequireTenantScopeBlocksMismatchedPath(t *testing.T) {
	tenants := tenancy.SampleTenants()
	repo := tenancy.NewRepository(tenants, nil)
	middleware := tenancy.RequireTenantScope(repo, nil)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-b/scans", nil)
	req.Header.Set("X-Duplynx-Tenant", "sample-tenant-a")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("tenantSlug", "sample-tenant-b")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for scope violation, got %d", rr.Code)
	}
}
