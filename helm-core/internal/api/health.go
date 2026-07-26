package api

import (
	"net/http"
)

// handleHealth returns a simple liveness check response.
// GET /health → {"status": "ok"}
func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}
