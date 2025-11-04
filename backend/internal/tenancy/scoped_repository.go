package tenancy

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/internal/scans"
)

// ScopedRepository enforces tenant boundaries when interacting with shared stores.
type ScopedRepository struct {
	scope       Scope
	scanRepo    *scans.Repository
	actionsRepo *actions.Repository
}

// NewScopedRepository constructs a scoped repository facade for the given tenant.
func NewScopedRepository(scope Scope, scanRepo *scans.Repository, actionsRepo *actions.Repository) *ScopedRepository {
	return &ScopedRepository{
		scope:       scope,
		scanRepo:    scanRepo,
		actionsRepo: actionsRepo,
	}
}

// ListScans returns scans belonging to the scoped tenant.
func (s *ScopedRepository) ListScans(ctx context.Context) ([]scans.ScanSummary, error) {
	if s.scanRepo == nil {
		return nil, nil
	}
	return s.scanRepo.ListByTenant(ctx, s.scope.TenantSlug)
}

// GetScan retrieves a scan by ID, returning ErrScanNotFound when outside the scope.
func (s *ScopedRepository) GetScan(ctx context.Context, scanID string) (scans.ScanSummary, error) {
	if s.scanRepo == nil {
		return scans.ScanSummary{}, scans.ErrScanNotFound
	}
	scan, err := s.scanRepo.Get(ctx, scanID)
	if err != nil {
		return scans.ScanSummary{}, err
	}
	if scan.TenantSlug != s.scope.TenantSlug {
		return scans.ScanSummary{}, scans.ErrScanNotFound
	}
	return scan, nil
}

// DuplicateGroups returns duplicate groups for the scan filtered by tenant scope.
func (s *ScopedRepository) DuplicateGroups(scanID string) []actions.DuplicateGroup {
	if s.actionsRepo == nil {
		return nil
	}

	id, err := uuid.Parse(scanID)
	if err != nil {
		return nil
	}

	groups, err := s.actionsRepo.ListByScan(context.Background(), id)
	if err != nil {
		return nil
	}
	filtered := make([]actions.DuplicateGroup, 0, len(groups))
	for _, group := range groups {
		if group.TenantSlug == s.scope.TenantSlug {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

// GetDuplicateGroup fetches a duplicate group when it belongs to the scoped tenant.
func (s *ScopedRepository) GetDuplicateGroup(groupID string) (*actions.DuplicateGroup, error) {
	if s.actionsRepo == nil {
		return nil, actions.ErrGroupNotFound
	}
	id, err := uuid.Parse(groupID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", actions.ErrInvalidGroupID, err)
	}

	group, err := s.actionsRepo.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	if group.TenantSlug != s.scope.TenantSlug {
		return nil, actions.ErrGroupNotFound
	}
	return &group, nil
}
