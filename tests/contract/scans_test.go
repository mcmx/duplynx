package contract_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
	"github.com/mcmx/duplynx/tests/testutil"
)

type boardHarness struct {
	router *chi.Mux
	seed   testutil.SeededClient
}

func setupBoardRouter(t *testing.T) boardHarness {
	seed := testutil.NewSeededClient(t)
	audit := &tenancy.AuditLogger{}
	tenancyRepo := tenancy.NewRepositoryFromClient(seed.Client, audit)
	scanRepo := scans.NewRepositoryFromClient(seed.Client)

	router := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo: tenancyRepo,
		ScanRepo:    scanRepo,
	})

	return boardHarness{router: router, seed: seed}
}

func TestScanListContract(t *testing.T) {
	harness := setupBoardRouter(t)
	tenantSlug := harness.seed.Dataset.Tenants[0].Slug

	req := httptest.NewRequest(http.MethodGet, "/tenants/"+tenantSlug+"/scans", nil)
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)
	rec := httptest.NewRecorder()

	harness.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestScanBoardContract(t *testing.T) {
	harness := setupBoardRouter(t)
	scanID := harness.seed.Dataset.Scans[0].ID.String()
	tenantSlug := testutil.TenantSlugFor(t, harness.seed.Dataset, harness.seed.Dataset.Scans[0].TenantID)

	req := httptest.NewRequest(http.MethodGet, "/scans/"+scanID, nil)
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)
	rec := httptest.NewRecorder()

	harness.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
