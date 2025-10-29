package templ

import (
	"html/template"
	"strings"

	"github.com/mcmx/duplynx/internal/tenancy"
)

// LaunchPage renders tenant and machine selection cards.
func LaunchPage(tenants []tenancy.Tenant) template.HTML {
	var b strings.Builder
	b.WriteString(`<section class="space-y-6">`)
	for _, tenant := range tenants {
		b.WriteString(`<article class="border border-slate-700 rounded-lg p-4">`)
		b.WriteString(`<h2 class="text-lg font-semibold">` + template.HTMLEscapeString(tenant.Name) + `</h2>`)
		if tenant.Description != "" {
			b.WriteString(`<p class="text-sm text-slate-400">` + template.HTMLEscapeString(tenant.Description) + `</p>`)
		}
		b.WriteString(`<ul class="mt-3 space-y-2">`)
		for _, machine := range tenant.Machines {
			b.WriteString(`<li class="flex items-center justify-between border border-slate-800 rounded px-3 py-2">`)
			b.WriteString(`<span>` + template.HTMLEscapeString(machine.Name) + `</span>`)
			b.WriteString(`<span class="text-xs uppercase tracking-wide text-slate-500">` + template.HTMLEscapeString(machine.Category) + `</span>`)
			b.WriteString(`</li>`)
		}
		b.WriteString(`</ul>`)
		b.WriteString(`</article>`)
	}
	b.WriteString(`</section>`)
	return template.HTML(b.String())
}
