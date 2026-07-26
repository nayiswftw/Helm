//go:build linux

package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// handleListContainers returns all containers.
// GET /api/v1/containers
func handleListContainers(containerSvc *service.ContainerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containers, err := containerSvc.List(r.Context())
		if err != nil {
			if errors.Is(err, service.ErrDockerUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "docker_unavailable", "Docker socket is not available")
				return
			}
			respondError(w, http.StatusInternalServerError, "container_error", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, containers)
	}
}

// handleStartContainer starts a container by ID.
// POST /api/v1/containers/{id}/start
func handleStartContainer(containerSvc *service.ContainerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Container ID is required")
			return
		}

		if err := containerSvc.Start(r.Context(), id); err != nil {
			if errors.Is(err, service.ErrDockerUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "docker_unavailable", "Docker socket is not available")
				return
			}
			respondError(w, http.StatusInternalServerError, "action_failed", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Container started successfully",
		})
	}
}

// handleStopContainer stops a container by ID.
// POST /api/v1/containers/{id}/stop
func handleStopContainer(containerSvc *service.ContainerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Container ID is required")
			return
		}

		if err := containerSvc.Stop(r.Context(), id); err != nil {
			if errors.Is(err, service.ErrDockerUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "docker_unavailable", "Docker socket is not available")
				return
			}
			respondError(w, http.StatusInternalServerError, "action_failed", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Container stopped successfully",
		})
	}
}

// handleRestartContainer restarts a container by ID.
// POST /api/v1/containers/{id}/restart
func handleRestartContainer(containerSvc *service.ContainerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Container ID is required")
			return
		}

		if err := containerSvc.Restart(r.Context(), id); err != nil {
			if errors.Is(err, service.ErrDockerUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "docker_unavailable", "Docker socket is not available")
				return
			}
			respondError(w, http.StatusInternalServerError, "action_failed", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Container restarted successfully",
		})
	}
}
