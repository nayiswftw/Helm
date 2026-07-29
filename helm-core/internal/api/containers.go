//go:build linux

package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

// handleGetContainerStats returns live CPU and Memory usage for a container.
// GET /api/v1/containers/{id}/stats
func handleGetContainerStats(containerSvc *service.ContainerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Container ID is required")
			return
		}

		stats, err := containerSvc.GetStats(r.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrDockerUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "docker_unavailable", "Docker socket is not available")
				return
			}
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				respondError(w, http.StatusNotFound, "not_found", "Container not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "stats_error", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, stats)
	}
}

// handleGetContainerLogs returns recent log lines for a container.
// GET /api/v1/containers/{id}/logs?tail=100
func handleGetContainerLogs(containerSvc *service.ContainerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Container ID is required")
			return
		}

		tail := 100
		if tailParam := r.URL.Query().Get("tail"); tailParam != "" {
			if parsed, err := strconv.Atoi(tailParam); err == nil && parsed > 0 {
				tail = parsed
			}
		}

		logs, err := containerSvc.GetLogs(r.Context(), id, tail)
		if err != nil {
			if errors.Is(err, service.ErrDockerUnavailable) {
				respondError(w, http.StatusServiceUnavailable, "docker_unavailable", "Docker socket is not available")
				return
			}
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				respondError(w, http.StatusNotFound, "not_found", "Container not found")
				return
			}
			respondError(w, http.StatusInternalServerError, "logs_error", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, logs)
	}
}
