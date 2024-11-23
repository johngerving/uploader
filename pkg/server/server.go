package server

import (
	"log"
	"net/http"

	"github.com/johngerving/uploader/repository"
)

func NewServer(
	logger *log.Logger,
	queries *repository.Queries,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, logger, queries)

	var handler http.Handler = mux

	return handler
}

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	queries *repository.Queries,
) {
	mux.Handle("GET /healthz", handleHealth())
	mux.Handle("POST /uploads", handlePostUpload(logger, queries))
	mux.Handle("POST /uploads/{id}/parts", handlePostPart(logger, queries))
}
