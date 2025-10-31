package contract_test

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

func setupActionsRouter() (*httptest.Server, *actions.AuditLogger) {
	audit := &actions.AuditLogger{}
	store := actions.NewStore(actions.SampleDuplicateGroups())
	dispatcher := actions.Dispatcher{Store: store, Audit: audit}
	repo := tenancy.NewRepository(tenancy.SampleTenants(), &tenancy.AuditLogger{})
	scanRepo := scans.NewRepository(scans.SampleScans())
	router := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       repo,
		ScanRepo:          scanRepo,
		ActionsStore:      store,
		ActionsDispatcher: dispatcher,
	})
	return httptest.NewServer(router), audit
}

func TestAssignKeeperContract(t *testing.T) {
	ts, audit := setupActionsRouter()
	defer ts.Close()

	payload := map[string]any{
		"tenantSlug":      "sample-tenant-a",
		"keeperMachineId": "helios-02",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/duplicate-groups/dg-001/keeper", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	entries := audit.Entries()
	if len(entries) == 0 || entries[0].Type != "assign_keeper" {
		t.Fatalf("expected assign_keeper audit entry, got %#v", entries)
	}
}

func TestActionEndpointContract(t *testing.T) {
	ts, audit := setupActionsRouter()
	defer ts.Close()

	payload := map[string]any{
		"tenantSlug":    "sample-tenant-a",
		"actionType":    "quarantine",
		"targetFileIds": []string{"f1"},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/duplicate-groups/dg-001/actions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}

	entries := audit.Entries()
	var found bool
	for _, entry := range entries {
		if entry.Type == "quarantine" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected quarantine audit entry, got %#v", entries)
	}
}
