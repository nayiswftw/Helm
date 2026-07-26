package api

import (
	"encoding/json"
	"net/http"
)

// handleHealth returns a simple liveness check response.
// GET /health → {"status": "ok"}
func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	}
}
