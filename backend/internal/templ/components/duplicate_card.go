package components

import (
	"html/template"
	"strings"

	"github.com/mcmx/duplynx/internal/actions"
)

// DuplicateCard renders a duplicate group card with htmx-enabled keeper/action controls.
func DuplicateCard(group actions.DuplicateGroup) template.HTML {
	var b strings.Builder
	b.WriteString(`<div class="flex flex-col gap-2 border border-slate-700 rounded-lg p-3 bg-slate-800">`)
	b.WriteString(`<div class="flex items-center justify-between">`)
	b.WriteString(`<span class="font-mono text-xs text-slate-300">` + template.HTMLEscapeString(group.Hash) + `</span>`)
	if group.KeeperMachineID != "" {
		b.WriteString(`<span class="text-xs text-emerald-400">Keeper: ` + template.HTMLEscapeString(group.KeeperMachineID) + `</span>`)
	}
	b.WriteString(`</div>`)

	b.WriteString(`<form hx-post="/duplicate-groups/` + template.HTMLEscapeString(group.ID) + `/keeper" hx-target="closest .duplicate-card" class="flex gap-2 items-center">`)
	b.WriteString(`<input type="hidden" name="tenantSlug" value="` + template.HTMLEscapeString(group.TenantSlug) + `">`)
	b.WriteString(`<input type="text" name="keeperMachineId" class="bg-slate-900 border border-slate-600 rounded px-2 py-1 text-xs" placeholder="Keeper machine">`)
	b.WriteString(`<button type="submit" class="text-xs px-2 py-1 bg-emerald-600 rounded">Assign keeper</button>`)
	b.WriteString(`</form>`)

	b.WriteString(`<form hx-post="/duplicate-groups/` + template.HTMLEscapeString(group.ID) + `/actions" hx-target="closest .duplicate-card" class="flex gap-2 items-center">`)
	b.WriteString(`<input type="hidden" name="tenantSlug" value="` + template.HTMLEscapeString(group.TenantSlug) + `">`)
	b.WriteString(`<select name="actionType" class="bg-slate-900 border border-slate-600 rounded px-2 py-1 text-xs">`)
	for _, action := range []string{"delete_copies", "create_hardlinks", "quarantine"} {
		b.WriteString(`<option value="` + action + `">` + action + `</option>`)
	}
	b.WriteString(`</select>`)
	b.WriteString(`<button type="submit" class="text-xs px-2 py-1 bg-amber-600 rounded">Run action</button>`)
	b.WriteString(`</form>`)

	b.WriteString(`<ul class="space-y-1 text-xs text-slate-400">`)
	for _, file := range group.Files {
		status := ""
		if file.Quarantined {
			status = " (quarantined)"
		}
		b.WriteString(`<li>` + template.HTMLEscapeString(file.MachineID+": "+file.Path) + status + `</li>`)
	}
	b.WriteString(`</ul>`)

	b.WriteString(`</div>`)
	return template.HTML(b.String())
}
