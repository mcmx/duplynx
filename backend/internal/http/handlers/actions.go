package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mcmx/duplynx/internal/actions"
)

// KeeperHandler assigns a keeper machine to a duplicate group.
type KeeperHandler struct {
	Dispatcher actions.Dispatcher
}

type keeperRequest struct {
	TenantSlug      string `json:"tenantSlug"`
	KeeperMachineID string `json:"keeperMachineId"`
}

func (h KeeperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req keeperRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if req.TenantSlug == "" {
		http.Error(w, "tenantSlug required", http.StatusBadRequest)
		return
	}
	groupID := chi.URLParam(r, "groupId")
	if err := h.Dispatcher.AssignKeeper(groupID, req.TenantSlug, req.KeeperMachineID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"message": "keeper assignment recorded",
	})
}

// ActionHandler triggers duplicate management actions (delete, hardlink, quarantine).
type ActionHandler struct {
	Dispatcher actions.Dispatcher
}

type actionRequest struct {
	TenantSlug    string             `json:"tenantSlug"`
	ActionType    actions.ActionType `json:"actionType"`
	TargetFileIDs []string           `json:"targetFileIds"`
	Notes         string             `json:"notes"`
}

func (h ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req actionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if req.TenantSlug == "" || req.ActionType == "" {
		http.Error(w, "tenantSlug and actionType required", http.StatusBadRequest)
		return
	}
	groupID := chi.URLParam(r, "groupId")
	payload := map[string]any{
		"targetFileIds": req.TargetFileIDs,
		"notes":         req.Notes,
	}
	if err := h.Dispatcher.PerformAction(groupID, req.TenantSlug, "system", req.ActionType, payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "accepted",
		"message": "action queued",
	})
}
