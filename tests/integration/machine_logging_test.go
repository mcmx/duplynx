package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func TestMachineSelectionLogging(t *testing.T) {
	audit := &tenancy.AuditLogger{}
	repo := tenancy.NewRepository(tenancy.SampleTenants(), audit)
	router := http.NewRouter(http.Dependencies{TenancyRepo: repo})

	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-a/machines?selected_machine=ares-laptop", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	entries := audit.Entries()
	var machineLogged bool
	for _, entry := range entries {
		if entry.Type == "machine_selection" && entry.MachineID == "ares-laptop" {
			machineLogged = true
			break
		}
	}

	if !machineLogged {
		t.Fatalf("expected machine selection audit entry, got %#v", entries)
	}
}
