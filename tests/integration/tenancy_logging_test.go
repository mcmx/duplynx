package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func TestTenantSelectionLogging(t *testing.T) {
	audit := &tenancy.AuditLogger{}
	repo := tenancy.NewRepository(tenancy.SampleTenants(), audit)
	r := http.NewRouter(http.Dependencies{TenancyRepo: repo})

	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-a/machines", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	entries := audit.Entries()
	if len(entries) == 0 || entries[0].Type != "tenant_selection" {
		t.Fatalf("expected tenant selection audit entry, got %#v", entries)
	}
}
