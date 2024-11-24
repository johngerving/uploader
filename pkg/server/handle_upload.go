package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/johngerving/uploader/repository"
)

func handlePostUpload(logger *log.Logger, queries *repository.Queries) http.Handler {
	type response struct {
		ID string `json:"id"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := uuid.New()

			uploadId, err := queries.CreateUpload(context.Background(), id.String())
			if err != nil {
				logger.Printf("error creating upload with ID %v: %v", uploadId, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			response := response{
				ID: uploadId,
			}

			encode(w, http.StatusOK, response)
		},
	)
}

func handleGetUpload(logger *log.Logger, queries *repository.Queries) http.Handler {
	type response struct {
		ID string `json:"id"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			_, err := queries.FindUploadById(context.Background(), id)
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				logger.Printf("error finding upload with ID %v: %v", id, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			resp := response{
				ID: id,
			}
			encode(w, http.StatusOK, &resp)
		},
	)
}
