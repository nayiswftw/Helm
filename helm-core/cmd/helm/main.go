package main

import (
	"log"

	"github.com/nayiswftw/helm/helm-core/internal/server"
)

func main() {
	log.Println("Starting Helm Core...")

	srv := server.New()

	log.Println("Listening on :8080")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}