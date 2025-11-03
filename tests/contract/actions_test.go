package contract_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcmx/duplynx/ent/enttest"
	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/internal/data"
	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
)

func setupActionsRouter(t *testing.T) (*httptest.Server, *actions.AuditLogger, data.DemoDataset) {
	t.Helper()

	audit := &actions.AuditLogger{}
	client := enttest.Open(t, "sqlite3", "file:actions-contract?mode=memory&cache=shared&_fk=1")

	dataset := data.CanonicalDemoDataset()
	if _, err := data.SeedDemoDataset(context.Background(), client, dataset); err != nil {
		t.Fatalf("failed to seed dataset: %v", err)
	}

	store := actions.NewStore(dataset.DuplicateGroupsForStore())
	dispatcher := actions.Dispatcher{Store: store, Audit: audit}
	repo := tenancy.NewRepositoryFromClient(client, &tenancy.AuditLogger{})
	scanRepo := scans.NewRepositoryFromClient(client)
	router := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       repo,
		ScanRepo:          scanRepo,
		ActionsStore:      store,
		ActionsDispatcher: dispatcher,
	})

	server := httptest.NewServer(router)
	t.Cleanup(func() {
		server.Close()
		_ = client.Close()
	})

	return server, audit, dataset
}

func TestAssignKeeperContract(t *testing.T) {
	ts, audit, dataset := setupActionsRouter(t)

	tenantSlug := dataset.Tenants[0].Slug
	groupID := dataset.DuplicateGroups[0].ID.String()
	keeperMachineID := dataset.Machines[1].ID.String()

	payload := map[string]any{
		"tenantSlug":      tenantSlug,
		"keeperMachineId": keeperMachineID,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/duplicate-groups/"+groupID+"/keeper", bytes.NewReader(body))
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
	ts, audit, dataset := setupActionsRouter(t)

	tenantSlug := dataset.Tenants[0].Slug
	groupID := dataset.DuplicateGroups[1].ID.String()

	var fileID string
	for _, file := range dataset.FileInstances {
		if file.DuplicateGroupID.String() == groupID {
			fileID = file.ID.String()
			break
		}
	}
	if fileID == "" {
		t.Fatalf("expected file instance for group %s", groupID)
	}

	payload := map[string]any{
		"tenantSlug":    tenantSlug,
		"actionType":    "quarantine",
		"targetFileIds": []string{fileID},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/duplicate-groups/"+groupID+"/actions", bytes.NewReader(body))
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
