package integration_test

import (
	"testing"

	"github.com/mcmx/duplynx/internal/actions"
)

func TestKeeperAssignmentDoesNotChangeStatus(t *testing.T) {
	store := actions.NewStore(actions.SampleDuplicateGroups())
	dispatcher := actions.Dispatcher{Store: store, Audit: &actions.AuditLogger{}}

	groupID := "dg-001"
	groupBefore, ok := store.Get(groupID)
	if !ok {
		t.Fatalf("seed group missing")
	}
	statusBefore := groupBefore.Status

	if err := dispatcher.AssignKeeper(groupID, groupBefore.TenantSlug, "helios-02"); err != nil {
		t.Fatalf("assign keeper failed: %v", err)
	}

	groupAfter, _ := store.Get(groupID)
	if groupAfter.Status != statusBefore {
		t.Fatalf("status changed unexpectedly: %s -> %s", statusBefore, groupAfter.Status)
	}
}
