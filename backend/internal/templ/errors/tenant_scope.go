package templerrors

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/mcmx/duplynx/internal/templ"
	"github.com/mcmx/duplynx/internal/tenancy"
)

// Renderer renders tenant scope violations using the shared templ layout.
type Renderer struct{}

// RenderTenantScopeViolation satisfies tenancy.ScopeViolationRenderer.
func (Renderer) RenderTenantScopeViolation(w http.ResponseWriter, r *http.Request, violation tenancy.ScopeViolation) {
	if violation.CurrentTenant != "" {
		w.Header().Set(tenancy.HeaderTenantSlug, violation.CurrentTenant)
	}
	body := RenderTenantScopeViolation(violation)
	markup := templ.RenderLayout("Tenant Scope Violation", violation.CurrentTenant, "", body)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(markup))
}

// RenderTenantScopeViolation builds the inner HTML fragment describing the violation.
func RenderTenantScopeViolation(violation tenancy.ScopeViolation) template.HTML {
	var b strings.Builder
	b.WriteString(`<section class="space-y-4">`)
	b.WriteString(`<h2 class="text-lg font-semibold text-amber-400">Tenant scope conflict</h2>`)
	b.WriteString(`<p class="text-sm text-slate-300">`)
	switch {
	case violation.CurrentTenant != "" && violation.RequestedSlug != "":
		b.WriteString(template.HTMLEscapeString(fmt.Sprintf("You are currently scoped to %s but attempted to access %s.", violation.CurrentTenant, violation.RequestedSlug)))
	case violation.RequestedSlug != "":
		b.WriteString(template.HTMLEscapeString(fmt.Sprintf("The requested tenant %s is unavailable in this session.", violation.RequestedSlug)))
	default:
		b.WriteString("A tenant scope is required to complete this request.")
	}
	b.WriteString(`</p>`)

	if violation.Reason != "" {
		b.WriteString(`<p class="text-xs text-slate-500">` + template.HTMLEscapeString(violation.Reason) + `</p>`)
	}

	b.WriteString(`<div class="flex gap-2">`)
	b.WriteString(`<a class="text-sm text-emerald-400 hover:text-emerald-300" href="/">Return to tenant selection</a>`)
	b.WriteString(`</div>`)
	b.WriteString(`</section>`)

	return template.HTML(b.String())
}
