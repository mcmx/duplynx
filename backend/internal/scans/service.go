package scans

import (
	"context"

	"github.com/mcmx/duplynx/internal/scans"
)

// Service provides higher-level scan aggregation helpers.
type Service struct {
	ScansRepo *scans.Repository
}

func (s Service) ListTenantScans(ctx context.Context, tenantSlug string) ([]scans.ScanSummary, error) {
	if s.ScansRepo == nil {
		return nil, nil
	}
	return s.ScansRepo.ListByTenant(ctx, tenantSlug)
}

func (s Service) GetScan(ctx context.Context, scanID string) (scans.ScanSummary, error) {
	return s.ScansRepo.Get(ctx, scanID)
}
