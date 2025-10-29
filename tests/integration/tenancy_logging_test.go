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

func TestMachineSelectionLogging(t *testing.T) {
	audit := &tenancy.AuditLogger{}
	repo := tenancy.NewRepository(tenancy.SampleTenants(), audit)
	r := http.NewRouter(http.Dependencies{TenancyRepo: repo})

	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-a/machines?selected_machine=ares-laptop", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	entries := audit.Entries()
	var machineLogged bool
	for _, entry := range entries {
		if entry.Type == "machine_selection" && entry.MachineID == "ares-laptop" {
			machineLogged = true
		}
	}

	if !machineLogged {
		t.Fatalf("expected machine selection audit entry, got %#v", entries)
	}
}
