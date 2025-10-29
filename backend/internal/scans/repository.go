package scans

import (
	"context"
	"errors"
)

var (
	ErrScanNotFound = errors.New("scan not found")
)

// Repository surfaces seeded scan summaries for the demo.
type Repository struct {
	scans map[string]ScanSummary
}

// NewRepository constructs a repository with immutable scan data.
func NewRepository(scans []ScanSummary) *Repository {
	index := make(map[string]ScanSummary, len(scans))
	for _, scan := range scans {
		index[scan.ID] = scan
	}
	return &Repository{scans: index}
}

// ListByTenant returns scan summaries for a tenant.
func (r *Repository) ListByTenant(ctx context.Context, tenantSlug string) ([]ScanSummary, error) {
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
	scan, ok := r.scans[scanID]
	if !ok {
		return ScanSummary{}, ErrScanNotFound
	}
	return scan, nil
}
