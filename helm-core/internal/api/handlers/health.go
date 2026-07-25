package handlers

import (
	"net/http"

	"github.com/nayiswftw/helm/helm-core/internal/api/response"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, HealthResponse{
		Status:  "ok",
		Version: "0.0.1",
	})
}