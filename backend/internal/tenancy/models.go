package tenancy

// Tenant represents an organization with associated machines.
type Tenant struct {
	Slug        string
	Name        string
	Description string
	Machines    []Machine
}

// Machine represents a predefined host that can participate in scans.
type Machine struct {
	ID         string
	TenantSlug string
	Name       string
	Category   string
	Hostname   string
	Role       string
}
