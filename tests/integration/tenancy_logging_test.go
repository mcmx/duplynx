package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/tenancy"
	"github.com/mcmx/duplynx/tests/testutil"
)

func TestTenantSelectionLogging(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	audit := &tenancy.AuditLogger{}
	repo := tenancy.NewRepositoryFromClient(seed.Client, audit)
	r := apphttp.NewRouter(apphttp.Dependencies{TenancyRepo: repo})

	tenantSlug := seed.Dataset.Tenants[0].Slug
	req := httptest.NewRequest(http.MethodGet, "/tenants/"+tenantSlug+"/machines", nil)
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	entries := audit.Entries()
	if len(entries) == 0 || entries[0].Type != "tenant_selection" {
		t.Fatalf("expected tenant selection audit entry, got %#v", entries)
	}
}
