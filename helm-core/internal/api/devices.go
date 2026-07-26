//go:build linux

package api

import (
	"encoding/json"
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(devices)
	}
}

// handleGetDevice returns details for a specific device.
// GET /api/v1/devices/{id}
func handleGetDevice(deviceSvc *service.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(apiError{
				Error: errorDetail{
					Code:    "bad_request",
					Message: "Device ID parameter is required",
				},
			})
			return
		}

		dev, err := deviceSvc.GetByID(id)
		if err != nil {
			if errors.Is(err, service.ErrDeviceNotFound) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(apiError{
					Error: errorDetail{
						Code:    "device_not_found",
						Message: "Device not found",
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(apiError{
				Error: errorDetail{
					Code:    "internal_error",
					Message: "Failed to retrieve device",
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dev)
	}
}
