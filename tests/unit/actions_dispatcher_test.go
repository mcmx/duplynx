package unit_test

import (
	"context"
	"testing"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/tests/testutil"
)

func TestAssignKeeperUpdatesGroupAndLogs(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	repo := actions.NewRepositoryFromClient(seed.Client)
	audit := &actions.AuditLogger{}
	d := actions.NewDispatcher(repo, audit)

	groupFixture := seed.Dataset.DuplicateGroups[0]
	groupID := groupFixture.ID.String()
	tenantSlug := testutil.TenantSlugFor(t, seed.Dataset, groupFixture.TenantID)
	machineIDs := testutil.MachineIDsForTenant(seed.Dataset, groupFixture.TenantID)
	if len(machineIDs) == 0 {
		t.Fatal("expected machines for tenant in dataset")
	}
	machineID := machineIDs[0].String()

	if err := d.AssignKeeper(context.Background(), groupID, tenantSlug, machineID); err != nil {
		t.Fatalf("assign keeper failed: %v", err)
	}

	updated, err := repo.Get(context.Background(), groupFixture.ID)
	if err != nil {
		t.Fatalf("reload duplicate group: %v", err)
	}
	if updated.KeeperMachineID != machineID {
		t.Fatalf("keeper not updated: %s", updated.KeeperMachineID)
	}
	entries := audit.Entries()
	if len(entries) == 0 || entries[0].Type != "assign_keeper" {
		t.Fatalf("expected assign_keeper audit entry, got %#v", entries)
	}
}

func TestPerformActionCreatesStubbedAudit(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	repo := actions.NewRepositoryFromClient(seed.Client)
	audit := &actions.AuditLogger{}
	d := actions.NewDispatcher(repo, audit)

	groupFixture := seed.Dataset.DuplicateGroups[0]
	groupID := groupFixture.ID.String()
	tenantSlug := testutil.TenantSlugFor(t, seed.Dataset, groupFixture.TenantID)

	if err := d.PerformAction(context.Background(), groupID, tenantSlug, "system", actions.ActionQuarantine, map[string]any{"targetFileIds": []string{}}); err != nil {
		t.Fatalf("perform action failed: %v", err)
	}

	entries := audit.Entries()
	var found bool
	for _, entry := range entries {
		if entry.Type == string(actions.ActionQuarantine) && entry.Stubbed {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected stubbed quarantine audit entry")
	}

	updated, err := repo.Get(context.Background(), groupFixture.ID)
	if err != nil {
		t.Fatalf("reload duplicate group: %v", err)
	}
	for _, file := range updated.Files {
		if !file.Quarantined {
			t.Fatalf("expected file %s to be quarantined", file.ID)
		}
	}
}
