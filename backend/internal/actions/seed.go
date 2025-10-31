package actions

// SampleDuplicateGroups returns seeded duplicate groups for the Create DupLynx demo.
func SampleDuplicateGroups() []DuplicateGroup {
	return []DuplicateGroup{
		{
			ID:              "dg-001",
			ScanID:          "baseline-sweep-2025-10-01",
			TenantSlug:      "sample-tenant-a",
			Status:          "review",
			KeeperMachineID: "helios-01",
			Hash:            "hash-a1",
			Files: []DuplicateFile{
				{ID: "f1", MachineID: "ares-laptop", Path: "/Users/demo/Documents/report.pdf", SizeBytes: 1024},
				{ID: "f2", MachineID: "helios-01", Path: "/srv/docs/report.pdf", SizeBytes: 1024},
			},
		},
		{
			ID:              "dg-002",
			ScanID:          "baseline-sweep-2025-10-01",
			TenantSlug:      "sample-tenant-a",
			Status:          "action_needed",
			KeeperMachineID: "helios-02",
			Hash:            "hash-b2",
			Files: []DuplicateFile{
				{ID: "f3", MachineID: "atlas-01", Path: "/srv/media/raw.mov", SizeBytes: 2048},
				{ID: "f4", MachineID: "helios-02", Path: "/srv/media/raw.mov", SizeBytes: 2048},
			},
		},
	}
}
