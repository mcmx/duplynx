package tenancy

// SampleTenants provides seeded data for the Create DupLynx demo.
func SampleTenants() []Tenant {
	return []Tenant{
		{
			Slug:        "sample-tenant-a",
			Name:        "Sample Tenant A",
			Description: "Demo organization for the Create DupLynx experience",
			Machines: []Machine{
				{ID: "ares-laptop", TenantSlug: "sample-tenant-a", Name: "Ares-Laptop", Category: "personal_laptop", Hostname: "ares.local"},
				{ID: "helios-01", TenantSlug: "sample-tenant-a", Name: "Helios-Server-01", Category: "server", Hostname: "helios-01.dc"},
				{ID: "helios-02", TenantSlug: "sample-tenant-a", Name: "Helios-Server-02", Category: "server", Hostname: "helios-02.dc"},
				{ID: "atlas-01", TenantSlug: "sample-tenant-a", Name: "Atlas-Server-01", Category: "server", Hostname: "atlas-01.dc"},
				{ID: "atlas-02", TenantSlug: "sample-tenant-a", Name: "Atlas-Server-02", Category: "server", Hostname: "atlas-02.dc"},
			},
		},
	}
}
