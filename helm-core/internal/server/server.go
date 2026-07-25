package server

import (
	"net/http"

	"github.com/nayiswftw/helm/helm-core/internal/api"
)

func New() *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: api.Router(),
	}
}