package server

import "net/http"

func handleHealth() http.Handler {
	type response struct {
		Status string `json:"status"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			response := response{
				Status: "up",
			}
			encode(w, 200, response)
		},
	)
}
