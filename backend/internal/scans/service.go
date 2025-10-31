package scans

import "context"

// Service provides higher-level scan aggregation helpers.
type Service struct {
	Repo *Repository
}

func (s Service) ListTenantScans(ctx context.Context, tenantSlug string) ([]ScanSummary, error) {
	if s.Repo == nil {
		return nil, nil
	}
	return s.Repo.ListByTenant(ctx, tenantSlug)
}

func (s Service) GetScan(ctx context.Context, scanID string) (ScanSummary, error) {
	if s.Repo == nil {
		return ScanSummary{}, ErrScanNotFound
	}
	return s.Repo.Get(ctx, scanID)
}
