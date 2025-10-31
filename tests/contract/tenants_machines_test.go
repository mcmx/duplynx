package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/http/handlers"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func setupTestRouter() *chi.Mux {
	audit := &tenancy.AuditLogger{}
	repo := tenancy.NewRepository(tenancy.SampleTenants(), audit)
	return apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo: repo,
	})
}

func TestListTenantsContract(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/tenants", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload struct {
		Tenants []handlers.TenantSummary `json:"tenants"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed parsing response: %v", err)
	}

	if len(payload.Tenants) == 0 {
		t.Fatalf("expected tenants in response")
	}
}

func TestListMachinesContract(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/tenants/sample-tenant-a/machines", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload struct {
		Machines []handlers.MachineSummary `json:"machines"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed parsing response: %v", err)
	}

	if len(payload.Machines) != 5 {
		t.Fatalf("expected 5 machines, got %d", len(payload.Machines))
	}
}
