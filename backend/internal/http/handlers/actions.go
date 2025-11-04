package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mcmx/duplynx/internal/actions"
	"github.com/mcmx/duplynx/internal/tenancy"
)

// KeeperHandler assigns a keeper machine to a duplicate group.
type KeeperHandler struct {
	Dispatcher *actions.Dispatcher
}

type keeperRequest struct {
	TenantSlug      string `json:"tenantSlug"`
	KeeperMachineID string `json:"keeperMachineId"`
}

func (h KeeperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Dispatcher == nil {
		http.Error(w, "actions dispatcher unavailable", http.StatusServiceUnavailable)
		return
	}

	var req keeperRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	scope, ok := tenancy.ScopeFromContext(r.Context())
	if !ok {
		http.Error(w, "tenant scope missing", http.StatusBadRequest)
		return
	}
	if req.TenantSlug != "" && req.TenantSlug != scope.TenantSlug {
		http.Error(w, "tenant scope violation", http.StatusNotFound)
		return
	}
	tenantSlug := scope.TenantSlug
	groupID := chi.URLParam(r, "groupId")
	if err := h.Dispatcher.AssignKeeper(r.Context(), groupID, tenantSlug, req.KeeperMachineID); err != nil {
		status := statusFromActionsError(err)
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"message": "keeper assignment recorded",
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// ActionHandler triggers duplicate management actions (delete, hardlink, quarantine).
type ActionHandler struct {
	Dispatcher *actions.Dispatcher
}

type actionRequest struct {
	TenantSlug    string             `json:"tenantSlug"`
	ActionType    actions.ActionType `json:"actionType"`
	TargetFileIDs []string           `json:"targetFileIds"`
	Notes         string             `json:"notes"`
}

func (h ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Dispatcher == nil {
		http.Error(w, "actions dispatcher unavailable", http.StatusServiceUnavailable)
		return
	}

	var req actionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	scope, ok := tenancy.ScopeFromContext(r.Context())
	if !ok {
		http.Error(w, "tenant scope missing", http.StatusBadRequest)
		return
	}
	if req.TenantSlug != "" && req.TenantSlug != scope.TenantSlug {
		http.Error(w, "tenant scope violation", http.StatusNotFound)
		return
	}
	if req.ActionType == "" {
		http.Error(w, "actionType required", http.StatusBadRequest)
		return
	}
	groupID := chi.URLParam(r, "groupId")
	payload := map[string]any{
		"targetFileIds": req.TargetFileIDs,
		"notes":         req.Notes,
	}
	if err := h.Dispatcher.PerformAction(r.Context(), groupID, scope.TenantSlug, "system", req.ActionType, payload); err != nil {
		status := statusFromActionsError(err)
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"status":  "accepted",
		"message": "action queued",
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func statusFromActionsError(err error) int {
	switch {
	case errors.Is(err, actions.ErrGroupNotFound):
		return http.StatusNotFound
	case errors.Is(err, actions.ErrKeeperMachineID),
		errors.Is(err, actions.ErrInvalidGroupID),
		errors.Is(err, actions.ErrInvalidMachineID):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
