package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
	"github.com/mcmx/duplynx/internal/actions"
)

func TestActionLoggingPipeline(t *testing.T) {
	audit := &actions.AuditLogger{}
	store := actions.NewStore(actions.SampleDuplicateGroups())
	dispatcher := actions.Dispatcher{Store: store, Audit: audit}
	repo := tenancy.NewRepository(tenancy.SampleTenants(), &tenancy.AuditLogger{})
	scanRepo := scans.NewRepository(scans.SampleScans())

	r := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       repo,
		ScanRepo:          scanRepo,
		ActionsStore:      store,
		ActionsDispatcher: dispatcher,
	})

	reqBody, _ := json.Marshal(map[string]any{
		"tenantSlug":      "sample-tenant-a",
		"keeperMachineId": "helios-02",
	})
	req := httptest.NewRequest(http.MethodPost, "/duplicate-groups/dg-001/keeper", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	entries := audit.Entries()
	if len(entries) == 0 || entries[0].Type != "assign_keeper" {
		t.Fatalf("expected assign_keeper audit entry, got %#v", entries)
	}
}

func TestActionEndpointLogging(t *testing.T) {
	audit := &actions.AuditLogger{}
	store := actions.NewStore(actions.SampleDuplicateGroups())
	dispatcher := actions.Dispatcher{Store: store, Audit: audit}
	repo := tenancy.NewRepository(tenancy.SampleTenants(), &tenancy.AuditLogger{})
	scanRepo := scans.NewRepository(scans.SampleScans())

	r := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       repo,
		ScanRepo:          scanRepo,
		ActionsStore:      store,
		ActionsDispatcher: dispatcher,
	})

	reqBody, _ := json.Marshal(map[string]any{
		"tenantSlug":    "sample-tenant-a",
		"actionType":    "delete_copies",
		"targetFileIds": []string{"f1"},
	})
	req := httptest.NewRequest(http.MethodPost, "/duplicate-groups/dg-001/actions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	entries := audit.Entries()
	var found bool
	for _, entry := range entries {
		if entry.Type == "delete_copies" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected delete_copies entry, got %#v", entries)
	}
}
