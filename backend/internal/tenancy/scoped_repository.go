package tenancy

import (
	"context"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/internal/scans"
)

// ScopedRepository enforces tenant boundaries when interacting with shared stores.
type ScopedRepository struct {
	scope        Scope
	scanRepo     *scans.Repository
	actionsStore *actions.Store
}

// NewScopedRepository constructs a scoped repository facade for the given tenant.
func NewScopedRepository(scope Scope, scanRepo *scans.Repository, actionsStore *actions.Store) *ScopedRepository {
	return &ScopedRepository{
		scope:        scope,
		scanRepo:     scanRepo,
		actionsStore: actionsStore,
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
	if s.actionsStore == nil {
		return nil
	}

	groups := s.actionsStore.ListByScan(scanID)
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
	if s.actionsStore == nil {
		return nil, actions.ErrGroupNotFound
	}
	group, ok := s.actionsStore.Get(groupID)
	if !ok {
		return nil, actions.ErrGroupNotFound
	}
	if group.TenantSlug != s.scope.TenantSlug {
		return nil, actions.ErrGroupNotFound
	}
	return group, nil
}
