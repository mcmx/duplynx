package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ActionHTMXHandler returns a simple acknowledgement fragment for htmx updates.
func ActionHTMXHandler(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupId")
	fragment := fmt.Sprintf(`<div class="text-xs text-emerald-400">Action accepted for group %s</div>`, groupID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(fragment))
}
