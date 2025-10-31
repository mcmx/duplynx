package contract_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func setupBoardRouter() *chi.Mux {
	audit := &tenancy.AuditLogger{}
	tenancyRepo := tenancy.NewRepository(tenancy.SampleTenants(), audit)
	scanRepo := scans.NewRepository(scans.SampleScans())
	svc := scans.Service{ScansRepo: scanRepo}
	return apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo: tenancyRepo,
		ScanService: svc,
	})
}

func TestScanListContract(t *testing.T) {
	r := setupBoardRouter()
	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-a/scans", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestScanBoardContract(t *testing.T) {
	r := setupBoardRouter()
	req := httptest.NewRequest(http.MethodGet, "/scans/baseline-sweep-2025-10-01", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
