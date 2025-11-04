package contract_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcmx/duplynx/internal/actions"
	apphttp "github.com/mcmx/duplynx/internal/http"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
	"github.com/mcmx/duplynx/tests/testutil"
)

type actionsHarness struct {
	server  *httptest.Server
	audit   *actions.AuditLogger
	dataset testutil.SeededClient
	repo    *actions.Repository
}

func setupActionsRouter(t *testing.T) actionsHarness {
	t.Helper()

	seed := testutil.NewSeededClient(t)
	audit := &actions.AuditLogger{}
	actionsRepo := actions.NewRepositoryFromClient(seed.Client)
	dispatcher := actions.NewDispatcher(actionsRepo, audit)

	router := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       tenancy.NewRepositoryFromClient(seed.Client, &tenancy.AuditLogger{}),
		ScanRepo:          scans.NewRepositoryFromClient(seed.Client),
		ActionsRepo:       actionsRepo,
		ActionsDispatcher: dispatcher,
	})

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	return actionsHarness{
		server:  server,
		audit:   audit,
		dataset: seed,
		repo:    actionsRepo,
	}
}

func TestAssignKeeperContract(t *testing.T) {
	harness := setupActionsRouter(t)

	group := harness.dataset.Dataset.DuplicateGroups[0]
	groupID := group.ID.String()
	tenantSlug := testutil.TenantSlugFor(t, harness.dataset.Dataset, group.TenantID)
	machines := testutil.MachineIDsForTenant(harness.dataset.Dataset, group.TenantID)
	if len(machines) < 2 {
		t.Fatalf("expected at least two machines for tenant %s", tenantSlug)
	}
	keeperMachineID := machines[1].String()

	payload := map[string]any{
		"tenantSlug":      tenantSlug,
		"keeperMachineId": keeperMachineID,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, harness.server.URL+"/duplicate-groups/"+groupID+"/keeper", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	if len(harness.audit.Entries()) == 0 || harness.audit.Entries()[0].Type != "assign_keeper" {
		t.Fatalf("expected assign_keeper audit entry, got %#v", harness.audit.Entries())
	}

	updated, err := harness.repo.Get(context.Background(), group.ID)
	if err != nil {
		t.Fatalf("reload duplicate group: %v", err)
	}
	if updated.KeeperMachineID != keeperMachineID {
		t.Fatalf("keeper machine not updated: %s", updated.KeeperMachineID)
	}
}

func TestActionEndpointContract(t *testing.T) {
	harness := setupActionsRouter(t)

	group := harness.dataset.Dataset.DuplicateGroups[1]
	groupID := group.ID.String()
	tenantSlug := testutil.TenantSlugFor(t, harness.dataset.Dataset, group.TenantID)

	payload := map[string]any{
		"tenantSlug":    tenantSlug,
		"actionType":    string(actions.ActionQuarantine),
		"targetFileIds": []string{},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, harness.server.URL+"/duplicate-groups/"+groupID+"/actions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}

	var found bool
	for _, entry := range harness.audit.Entries() {
		if entry.Type == string(actions.ActionQuarantine) {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected quarantine audit entry, got %#v", harness.audit.Entries())
	}

	updated, err := harness.repo.Get(context.Background(), group.ID)
	if err != nil {
		t.Fatalf("reload duplicate group: %v", err)
	}
	for _, file := range updated.Files {
		if !file.Quarantined {
			t.Fatalf("expected file %s to be quarantined", file.ID)
		}
	}
}
