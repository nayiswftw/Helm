package api

import "github.com/nayiswftw/helm/helm-core/internal/app"

type Handler struct {
	app *app.Application
}

func NewHandler(a *app.Application) *Handler {
	return &Handler{app: a}
}
