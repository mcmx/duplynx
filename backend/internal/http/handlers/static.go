package handlers

import (
	"net/http"

	"github.com/mcmx/duplynx/internal/tenancy"
)

// StaticHandler serves static assets while preserving tenant headers for caches and clients.
type StaticHandler struct {
	Root http.FileSystem
}

func (h StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Root == nil {
		http.NotFound(w, r)
		return
	}

	if tenantSlug := r.Header.Get(tenancy.HeaderTenantSlug); tenantSlug != "" {
		w.Header().Set(tenancy.HeaderTenantSlug, tenantSlug)
	}
	w.Header().Add("Vary", tenancy.HeaderTenantSlug)
	w.Header().Set("Cache-Control", "public, max-age=300")

	http.StripPrefix("/static/", http.FileServer(h.Root)).ServeHTTP(w, r)
}
