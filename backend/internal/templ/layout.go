package templ

import (
	"fmt"
	"html/template"
	"strings"
)

// Breadcrumb builds the tenant/machine breadcrumb string.
func Breadcrumb(tenantLabel, machineLabel string) string {
	tenantLabel = strings.TrimSpace(tenantLabel)
	machineLabel = strings.TrimSpace(machineLabel)
	switch {
	case tenantLabel == "" && machineLabel == "":
		return ""
	case machineLabel == "":
		return tenantLabel
	case tenantLabel == "":
		return machineLabel
	default:
		return fmt.Sprintf("%s / %s", tenantLabel, machineLabel)
	}
}

// RenderLayout returns HTML with a consistent DupLynx shell.
func RenderLayout(title, tenantLabel, machineLabel string, body template.HTML) template.HTML {
	breadcrumb := Breadcrumb(tenantLabel, machineLabel)
	return template.HTML(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>%s</title>
    <link rel="stylesheet" href="/static/app.css">
  </head>
  <body class="min-h-screen bg-slate-900 text-slate-100">
    <header class="border-b border-slate-700 py-4">
      <div class="max-w-5xl mx-auto px-6">
        <h1 class="text-xl font-semibold">DupLynx</h1>
        <p class="text-sm text-slate-400" aria-label="Current context">%s</p>
      </div>
    </header>
    <main class="max-w-5xl mx-auto px-6 py-8">
      %s
    </main>
  </body>
</html>`,
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(breadcrumb),
		body))
}
