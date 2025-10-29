package unit_test

import (
	"html/template"
	"strings"
	"testing"

	"github.com/mcmx/duplynx/internal/templ"
)

func TestBreadcrumbRendersTenantAndMachine(t *testing.T) {
	breadcrumb := templ.Breadcrumb("Sample Tenant A", "Ares-Laptop")
	if breadcrumb != "Sample Tenant A / Ares-Laptop" {
		t.Fatalf("unexpected breadcrumb: %q", breadcrumb)
	}

	html := templ.RenderLayout("DupLynx", "Sample Tenant A", "Ares-Laptop", template.HTML("<p>Body</p>"))
	if !strings.Contains(string(html), "Sample Tenant A / Ares-Laptop") {
		t.Fatalf("layout did not include breadcrumb: %s", html)
	}
}
