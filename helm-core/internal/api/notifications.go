//go:build linux

package api

import (
	"net/http"

	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// handleTestNotification sends a test notification for debugging.
// POST /api/v1/notifications/test
func handleTestNotification(notifySvc *service.NotificationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := notifySvc.SendTest(); err != nil {
			respondError(w, http.StatusServiceUnavailable, "notifications_unavailable", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Test notification sent",
		})
	}
}
