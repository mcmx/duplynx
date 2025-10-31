package perf_test

import (
	"fmt"
	"testing"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/templ"
)

func BenchmarkBoardRendering(b *testing.B) {
	summary := scans.ScanSummary{
		ID:         "baseline-sweep-2025-10-01",
		TenantSlug: "sample-tenant-a",
		StatusCounts: map[string]int{
			"review":        50,
			"action_needed": 40,
			"resolved":      60,
			"archived":      30,
		},
	}

	groups := map[string][]actions.DuplicateGroup{
		"review":        make([]actions.DuplicateGroup, 0, summary.StatusCounts["review"]),
		"action_needed": make([]actions.DuplicateGroup, 0, summary.StatusCounts["action_needed"]),
		"resolved":      make([]actions.DuplicateGroup, 0, summary.StatusCounts["resolved"]),
		"archived":      make([]actions.DuplicateGroup, 0, summary.StatusCounts["archived"]),
	}

	for status, total := range summary.StatusCounts {
		for i := 0; i < total; i++ {
			group := actions.DuplicateGroup{
				ID:         fmt.Sprintf("%s-%03d", status, i),
				ScanID:     summary.ID,
				TenantSlug: summary.TenantSlug,
				Status:     status,
				Hash:       fmt.Sprintf("hash-%s-%03d", status, i),
				Files: []actions.DuplicateFile{
					{ID: fmt.Sprintf("%s-%03d-a", status, i), MachineID: "ares-laptop", Path: "/path/a", SizeBytes: 1024},
					{ID: fmt.Sprintf("%s-%03d-b", status, i), MachineID: "helios-01", Path: "/path/b", SizeBytes: 1024},
					{ID: fmt.Sprintf("%s-%03d-c", status, i), MachineID: "atlas-01", Path: "/path/c", SizeBytes: 1024},
				},
			}
			groups[status] = append(groups[status], group)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		markup := templ.BoardPage(summary, groups)
		if len(markup) == 0 {
			b.Fatal("expected board markup")
		}
	}
}
