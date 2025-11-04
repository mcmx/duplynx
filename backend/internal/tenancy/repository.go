package tenancy

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/mcmx/duplynx/ent"
	entmachine "github.com/mcmx/duplynx/ent/machine"
	enttenant "github.com/mcmx/duplynx/ent/tenant"
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
	client      *ent.Client
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

// NewRepositoryFromClient constructs a repository backed by the DupLynx Ent client.
func NewRepositoryFromClient(client *ent.Client, audit *AuditLogger) *Repository {
	return &Repository{
		client: client,
		audit:  audit,
	}
}

// ListTenants returns tenants sorted by seed order.
func (r *Repository) ListTenants(ctx context.Context) ([]Tenant, error) {
	if r.client != nil {
		return r.listTenantsFromClient(ctx)
	}

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
	if r.client != nil {
		machines, err := r.listMachinesFromClient(ctx, tenantSlug)
		if err != nil {
			return nil, err
		}
		if r.audit != nil {
			r.audit.LogTenantSelection(tenantSlug)
		}
		return machines, nil
	}

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
	if r.client != nil {
		return r.findMachineFromClient(ctx, tenantSlug, machineID)
	}

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
	if r.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		record, err := r.client.Tenant.
			Query().
			Where(enttenant.SlugEQ(slug)).
			Only(ctx)
		if err != nil {
			return Tenant{}, false
		}
		return Tenant{
			Slug:        record.Slug,
			Name:        record.Name,
			Description: record.Description,
		}, true
	}

	tenant, ok := r.tenants[slug]
	return tenant, ok
}

func (r *Repository) listTenantsFromClient(ctx context.Context) ([]Tenant, error) {
	if r.client == nil {
		return nil, errors.New("ent client not configured")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	records, err := r.client.Tenant.
		Query().
		WithMachines(func(mq *ent.MachineQuery) {
			mq.Order(entmachine.ByName()).
				WithTenant()
		}).
		Order(enttenant.ByName()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Tenant, 0, len(records))
	for _, record := range records {
		out = append(out, convertTenant(record))
	}
	return out, nil
}

func (r *Repository) listMachinesFromClient(ctx context.Context, tenantSlug string) ([]Machine, error) {
	if r.client == nil {
		return nil, errors.New("ent client not configured")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	records, err := r.client.Machine.
		Query().
		Where(entmachine.HasTenantWith(enttenant.SlugEQ(tenantSlug))).
		Order(entmachine.ByName()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	machines := make([]Machine, 0, len(records))
	for _, record := range records {
		machines = append(machines, convertMachine(record))
	}
	return machines, nil
}

func (r *Repository) findMachineFromClient(ctx context.Context, tenantSlug, machineID string) (Machine, error) {
	if r.client == nil {
		return Machine{}, errors.New("ent client not configured")
	}
	u, err := uuid.Parse(machineID)
	if err != nil {
		return Machine{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	record, err := r.client.Machine.
		Query().
		Where(entmachine.IDEQ(u)).
		WithTenant().
		Only(ctx)
	if err != nil {
		return Machine{}, err
	}
	if record.Edges.Tenant == nil || record.Edges.Tenant.Slug != tenantSlug {
		return Machine{}, errors.New("machine not found")
	}
	return convertMachine(record), nil
}

func convertTenant(record *ent.Tenant) Tenant {
	if record == nil {
		return Tenant{}
	}
	tenant := Tenant{
		Slug:        record.Slug,
		Name:        record.Name,
		Description: record.Description,
	}
	for _, machine := range record.Edges.Machines {
		m := convertMachine(machine)
		if m.TenantSlug == "" {
			m.TenantSlug = record.Slug
		}
		tenant.Machines = append(tenant.Machines, m)
	}
	return tenant
}

func convertMachine(record *ent.Machine) Machine {
	if record == nil {
		return Machine{}
	}
	var tenantSlug string
	if record.Edges.Tenant != nil {
		tenantSlug = record.Edges.Tenant.Slug
	}
	return Machine{
		ID:         record.ID.String(),
		TenantSlug: tenantSlug,
		Name:       record.Name,
		Category:   string(record.Category),
		Hostname:   record.Hostname,
		Role:       record.Role,
	}
}
