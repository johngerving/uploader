package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

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

			upload, err := queries.CreateUpload(context.Background(), id.String())
			if err != nil {
				logger.Printf("error creating upload with ID %v: %v", upload.ID, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			response := response{
				ID: upload.ID,
			}

			encode(w, http.StatusOK, response)
		},
	)
}

func handlePutUpload(logger *log.Logger, queries *repository.Queries) http.Handler {
	type responseError struct {
		Message string `json:"message"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			complete := r.URL.Query().Get("complete")

			if strings.ToLower(complete) != "true" {
				resp := responseError{
					Message: "invalid value of 'complete' query parameter",
				}
				encode(w, http.StatusBadRequest, resp)
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

			parts, err := queries.FindUploadPartsById(context.Background(), id)
			if err == sql.ErrNoRows {
				resp := responseError{
					Message: fmt.Sprintf("no parts uploaded for upload with ID %v", id),
				}
				encode(w, http.StatusBadRequest, resp)
				return
			}

			if parts[len(parts)-1] != int64(len(parts)) {
				resp := responseError{
					Message: fmt.Sprintf("missing parts for upload with ID %v", id),
				}
				encode(w, http.StatusBadRequest, resp)
				return
			}

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
			encode(w, http.StatusOK, resp)
		},
	)
}
