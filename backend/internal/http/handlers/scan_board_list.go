package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mcmx/duplynx/internal/scans"
)

// ScanListHandler lists scans for a tenant.
type ScanListHandler struct {
	Service scans.Service
}

func (h ScanListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenantSlug := chi.URLParam(r, "tenantSlug")
	scanSummaries, err := h.Service.ListTenantScans(r.Context(), tenantSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Scans []scans.ScanSummary `json:"scans"`
	}{Scans: scanSummaries})
}
