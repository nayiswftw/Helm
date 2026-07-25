package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Router constructs the HTTP router and registers API routes.
func Router(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.Health)
		r.Get("/dashboard", h.Dashboard)
		r.Get("/devices", h.ListDevices)
		r.Get("/devices/{id}", h.GetDevice)
	})

	return r
}
