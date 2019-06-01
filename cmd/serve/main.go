package main

import (
	"net/http"

	"github.com/lob/logger-go"
	"github.com/shreyasj2006/goyagi/pkg/server"
)

func main() {
	log := logger.New()

	srv := server.New()

	log.Info("server started")

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Err(err).Fatal("server stopped")
	}

	log.Info("server stopped")
}
