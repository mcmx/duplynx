package integration_test

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

func TestActionLoggingPipeline(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	repo := tenancy.NewRepositoryFromClient(seed.Client, &tenancy.AuditLogger{})
	scanRepo := scans.NewRepositoryFromClient(seed.Client)
	actionsRepo := actions.NewRepositoryFromClient(seed.Client)
	audit := &actions.AuditLogger{}
	dispatcher := actions.NewDispatcher(actionsRepo, audit)

	group := seed.Dataset.DuplicateGroups[0]
	groupID := group.ID.String()
	tenantSlug := testutil.TenantSlugFor(t, seed.Dataset, group.TenantID)
	machines := testutil.MachineIDsForTenant(seed.Dataset, group.TenantID)
	if len(machines) < 2 {
		t.Fatalf("expected at least two machines for tenant %s", tenantSlug)
	}

	r := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       repo,
		ScanRepo:          scanRepo,
		ActionsRepo:       actionsRepo,
		ActionsDispatcher: dispatcher,
	})

	reqBody, _ := json.Marshal(map[string]any{
		"tenantSlug":      tenantSlug,
		"keeperMachineId": machines[1].String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/duplicate-groups/"+groupID+"/keeper", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	entries := audit.Entries()
	if len(entries) == 0 || entries[0].Type != "assign_keeper" {
		t.Fatalf("expected assign_keeper audit entry, got %#v", entries)
	}
	updated, err := actionsRepo.Get(context.Background(), group.ID)
	if err != nil {
		t.Fatalf("reload duplicate group: %v", err)
	}
	if updated.KeeperMachineID != machines[1].String() {
		t.Fatalf("keeper not updated: %s", updated.KeeperMachineID)
	}
}

func TestActionEndpointLogging(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	repo := tenancy.NewRepositoryFromClient(seed.Client, &tenancy.AuditLogger{})
	scanRepo := scans.NewRepositoryFromClient(seed.Client)
	actionsRepo := actions.NewRepositoryFromClient(seed.Client)
	audit := &actions.AuditLogger{}
	dispatcher := actions.NewDispatcher(actionsRepo, audit)

	group := seed.Dataset.DuplicateGroups[0]
	groupID := group.ID.String()
	tenantSlug := testutil.TenantSlugFor(t, seed.Dataset, group.TenantID)

	r := apphttp.NewRouter(apphttp.Dependencies{
		TenancyRepo:       repo,
		ScanRepo:          scanRepo,
		ActionsRepo:       actionsRepo,
		ActionsDispatcher: dispatcher,
	})

	reqBody, _ := json.Marshal(map[string]any{
		"tenantSlug":    tenantSlug,
		"actionType":    "delete_copies",
		"targetFileIds": []string{},
	})
	req := httptest.NewRequest(http.MethodPost, "/duplicate-groups/"+groupID+"/actions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tenancy.HeaderTenantSlug, tenantSlug)
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
