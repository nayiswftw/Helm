//go:build linux

package api

import (
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
		respondJSON(w, http.StatusOK, actions)
	}
}

// handleExecuteAction executes a predefined action by ID.
// POST /api/v1/actions/{id}/execute
func handleExecuteAction(actionSvc *service.ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Action ID parameter is required")
			return
		}

		result, err := actionSvc.Execute(id)
		if err != nil {
			if errors.Is(err, service.ErrActionNotFound) {
				respondError(w, http.StatusNotFound, "action_not_found", "Action not found")
				return
			}

			respondError(w, http.StatusInternalServerError, "execution_failed", err.Error())
			return
		}

		respondJSON(w, http.StatusAccepted, result)
	}
}
