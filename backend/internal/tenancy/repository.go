package tenancy

import (
	"context"
	"errors"
	"sort"
)

var (
	// ErrTenantNotFound is returned when a tenant slug cannot be resolved.
	ErrTenantNotFound = errors.New("tenant not found")
)

// Repository exposes tenant and machine lookup operations backed by seed data.
type Repository struct {
	audit       *AuditLogger
	tenants     map[string]Tenant
	tenantOrder []string
}

// NewRepository constructs an in-memory repository using the supplied tenants.
func NewRepository(tenants []Tenant, audit *AuditLogger) *Repository {
	m := make(map[string]Tenant, len(tenants))
	order := make([]string, 0, len(tenants))
	for _, tenant := range tenants {
		m[tenant.Slug] = tenant
		order = append(order, tenant.Slug)
	}
	return &Repository{audit: audit, tenants: m, tenantOrder: order}
}

// ListTenants returns tenants sorted by seed order.
func (r *Repository) ListTenants(ctx context.Context) ([]Tenant, error) {
	out := make([]Tenant, 0, len(r.tenants))
	for _, slug := range r.tenantOrder {
		if tenant, ok := r.tenants[slug]; ok {
			out = append(out, tenant)
		}
	}
	return out, nil
}

// ListMachines returns machines for a tenant, logging the selection.
func (r *Repository) ListMachines(ctx context.Context, tenantSlug string) ([]Machine, error) {
	tenant, ok := r.tenants[tenantSlug]
	if !ok {
		return nil, ErrTenantNotFound
	}
	if r.audit != nil {
		r.audit.LogTenantSelection(tenantSlug)
	}

	machines := make([]Machine, len(tenant.Machines))
	copy(machines, tenant.Machines)

	sort.SliceStable(machines, func(i, j int) bool {
		return machines[i].Name < machines[j].Name
	})

	return machines, nil
}

// FindMachine retrieves a specific machine by ID.
func (r *Repository) FindMachine(ctx context.Context, tenantSlug, machineID string) (Machine, error) {
	machines, err := r.ListMachines(ctx, tenantSlug)
	if err != nil {
		return Machine{}, err
	}
	for _, m := range machines {
		if m.ID == machineID {
			return m, nil
		}
	}
	return Machine{}, errors.New("machine not found")
}

// LogMachineSelection delegates to the audit logger for tracking.
func (r *Repository) LogMachineSelection(tenantSlug string, machine Machine) {
	if r.audit != nil {
		r.audit.LogMachineSelection(tenantSlug, machine)
	}
}

// Audit returns the current audit logger (may be nil).
func (r *Repository) Audit() *AuditLogger {
	return r.audit
}

// Tenant returns a tenant by slug if it exists.
func (r *Repository) Tenant(slug string) (Tenant, bool) {
	tenant, ok := r.tenants[slug]
	return tenant, ok
}
