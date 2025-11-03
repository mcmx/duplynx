package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/ent"
	"github.com/mcmx/duplynx/ent/actionaudit"
	"github.com/mcmx/duplynx/ent/duplicategroup"
	"github.com/mcmx/duplynx/ent/machine"
	"github.com/mcmx/duplynx/ent/migrate"
)

// SeedReport captures aggregate counts after seeding completes.
type SeedReport struct {
	Tenants         int
	Machines        int
	Scans           int
	DuplicateGroups int
	FileInstances   int
	ActionAudits    int
}

// SeedDemoDataset resets the database and loads the canonical dataset.
func SeedDemoDataset(ctx context.Context, client *ent.Client, dataset DemoDataset) (SeedReport, error) {
	if client == nil {
		return SeedReport{}, errors.New("ent client is nil")
	}

	if dataset.Tenants == nil || len(dataset.Tenants) == 0 {
		dataset = CanonicalDemoDataset()
	} else {
		dataset = dataset.Clone()
	}

	if ctx == nil {
		ctx = context.Background()
	}

	migrateCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := client.Schema.Create(migrateCtx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return SeedReport{}, fmt.Errorf("apply migrations: %w", err)
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return SeedReport{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := clearExisting(ctx, tx); err != nil {
		return SeedReport{}, err
	}

	if err := insertTenants(ctx, tx, dataset.Tenants); err != nil {
		return SeedReport{}, err
	}
	if err := insertMachines(ctx, tx, dataset.Machines); err != nil {
		return SeedReport{}, err
	}
	if err := insertScans(ctx, tx, dataset.Scans); err != nil {
		return SeedReport{}, err
	}
	if err := insertDuplicateGroups(ctx, tx, dataset.DuplicateGroups); err != nil {
		return SeedReport{}, err
	}
	if err := insertFileInstances(ctx, tx, dataset.FileInstances); err != nil {
		return SeedReport{}, err
	}
	if err := insertActionAudits(ctx, tx, dataset.ActionAudits); err != nil {
		return SeedReport{}, err
	}

	if err := tx.Commit(); err != nil {
		return SeedReport{}, fmt.Errorf("commit transaction: %w", err)
	}

	return SeedReport{
		Tenants:         len(dataset.Tenants),
		Machines:        len(dataset.Machines),
		Scans:           len(dataset.Scans),
		DuplicateGroups: len(dataset.DuplicateGroups),
		FileInstances:   len(dataset.FileInstances),
		ActionAudits:    len(dataset.ActionAudits),
	}, nil
}

func clearExisting(ctx context.Context, tx *ent.Tx) error {
	if _, err := tx.ActionAudit.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("clear action audits: %w", err)
	}
	if _, err := tx.FileInstance.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("clear file instances: %w", err)
	}
	if _, err := tx.DuplicateGroup.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("clear duplicate groups: %w", err)
	}
	if _, err := tx.Scan.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("clear scans: %w", err)
	}
	if _, err := tx.Machine.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("clear machines: %w", err)
	}
	if _, err := tx.Tenant.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("clear tenants: %w", err)
	}
	return nil
}

func insertTenants(ctx context.Context, tx *ent.Tx, tenants []TenantFixture) error {
	for _, tenantFixture := range tenants {
		builder := tx.Tenant.Create().
			SetID(tenantFixture.ID).
			SetSlug(tenantFixture.Slug).
			SetName(tenantFixture.Name).
			SetCreateTime(tenantFixture.CreatedAt).
			SetUpdateTime(tenantFixture.UpdatedAt)

		if tenantFixture.Description != "" {
			builder.SetDescription(tenantFixture.Description)
		}
		if tenantFixture.PrimaryContact != "" {
			builder.SetPrimaryContact(tenantFixture.PrimaryContact)
		}

		if err := builder.Exec(ctx); err != nil {
			return fmt.Errorf("insert tenant %q: %w", tenantFixture.Slug, err)
		}
	}
	return nil
}

