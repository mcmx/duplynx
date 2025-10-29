package templ

import (
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/mcmx/duplynx/internal/scans"
)

var statusOrder = []string{"review", "action_needed", "resolved", "archived"}

// BoardPage renders the board columns for a scan.
func BoardPage(summary scans.ScanSummary, groups map[string][]scans.DuplicateGroupSummary) template.HTML {
	var b strings.Builder
	b.WriteString(`<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">`)

	for _, status := range statusOrder {
		b.WriteString(`<section class="bg-slate-800 border border-slate-700 rounded-lg">`)
		b.WriteString(`<header class="flex items-center justify-between px-4 py-3 border-b border-slate-700">`)
		b.WriteString(`<h2 class="text-sm font-semibold uppercase tracking-wide">` + template.HTMLEscapeString(strings.ReplaceAll(status, "_", " ")) + `</h2>`)
		b.WriteString(`<span class="text-xs text-slate-400">` + fmt.Sprint(summary.StatusCounts[status]) + `</span>`)
		b.WriteString(`</header>`)

		b.WriteString(`<ul class="divide-y divide-slate-700">`)
		list := groups[status]
		sort.SliceStable(list, func(i, j int) bool { return list[i].Hash < list[j].Hash })
		for _, group := range list {
			b.WriteString(`<li class="px-4 py-3">`)
			b.WriteString(`<p class="font-mono text-xs text-slate-300">` + template.HTMLEscapeString(group.Hash) + `</p>`)
			b.WriteString(`<p class="text-sm text-white">` + fmt.Sprintf("%d files • %d bytes", group.FileCount, group.TotalSizeBytes) + `</p>`)
			b.WriteString(`</li>`)
		}
		b.WriteString(`</ul>`)
		b.WriteString(`</section>`)
	}

	b.WriteString(`</div>`)
	return template.HTML(b.String())
}
