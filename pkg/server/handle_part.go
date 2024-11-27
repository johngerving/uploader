package server

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/johngerving/uploader/repository"
	"github.com/mattn/go-sqlite3"
)

func handlePostPart(logger *log.Logger, queries *repository.Queries) http.Handler {
	type responseError struct {
		Message string `json:"message"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			uploadId := r.PathValue("id")
			partString := r.PathValue("part")

			if uploadId == "" {
				resp := responseError{
					Message: "invalid 'id' path value",
				}
				encode(w, http.StatusBadRequest, resp)
				return
			}

			part, err := strconv.Atoi(partString)
			if err != nil || part < 1 {
				resp := responseError{
					Message: "invalid 'part' path value",
				}
				encode(w, http.StatusBadRequest, resp)
				return
			}

			upload, err := queries.FindUploadById(context.Background(), uploadId)
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				logger.Printf("error getting row with ID %v: %v\n", upload.ID, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Printf("error reading body: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if len(body) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			md5Header := r.Header.Get("Content-MD5")
			if md5Header == "" {
				resp := responseError{
					Message: "missing Content-MD5 header",
				}
				encode(w, http.StatusBadRequest, resp)
				return
			}

			hash := md5.Sum(body)

			if hex.EncodeToString(hash[:]) != md5Header {
				resp := responseError{
					Message: "Content-MD5 header does not match body MD5",
				}
				encode(w, http.StatusBadRequest, resp)
				return
			}

			params := repository.CreatePartParams{
				UploadID: uploadId,
				ID:       int64(part),
				Data:     body,
			}
			err = queries.CreatePart(context.Background(), params)

			if err != nil {
				var sqliteErr sqlite3.Error
				if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
					fmt.Println(sqliteErr.Code == sqlite3.ErrConstraint)
					resp := responseError{
						Message: fmt.Sprintf("part %v already exists for upload with ID %v", part, uploadId),
					}
					encode(w, http.StatusBadRequest, resp)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
	)
}
