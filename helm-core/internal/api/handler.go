package api

import "github.com/nayiswftw/helm/helm-core/internal/app"

// Handler holds API route handlers and application dependencies.
type Handler struct {
	app *app.Application
}

// NewHandler constructs an API Handler.
func NewHandler(a *app.Application) *Handler {
	return &Handler{app: a}
}
