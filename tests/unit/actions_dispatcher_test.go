package unit_test

import (
	"testing"

	"github.com/mcmx/duplynx/internal/actions"
)

func TestAssignKeeperUpdatesGroupAndLogs(t *testing.T) {
	store := actions.NewStore(actions.SampleDuplicateGroups())
	audit := &actions.AuditLogger{}
	d := actions.Dispatcher{Store: store, Audit: audit}

	if err := d.AssignKeeper("dg-001", "sample-tenant-a", "atlas-01"); err != nil {
		t.Fatalf("assign keeper failed: %v", err)
	}
	group, _ := store.Get("dg-001")
	if group.KeeperMachineID != "atlas-01" {
		t.Fatalf("keeper not updated: %s", group.KeeperMachineID)
	}
	entries := audit.Entries()
	if len(entries) == 0 || entries[0].Type != "assign_keeper" {
		t.Fatalf("expected assign_keeper audit entry, got %#v", entries)
	}
}

func TestPerformActionCreatesStubbedAudit(t *testing.T) {
	store := actions.NewStore(actions.SampleDuplicateGroups())
	audit := &actions.AuditLogger{}
	d := actions.Dispatcher{Store: store, Audit: audit}

	if err := d.PerformAction("dg-001", "sample-tenant-a", "system", actions.ActionQuarantine, map[string]any{"targetFileIds": []string{"f1"}}); err != nil {
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

	group, _ := store.Get("dg-001")
	if !group.Files[0].Quarantined {
		t.Fatalf("expected file to be marked quarantined")
	}
}
