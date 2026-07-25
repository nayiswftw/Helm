package api

import "net/http"

// Health handles GET /api/v1/health.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"name":   h.app.Config.Name,
	})
}
