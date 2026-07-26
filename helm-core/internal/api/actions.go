//go:build linux

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// handleListActions returns all available predefined administrative actions.
// GET /api/v1/actions
func handleListActions(actionSvc *service.ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actions := actionSvc.ListActions()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(actions)
	}
}

// handleExecuteAction executes a predefined action by ID.
// POST /api/v1/actions/{id}/execute
func handleExecuteAction(actionSvc *service.ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiError{
				Error: errorDetail{
					Code:    "bad_request",
					Message: "Action ID parameter is required",
				},
			})
			return
		}

		result, err := actionSvc.Execute(id)
		if err != nil {
			if errors.Is(err, service.ErrActionNotFound) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(apiError{
					Error: errorDetail{
						Code:    "action_not_found",
						Message: "Action not found",
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(apiError{
				Error: errorDetail{
					Code:    "execution_failed",
					Message: err.Error(),
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(result)
	}
}
