package scans

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/ent"
	entscan "github.com/mcmx/duplynx/ent/scan"
	enttenant "github.com/mcmx/duplynx/ent/tenant"
)

var (
	ErrScanNotFound = errors.New("scan not found")
)

// Repository surfaces seeded scan summaries for the demo.
type Repository struct {
	scans  map[string]ScanSummary
	client *ent.Client
}

// NewRepository constructs a repository with immutable scan data.
func NewRepository(scans []ScanSummary) *Repository {
	index := make(map[string]ScanSummary, len(scans))
	for _, scan := range scans {
		index[scan.ID] = scan
	}
	return &Repository{scans: index}
}

// NewRepositoryFromClient constructs a repository backed by Ent.
func NewRepositoryFromClient(client *ent.Client) *Repository {
	return &Repository{client: client}
}

// ListByTenant returns scan summaries for a tenant.
func (r *Repository) ListByTenant(ctx context.Context, tenantSlug string) ([]ScanSummary, error) {
	if r.client != nil {
		return r.listByTenantFromClient(ctx, tenantSlug)
	}

	var out []ScanSummary
	for _, scan := range r.scans {
		if scan.TenantSlug == tenantSlug {
			out = append(out, scan)
		}
	}
	return out, nil
}

// Get returns a scan summary by ID.
func (r *Repository) Get(ctx context.Context, scanID string) (ScanSummary, error) {
	if r.client != nil {
		return r.getFromClient(ctx, scanID)
	}

	scan, ok := r.scans[scanID]
	if !ok {
		return ScanSummary{}, ErrScanNotFound
	}
	return scan, nil
}

func (r *Repository) listByTenantFromClient(ctx context.Context, tenantSlug string) ([]ScanSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	records, err := r.client.Scan.
		Query().
		Where(entscan.HasTenantWith(enttenant.SlugEQ(tenantSlug))).
		WithTenant().
		WithDuplicateGroups(func(q *ent.DuplicateGroupQuery) {
			q.WithKeeperMachine()
		}).
		Order(entscan.ByStartedAt()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ScanSummary, 0, len(records))
	for _, record := range records {
		out = append(out, convertScan(record))
	}
	return out, nil
}

func (r *Repository) getFromClient(ctx context.Context, scanID string) (ScanSummary, error) {
	u, err := uuid.Parse(scanID)
	if err != nil {
		return ScanSummary{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	record, err := r.client.Scan.
		Query().
		Where(entscan.IDEQ(u)).
		WithTenant().
		WithDuplicateGroups(func(q *ent.DuplicateGroupQuery) {
			q.WithKeeperMachine()
		}).
		Only(ctx)
	if err != nil {
		return ScanSummary{}, err
	}
	return convertScan(record), nil
}

func convertScan(record *ent.Scan) ScanSummary {
	if record == nil {
		return ScanSummary{}
	}

	summary := ScanSummary{
		ID:                  record.ID.String(),
		Name:                record.Name,
		InitiatedMachineID:  uuidToString(record.InitiatedMachineID),
		StartedAt:           record.StartedAt,
		CompletedAt:         record.CompletedAt,
		DuplicateGroupCount: record.DuplicateGroupCount,
		StatusCounts:        map[string]int{},
	}

	if record.Edges.Tenant != nil {
		summary.TenantSlug = record.Edges.Tenant.Slug
	}

	for _, group := range record.Edges.DuplicateGroups {
		status := string(group.Status)
		summary.StatusCounts[status]++
	}

	if summary.DuplicateGroupCount == 0 {
		summary.DuplicateGroupCount = len(record.Edges.DuplicateGroups)
	}

	return summary
}

func uuidToString(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}
