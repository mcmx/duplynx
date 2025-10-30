package unit_test

import (
	"testing"
	"time"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/templ"
)

func TestBoardPageRendersStatusCounts(t *testing.T) {
	summary := scans.ScanSummary{
		ID:                  "baseline-sweep-2025-10-01",
		TenantSlug:          "sample-tenant-a",
		Name:                "Baseline Sweep 2025-10-01",
		StartedAt:           time.Now(),
		CompletedAt:         time.Now(),
		DuplicateGroupCount: 10,
		StatusCounts: map[string]int{
			"review":        3,
			"action_needed": 2,
			"resolved":      4,
			"archived":      1,
		},
	}
	groups := map[string][]actions.DuplicateGroup{
		"review": {
			{ID: "dg-001", TenantSlug: "sample-tenant-a", Hash: "hash-a1", Files: []actions.DuplicateFile{{MachineID: "ares-laptop", Path: "/file"}}},
		},
	}

	markup := templ.BoardPage(summary, groups)
	if len(markup) == 0 {
		t.Fatalf("expected board markup to render content")
	}
}
