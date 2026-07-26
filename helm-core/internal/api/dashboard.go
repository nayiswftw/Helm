//go:build linux

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// apiError is the standard error response format.
type apiError struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// handleDashboard returns the current system metrics snapshot.
// GET /api/v1/dashboard
func handleDashboard(dashboard *service.DashboardService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := dashboard.GetMetrics()
		if err != nil {
			logger.Error("failed to collect metrics", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(apiError{
				Error: errorDetail{
					Code:    "metrics_unavailable",
					Message: "Failed to collect system metrics",
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	}
}
