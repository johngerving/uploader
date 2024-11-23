package server

import (
	"context"
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

			encode(w, r, http.StatusOK, response)
		},
	)
}
