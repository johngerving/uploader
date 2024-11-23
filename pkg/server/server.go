package server

import (
	"log"
	"net/http"
)

func NewServer(
	logger *log.Logger,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux)

	var handler http.Handler = mux

	return handler
}

func addRoutes(
	mux *http.ServeMux,
	// logger *log.Logger,
) {
	mux.Handle("/healthz", handleHealth())
}