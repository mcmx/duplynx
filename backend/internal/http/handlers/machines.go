package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mcmx/duplynx/internal/tenancy"
)

// MachinesHandler lists machines for a tenant and records selection analytics.
type MachinesHandler struct {
	Repo *tenancy.Repository
}

func (h MachinesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Repo == nil {
		http.Error(w, "tenancy repository not configured", http.StatusInternalServerError)
		return
	}

	scope, ok := tenancy.ScopeFromContext(r.Context())
	if !ok {
		http.Error(w, "tenant scope missing", http.StatusBadRequest)
		return
	}

	machines, err := h.Repo.ListMachines(r.Context(), scope.TenantSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	selectedID := r.URL.Query().Get("selected_machine")
	if selectedID != "" {
		if machine, err := h.Repo.FindMachine(r.Context(), scope.TenantSlug, selectedID); err == nil {
			h.Repo.LogMachineSelection(scope.TenantSlug, machine)
		}
	}

	resp := struct {
		Machines []MachineSummary `json:"machines"`
	}{}

	for _, machine := range machines {
		resp.Machines = append(resp.Machines, MachineSummary{
			ID:       machine.ID,
			Name:     machine.Name,
			Category: machine.Category,
			Hostname: machine.Hostname,
			Role:     machine.Role,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// MachineSummary is exposed to clients selecting keepers and running scans.
type MachineSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Hostname string `json:"hostname"`
	Role     string `json:"role"`
}
