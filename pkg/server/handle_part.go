package server

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/johngerving/uploader/repository"
)

func handlePostPart(logger *log.Logger, queries *repository.Queries) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			uploadId := r.PathValue("id")
			if uploadId == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			_, err := queries.FindUploadById(context.Background(), uploadId)
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				logger.Printf("error getting row with ID %v: %v\n", uploadId, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Printf("error reading body: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Println(string(body))
		},
	)
}
