package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/tenancy"
	"github.com/mcmx/duplynx/tests/testutil"
)

func TestMachineSelectionLogging(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	audit := &tenancy.AuditLogger{}
	repo := tenancy.NewRepositoryFromClient(seed.Client, audit)
	router := apphttp.NewRouter(apphttp.Dependencies{TenancyRepo: repo})

	tenantSlug := seed.Dataset.Tenants[0].Slug
	machines := testutil.MachineIDsForTenant(seed.Dataset, seed.Dataset.Tenants[0].ID)
	if len(machines) == 0 {
		t.Fatalf("expected machines for tenant %s", tenantSlug)
	}
	selected := machines[0].String()

	req := httptest.NewRequest(http.MethodGet, "/tenants/"+tenantSlug+"/machines?selected_machine="+selected, nil)
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	entries := audit.Entries()
	var machineLogged bool
	for _, entry := range entries {
		if entry.Type == "machine_selection" && entry.MachineID == selected {
			machineLogged = true
			break
		}
	}

	if !machineLogged {
		t.Fatalf("expected machine selection audit entry, got %#v", entries)
	}
}
