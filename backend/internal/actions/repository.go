package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/ent"
	entduplicategroup "github.com/mcmx/duplynx/ent/duplicategroup"
	entfileinstance "github.com/mcmx/duplynx/ent/fileinstance"
)

// Repository surfaces duplicate group data backed by Ent.
type Repository struct {
	client *ent.Client
}

// NewRepositoryFromClient constructs a repository using the supplied Ent client.
func NewRepositoryFromClient(client *ent.Client) *Repository {
	if client == nil {
		return nil
	}
	return &Repository{client: client}
}

// ListByScan returns duplicate groups for a scan ordered by status and hash.
func (r *Repository) ListByScan(ctx context.Context, scanID uuid.UUID) ([]DuplicateGroup, error) {
	if r == nil || r.client == nil {
		return nil, errors.New("actions repository not configured")
	}
	records, err := r.client.DuplicateGroup.
		Query().
		Where(entduplicategroup.ScanID(scanID)).
		WithTenant().
		WithFileInstances(func(q *ent.FileInstanceQuery) {
			q.Order(entfileinstance.ByPath())
		}).
		Order(
			entduplicategroup.ByStatus(),
			entduplicategroup.ByHash(),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list duplicate groups: %w", err)
	}

	out := make([]DuplicateGroup, 0, len(records))
	for _, record := range records {
		out = append(out, convertDuplicateGroup(record))
	}
	return out, nil
}

// Get returns a single duplicate group by ID.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (DuplicateGroup, error) {
	if r == nil || r.client == nil {
		return DuplicateGroup{}, errors.New("actions repository not configured")
	}
	record, err := r.client.DuplicateGroup.
		Query().
		Where(entduplicategroup.IDEQ(id)).
		WithTenant().
		WithFileInstances(func(q *ent.FileInstanceQuery) {
			q.Order(entfileinstance.ByPath())
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return DuplicateGroup{}, ErrGroupNotFound
		}
		return DuplicateGroup{}, fmt.Errorf("load duplicate group: %w", err)
	}
	return convertDuplicateGroup(record), nil
}

// UpdateKeeper sets the keeper machine for a duplicate group.
func (r *Repository) UpdateKeeper(ctx context.Context, id uuid.UUID, machineID uuid.UUID) error {
	if r == nil || r.client == nil {
		return errors.New("actions repository not configured")
	}
	update := r.client.DuplicateGroup.UpdateOneID(id)
	if machineID == uuid.Nil {
		update = update.ClearKeeperMachineID()
	} else {
		update = update.SetKeeperMachineID(machineID)
	}
	if err := update.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return ErrGroupNotFound
		}
		return fmt.Errorf("update keeper machine: %w", err)
	}
	return nil
}

// QuarantineFiles marks all file instances for the duplicate group as quarantined.
func (r *Repository) QuarantineFiles(ctx context.Context, id uuid.UUID) error {
	if r == nil || r.client == nil {
		return errors.New("actions repository not configured")
	}
	_, err := r.client.FileInstance.
		Update().
		Where(entfileinstance.HasDuplicateGroupWith(entduplicategroup.IDEQ(id))).
		SetQuarantined(true).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("quarantine duplicate files: %w", err)
	}
	return nil
}

func convertDuplicateGroup(record *ent.DuplicateGroup) DuplicateGroup {
	if record == nil {
		return DuplicateGroup{}
	}
	var tenantSlug string
	if record.Edges.Tenant != nil {
		tenantSlug = record.Edges.Tenant.Slug
	}

	var keeperMachineID string
	if record.KeeperMachineID != uuid.Nil {
		keeperMachineID = record.KeeperMachineID.String()
	}

	files := make([]DuplicateFile, 0, len(record.Edges.FileInstances))
	for _, file := range record.Edges.FileInstances {
		files = append(files, DuplicateFile{
			ID:          file.ID.String(),
			MachineID:   file.MachineID.String(),
			Path:        file.Path,
			SizeBytes:   file.SizeBytes,
			Quarantined: file.Quarantined,
		})
	}

	return DuplicateGroup{
		ID:              record.ID.String(),
		ScanID:          record.ScanID.String(),
		TenantSlug:      tenantSlug,
		Status:          string(record.Status),
		KeeperMachineID: keeperMachineID,
		Hash:            record.Hash,
		Files:           files,
	}
}
