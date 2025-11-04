package integration_test

import (
	"context"
	"testing"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/tests/testutil"
)

func TestKeeperAssignmentDoesNotChangeStatus(t *testing.T) {
	seed := testutil.NewSeededClient(t)
	repo := actions.NewRepositoryFromClient(seed.Client)
	dispatcher := actions.NewDispatcher(repo, &actions.AuditLogger{})

	groupFixture := seed.Dataset.DuplicateGroups[0]
	groupBefore, err := repo.Get(context.Background(), groupFixture.ID)
	if err != nil {
		t.Fatalf("load duplicate group: %v", err)
	}
	statusBefore := groupBefore.Status

	machines := testutil.MachineIDsForTenant(seed.Dataset, groupFixture.TenantID)
	if len(machines) == 0 {
		t.Fatal("expected machines for tenant")
	}

	if err := dispatcher.AssignKeeper(context.Background(), groupBefore.ID, groupBefore.TenantSlug, machines[0].String()); err != nil {
		t.Fatalf("assign keeper failed: %v", err)
	}

	groupAfter, err := repo.Get(context.Background(), groupFixture.ID)
	if err != nil {
		t.Fatalf("reload duplicate group: %v", err)
	}
	if groupAfter.Status != statusBefore {
		t.Fatalf("status changed unexpectedly: %s -> %s", statusBefore, groupAfter.Status)
	}
}
