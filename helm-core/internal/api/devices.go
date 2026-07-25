package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListDevices handles GET /api/v1/devices.
func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, h.app.Devices.List())
}

// GetDevice handles GET /api/v1/devices/{id}.
func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	device, ok := h.app.Devices.GetByID(id)
	if !ok {
		Error(w, http.StatusNotFound, "not_found", fmt.Sprintf("device not found: %s", id))
		return
	}

	JSON(w, http.StatusOK, device)
}
