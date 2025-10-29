package unit_test

import (
	"testing"
	"time"

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
	markup := templ.BoardPage(summary, map[string][]scans.DuplicateGroupSummary{})
	if len(markup) == 0 {
		t.Fatalf("expected board markup to render content")
	}
}
