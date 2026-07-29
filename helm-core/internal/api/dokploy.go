//go:build linux

package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// handleListProjects returns all Dokploy projects.
// GET /api/v1/dokploy/projects
func handleListProjects(dokploySvc *service.DokployService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := dokploySvc.ListProjects(r.Context())
		if err != nil {
			if errors.Is(err, service.ErrDokployUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "dokploy_unavailable", "Dokploy is not configured")
				return
			}
			respondError(w, http.StatusInternalServerError, "dokploy_error", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, projects)
	}
}

// handleGetApplication returns details for a specific Dokploy application.
// GET /api/v1/dokploy/applications/{id}
func handleGetApplication(dokploySvc *service.DokployService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Application ID is required")
			return
		}

		app, err := dokploySvc.GetApplication(r.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrDokployUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "dokploy_unavailable", "Dokploy is not configured")
				return
			}
			respondError(w, http.StatusInternalServerError, "dokploy_error", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, app)
	}
}

// handleDeployApplication triggers a deployment.
// POST /api/v1/dokploy/applications/{id}/deploy
func handleDeployApplication(dokploySvc *service.DokployService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Application ID is required")
			return
		}

		if err := dokploySvc.Deploy(r.Context(), id); err != nil {
			if errors.Is(err, service.ErrDokployUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "dokploy_unavailable", "Dokploy is not configured")
				return
			}
			respondError(w, http.StatusInternalServerError, "deploy_failed", err.Error())
			return
		}

		respondJSON(w, http.StatusAccepted, map[string]string{
			"status":  "accepted",
			"message": "Deployment triggered",
		})
	}
}

// handleRedeployApplication triggers a redeployment.
// POST /api/v1/dokploy/applications/{id}/redeploy
func handleRedeployApplication(dokploySvc *service.DokployService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Application ID is required")
			return
		}

		if err := dokploySvc.Redeploy(r.Context(), id); err != nil {
			if errors.Is(err, service.ErrDokployUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "dokploy_unavailable", "Dokploy is not configured")
				return
			}
			respondError(w, http.StatusInternalServerError, "redeploy_failed", err.Error())
			return
		}

		respondJSON(w, http.StatusAccepted, map[string]string{
			"status":  "accepted",
			"message": "Redeployment triggered",
		})
	}
}

// handleListDeployments returns deployment history for an application.
// GET /api/v1/dokploy/applications/{id}/deployments
func handleListDeployments(dokploySvc *service.DokployService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Application ID is required")
			return
		}

		deployments, err := dokploySvc.ListDeployments(r.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrDokployUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "dokploy_unavailable", "Dokploy is not configured")
				return
			}
			respondError(w, http.StatusInternalServerError, "dokploy_error", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, deployments)
	}
}
