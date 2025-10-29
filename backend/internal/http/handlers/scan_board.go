package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mcmx/duplynx/internal/scans"
)

// ScanBoardHandler returns the scan summary with duplicate status counts.
type ScanBoardHandler struct {
	Service scans.Service
}

func (h ScanBoardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	scanID := chi.URLParam(r, "scanID")
	summary, err := h.Service.GetScan(r.Context(), scanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
