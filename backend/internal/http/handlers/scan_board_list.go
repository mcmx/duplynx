package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mcmx/duplynx/internal/scans"
	"github.com/mcmx/duplynx/internal/tenancy"
)

// ScanListHandler lists scans for a tenant.
type ScanListHandler struct {
	Service scans.Service
}

func (h ScanListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	scope, ok := tenancy.ScopeFromContext(r.Context())
	if !ok {
		http.Error(w, "tenant scope missing", http.StatusBadRequest)
		return
	}

	scoped := tenancy.NewScopedRepository(scope, h.Service.Repo, nil)
	scanSummaries, err := scoped.ListScans(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Scans []scans.ScanSummary `json:"scans"`
	}{Scans: scanSummaries}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
