package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nayiswftw/helm/helm-core/internal/api/handlers"
	"github.com/nayiswftw/helm/helm-core/internal/api/response"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", handlers.Root)
		r.Get("/health", handlers.Health)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusNotFound, map[string]any{
			"error": map[string]string{
				"code":    "route_not_found",
				"message": "The requested route does not exist.",
			},
		})
	})

	return r
}