package handlers

import (
	"net/http"

	"github.com/nayiswftw/helm/helm-core/internal/api/response"
)

type RootResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Root(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, RootResponse{
		Name:    "Helm",
		Version: "0.0.1",
	})
}