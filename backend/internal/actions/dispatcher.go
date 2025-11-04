package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrGroupNotFound    = errors.New("duplicate group not found")
	ErrKeeperMachineID  = errors.New("keeper machine id required")
	ErrInvalidGroupID   = errors.New("invalid duplicate group id")
	ErrInvalidMachineID = errors.New("invalid machine id")
)

// Dispatcher coordinates keeper assignments and duplicate actions.
type Dispatcher struct {
	Repo  *Repository
	Audit *AuditLogger
}

// NewDispatcher constructs a dispatcher backed by the supplied repository and audit logger.
func NewDispatcher(repo *Repository, audit *AuditLogger) *Dispatcher {
	return &Dispatcher{Repo: repo, Audit: audit}
}

func (d *Dispatcher) repository() (*Repository, error) {
	if d == nil || d.Repo == nil {
		return nil, errors.New("actions repository not configured")
	}
	return d.Repo, nil
}

// AssignKeeper records a keeper machine selection for the duplicate group.
func (d *Dispatcher) AssignKeeper(ctx context.Context, groupID, tenantSlug, machineID string) error {
	if machineID == "" {
		return ErrKeeperMachineID
	}
	repo, err := d.repository()
	if err != nil {
		return err
	}

	gid, err := uuid.Parse(groupID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidGroupID, err)
	}
	mid, err := uuid.Parse(machineID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMachineID, err)
	}

	group, err := repo.Get(ctx, gid)
	if err != nil {
		return err
	}
	if group.TenantSlug != tenantSlug {
		return ErrGroupNotFound
	}

	if err := repo.UpdateKeeper(ctx, gid, mid); err != nil {
		return err
	}

	if d.Audit != nil {
		d.Audit.Log(AuditEntry{
			Type:            "assign_keeper",
			GroupID:         groupID,
			TenantSlug:      tenantSlug,
			KeeperMachineID: machineID,
		})
	}
	return nil
}

// ActionType enumerates supported duplicate actions.
type ActionType string

const (
	ActionDelete     ActionType = "delete_copies"
	ActionHardlink   ActionType = "create_hardlinks"
	ActionQuarantine ActionType = "quarantine"
)

// PerformAction executes the duplicate action within the tenant scope.
func (d *Dispatcher) PerformAction(ctx context.Context, groupID, tenantSlug, actor string, action ActionType, payload map[string]any) error {
	repo, err := d.repository()
	if err != nil {
		return err
	}

	gid, err := uuid.Parse(groupID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidGroupID, err)
	}

	group, err := repo.Get(ctx, gid)
	if err != nil {
		return err
	}
	if group.TenantSlug != tenantSlug {
		return ErrGroupNotFound
	}

	if action == ActionQuarantine {
		if err := repo.QuarantineFiles(ctx, gid); err != nil {
			return err
		}
	}

	if d.Audit != nil {
		d.Audit.Log(AuditEntry{
			Type:       string(action),
			GroupID:    groupID,
			TenantSlug: tenantSlug,
			Actor:      actor,
			Payload:    payload,
			Stubbed:    true,
		})
	}
	return nil
}
