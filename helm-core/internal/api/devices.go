//go:build linux

package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayiswftw/helm/helm-core/internal/service"
)

// handleListDevices returns all registered devices.
// GET /api/v1/devices
func handleListDevices(deviceSvc *service.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices := deviceSvc.GetAll()
		respondJSON(w, http.StatusOK, devices)
	}
}

// handleGetDevice returns details for a specific device.
// GET /api/v1/devices/{id}
func handleGetDevice(deviceSvc *service.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			respondError(w, http.StatusBadRequest, "bad_request", "Device ID parameter is required")
			return
		}

		dev, err := deviceSvc.GetByID(id)
		if err != nil {
			if errors.Is(err, service.ErrDeviceNotFound) {
				respondError(w, http.StatusNotFound, "device_not_found", "Device not found")
				return
			}

			respondError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve device")
			return
		}

		respondJSON(w, http.StatusOK, dev)
	}
}
