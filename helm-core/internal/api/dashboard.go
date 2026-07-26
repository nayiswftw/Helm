//go:build linux

package api

import (
	"log/slog"
	"net/http"

	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// handleDashboard returns the current system metrics snapshot.
// GET /api/v1/dashboard
func handleDashboard(dashboard *service.DashboardService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := dashboard.GetMetrics()
		if err != nil {
			logger.Error("failed to collect metrics", "error", err)
			respondError(w, http.StatusInternalServerError, "metrics_unavailable", "Failed to collect system metrics")
			return
		}

		respondJSON(w, http.StatusOK, metrics)
	}
}
