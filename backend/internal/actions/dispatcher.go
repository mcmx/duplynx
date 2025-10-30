package actions

import (
	"errors"
)

var (
	ErrGroupNotFound   = errors.New("duplicate group not found")
	ErrKeeperMachineID = errors.New("keeper machine id required")
)

// Dispatcher coordinates keeper assignments and duplicate actions.
type Dispatcher struct {
	Store *Store
	Audit *AuditLogger
}

func (d Dispatcher) AssignKeeper(groupID, tenantSlug, machineID string) error {
	if machineID == "" {
		return ErrKeeperMachineID
	}
	group, ok := d.Store.Get(groupID)
	if !ok {
		return ErrGroupNotFound
	}
	group.KeeperMachineID = machineID
	d.Store.Update(group)
	d.Audit.Log(AuditEntry{
		Type:            "assign_keeper",
		GroupID:         groupID,
		TenantSlug:      tenantSlug,
		KeeperMachineID: machineID,
	})
	return nil
}

// ActionType enumerates supported duplicate actions.
type ActionType string

const (
	ActionDelete     ActionType = "delete_copies"
	ActionHardlink   ActionType = "create_hardlinks"
	ActionQuarantine ActionType = "quarantine"
)

func (d Dispatcher) PerformAction(groupID, tenantSlug, actor string, action ActionType, payload map[string]any) error {
	group, ok := d.Store.Get(groupID)
	if !ok {
		return ErrGroupNotFound
	}
	// Phase 1 implementation stubs remote execution, only logs audit entry.
	d.Audit.Log(AuditEntry{
		Type:       string(action),
		GroupID:    groupID,
		TenantSlug: tenantSlug,
		Actor:      actor,
		Payload:    payload,
		Stubbed:    true,
	})
	if action == ActionQuarantine {
		for i := range group.Files {
			group.Files[i].Quarantined = true
		}
		d.Store.Update(group)
	}
	return nil
}
