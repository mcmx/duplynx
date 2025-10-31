package scans

import "time"

// ScanSummary describes a scan available to the dashboard.
type ScanSummary struct {
	ID                  string
	TenantSlug          string
	Name                string
	InitiatedMachineID  string
	StartedAt           time.Time
	CompletedAt         time.Time
	DuplicateGroupCount int
	StatusCounts        map[string]int
}

// DuplicateGroupSummary outlines aggregate data for a duplicate group.
type DuplicateGroupSummary struct {
	ID              string
	Hash            string
	Status          string
	FileCount       int
	TotalSizeBytes  int64
	PreviewMachines []string
}
