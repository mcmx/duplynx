package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/ent/enttest"
	"github.com/mcmx/duplynx/internal/data"
	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func TestTenantGuardBlocksCrossTenantAccess(t *testing.T) {
	ctx := context.Background()

	client := enttest.Open(t, "sqlite3", "file:tenant-guard?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() {
		_ = client.Close()
	})

	dataset := data.CanonicalDemoDataset()
	if _, err := data.SeedDemoDataset(ctx, client, dataset); err != nil {
		t.Fatalf("failed to seed dataset: %v", err)
	}

	repo := tenancy.NewRepositoryFromClient(client, &tenancy.AuditLogger{})
	scanRepo := scans.NewRepositoryFromClient(client)

	router := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo: repo,
		ScanRepo:    scanRepo,
	})

	var tenantA, tenantB string
	var tenantBID uuid.UUID
	for _, tenant := range dataset.Tenants {
		if tenant.Slug == "orion-analytics" {
			tenantA = tenant.Slug
		}
		if tenant.Slug == "selene-research" {
			tenantB = tenant.Slug
			tenantBID = tenant.ID
		}
	}
	if tenantA == "" || tenantB == "" || tenantBID == uuid.Nil {
		t.Fatalf("canonical dataset missing expected tenants")
	}

	var tenantBScanID string
	for _, scan := range dataset.Scans {
		if scan.TenantID == tenantBID {
			tenantBScanID = scan.ID.String()
			break
		}
	}
	if tenantBScanID == "" {
		t.Fatalf("canonical dataset missing scan for tenant %s", tenantB)
	}

	req := httptest.NewRequest(http.MethodGet, "/tenants/"+tenantB+"/scans", nil)
	req.Header.Set("X-Duplynx-Tenant", tenantA)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for cross-tenant scan list access, got %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/scans/"+tenantBScanID, nil)
	req.Header.Set("X-Duplynx-Tenant", tenantA)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for cross-tenant scan view, got %d", rr.Code)
	}
}
