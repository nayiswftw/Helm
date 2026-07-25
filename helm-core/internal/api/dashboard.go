package api

import "net/http"

// Dashboard handles GET /api/v1/dashboard.
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, h.app.Dashboard.Get())
}
