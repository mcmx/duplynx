package scans

import "time"

// SampleScans provides seeded scans for the Create DupLynx demo.
func SampleScans() []ScanSummary {
	return []ScanSummary{
		{
			ID:                  "baseline-sweep-2025-10-01",
			TenantSlug:          "sample-tenant-a",
			Name:                "Baseline Sweep 2025-10-01",
			InitiatedMachineID:  "helios-01",
			StartedAt:           time.Date(2025, 10, 1, 9, 0, 0, 0, time.UTC),
			CompletedAt:         time.Date(2025, 10, 1, 10, 0, 0, 0, time.UTC),
			DuplicateGroupCount: 48,
			StatusCounts: map[string]int{
				"review":        20,
				"action_needed": 10,
				"resolved":      12,
				"archived":      6,
			},
		},
		{
			ID:                  "media-audit-2025-10-10",
			TenantSlug:          "sample-tenant-a",
			Name:                "Media Audit 2025-10-10",
			InitiatedMachineID:  "helios-02",
			StartedAt:           time.Date(2025, 10, 10, 9, 0, 0, 0, time.UTC),
			CompletedAt:         time.Date(2025, 10, 10, 11, 30, 0, 0, time.UTC),
			DuplicateGroupCount: 36,
			StatusCounts: map[string]int{
				"review":        12,
				"action_needed": 8,
				"resolved":      10,
				"archived":      6,
			},
		},
		{
			ID:                  "archive-sync-2025-10-20",
			TenantSlug:          "sample-tenant-a",
			Name:                "Archive Sync 2025-10-20",
			InitiatedMachineID:  "atlas-01",
			StartedAt:           time.Date(2025, 10, 20, 9, 0, 0, 0, time.UTC),
			CompletedAt:         time.Date(2025, 10, 20, 12, 45, 0, 0, time.UTC),
			DuplicateGroupCount: 52,
			StatusCounts: map[string]int{
				"review":        18,
				"action_needed": 14,
				"resolved":      15,
				"archived":      5,
			},
		},
	}
}