func insertMachines(ctx context.Context, tx *ent.Tx, machines []MachineFixture) error {
	for _, machineFixture := range machines {
		builder := tx.Machine.Create().
			SetID(machineFixture.ID).
			SetTenantID(machineFixture.TenantID).
			SetName(machineFixture.Name).
			SetCategory(machine.Category(machineFixture.Category)).
			SetCreateTime(machineFixture.CreatedAt).
			SetUpdateTime(machineFixture.UpdatedAt)

		if machineFixture.Hostname != "" {
			builder.SetHostname(machineFixture.Hostname)
		}
		if machineFixture.Role != "" {
			builder.SetRole(machineFixture.Role)
		}
		if !machineFixture.LastScan.IsZero() {
			builder.SetLastScanAt(machineFixture.LastScan)
		}

		if err := builder.Exec(ctx); err != nil {
			return fmt.Errorf("insert machine %q: %w", machineFixture.Name, err)
		}
	}
	return nil
}

func insertScans(ctx context.Context, tx *ent.Tx, scans []ScanFixture) error {
	for _, scanFixture := range scans {
		builder := tx.Scan.Create().
			SetID(scanFixture.ID).
			SetTenantID(scanFixture.TenantID).
			SetName(scanFixture.Name).
			SetStartedAt(scanFixture.StartedAt).
			SetDuplicateGroupCount(scanFixture.DuplicateGroupCount).
			SetCreateTime(scanFixture.CreatedAt).
			SetUpdateTime(scanFixture.UpdatedAt)

		if scanFixture.Description != "" {
			builder.SetDescription(scanFixture.Description)
		}
		if !scanFixture.CompletedAt.IsZero() {
			builder.SetCompletedAt(scanFixture.CompletedAt)
		}
		if scanFixture.InitiatedMachineID != uuid.Nil {
			builder.SetInitiatedMachineID(scanFixture.InitiatedMachineID)
		}

		if err := builder.Exec(ctx); err != nil {
			return fmt.Errorf("insert scan %q: %w", scanFixture.Name, err)
		}
	}
	return nil
}

func insertDuplicateGroups(ctx context.Context, tx *ent.Tx, groups []DuplicateGroupFixture) error {
	for _, group := range groups {
		builder := tx.DuplicateGroup.Create().
			SetID(group.ID).
			SetTenantID(group.TenantID).
			SetScanID(group.ScanID).
			SetHash(group.Hash).
			SetStatus(duplicategroup.Status(group.Status)).
			SetFileCount(group.FileCount).
			SetTotalSizeBytes(group.TotalSizeBytes).
			SetCreateTime(group.CreatedAt).
			SetUpdateTime(group.UpdatedAt)

		if group.KeeperMachineID != uuid.Nil {
			builder.SetKeeperMachineID(group.KeeperMachineID)
		}

		if err := builder.Exec(ctx); err != nil {
			return fmt.Errorf("insert duplicate group %s: %w", group.ID, err)
		}
	}
	return nil
}

func insertFileInstances(ctx context.Context, tx *ent.Tx, files []FileInstanceFixture) error {
	for _, file := range files {
		builder := tx.FileInstance.Create().
			SetID(file.ID).
			SetDuplicateGroupID(file.DuplicateGroupID).
			SetMachineID(file.MachineID).
			SetPath(file.Path).
			SetSizeBytes(file.SizeBytes).
			SetChecksum(file.Checksum).
			SetLastSeenAt(file.LastSeenAt).
			SetQuarantined(file.Quarantined).
			SetCreateTime(file.CreatedAt).
			SetUpdateTime(file.UpdatedAt)

		if err := builder.Exec(ctx); err != nil {
			return fmt.Errorf("insert file instance %s: %w", file.ID, err)
		}
	}
	return nil
}

func insertActionAudits(ctx context.Context, tx *ent.Tx, audits []ActionAuditFixture) error {
	for _, audit := range audits {
		builder := tx.ActionAudit.Create().
			SetID(audit.ID).
			SetTenantID(audit.TenantID).
			SetDuplicateGroupID(audit.DuplicateGroupID).
			SetActor(audit.Actor).
			SetActionType(actionaudit.ActionType(audit.ActionType)).
			SetPerformedAt(audit.PerformedAt).
			SetStubbed(audit.Stubbed).
			SetCreateTime(audit.CreatedAt).
			SetUpdateTime(audit.UpdatedAt)

		if audit.Payload != nil {
			builder.SetPayload(audit.Payload)
		}

		if err := builder.Exec(ctx); err != nil {
			return fmt.Errorf("insert action audit %s: %w", audit.ID, err)
		}
	}
	return nil
}
