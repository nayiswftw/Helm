package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nayiswftw/helm/helm-core/internal/api"
	"github.com/nayiswftw/helm/helm-core/internal/app"
	"github.com/nayiswftw/helm/helm-core/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	app := app.New(cfg)

	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: api.Router(api.NewHandler(app)),
	}

	go func() {
		app.Logger.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failure: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	app.Logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Error("shutdown failed", "error", err)
	}
}
