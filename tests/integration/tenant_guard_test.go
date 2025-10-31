package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcmx/duplynx/internal/actions"
	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func TestTenantGuardBlocksCrossTenantAccess(t *testing.T) {
	tenants := tenancy.SampleTenants()
	tenants = append(tenants, tenancy.Tenant{
		Slug: "sample-tenant-b",
		Name: "Sample Tenant B",
		Machines: []tenancy.Machine{
			{ID: "b-laptop", TenantSlug: "sample-tenant-b", Name: "B-Laptop", Category: "personal_laptop", Hostname: "b.local"},
		},
	})
	repo := tenancy.NewRepository(tenants, nil)

	scanSummaries := append(scans.SampleScans(), scans.ScanSummary{
		ID:                  "tenant-b-scan",
		TenantSlug:          "sample-tenant-b",
		Name:                "Tenant B Baseline",
		InitiatedMachineID:  "b-laptop",
		StartedAt:           time.Now(),
		CompletedAt:         time.Now(),
		DuplicateGroupCount: 0,
		StatusCounts:        map[string]int{},
	})
	scanRepo := scans.NewRepository(scanSummaries)

	store := actions.NewStore(actions.SampleDuplicateGroups())
	dispatcher := actions.Dispatcher{Store: store, Audit: &actions.AuditLogger{}}

	router := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       repo,
		ScanRepo:          scanRepo,
		ActionsStore:      store,
		ActionsDispatcher: dispatcher,
	})

	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-b/scans", nil)
	req.Header.Set("X-Duplynx-Tenant", "sample-tenant-a")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for cross-tenant scan list access, got %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/scans/tenant-b-scan", nil)
	req.Header.Set("X-Duplynx-Tenant", "sample-tenant-a")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for cross-tenant scan view, got %d", rr.Code)
	}
}
