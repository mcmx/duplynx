package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mcmx/duplynx/internal/tenancy"
)

// TenantsHandler responds with tenant summary data for onboarding flows.
type TenantsHandler struct {
	Repo *tenancy.Repository
}

func (h TenantsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Repo == nil {
		http.Error(w, "tenancy repository not configured", http.StatusInternalServerError)
		return
	}

	tenants, err := h.Repo.ListTenants(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Tenants []TenantSummary `json:"tenants"`
	}{}

	for _, tenant := range tenants {
		resp.Tenants = append(resp.Tenants, TenantSummary{
			Slug:         tenant.Slug,
			Name:         tenant.Name,
			Description:  tenant.Description,
			MachineCount: len(tenant.Machines),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// TenantSummary is a lightweight representation for API responses.
type TenantSummary struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	MachineCount int    `json:"machine_count"`
}
