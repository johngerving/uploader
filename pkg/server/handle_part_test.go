package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestHandlePart(t *testing.T) {
	type request struct {
		queryParam string
		body       []byte
	}
	type errorResp struct {
		Message string `json:"message"`
	}
	tests := []struct {
		description     string
		createUpload    bool // Whether to create an upload beforehand or not
		requests        []request
		wantStatusCodes []int
		wantRespBodies  []errorResp
	}{
		{
			"Normal request should return successfully",
			true,
			[]request{
				{
					"?part=1",
					[]byte("hello world"),
				},
				{
					"?part=2",
					[]byte("testing"),
				},
			},
			[]int{http.StatusOK, http.StatusOK},
			[]errorResp{{}, {}},
		},
		{
			"Missing query params should result in a response error",
			true,
			[]request{
				{
					"",
					[]byte("hello world"),
				},
			},
			[]int{http.StatusBadRequest},
			[]errorResp{
				{
					Message: "missing query parameter 'part'",
				},
			},
		},
		{
			"Invalid query params should result in a response error",
			true,
			[]request{
				{
					"?part=bad",
					[]byte("hello world"),
				},
			},
			[]int{http.StatusBadRequest},
			[]errorResp{
				{
					Message: "invalid query parameter 'part'",
				},
			},
		},
		{
			"Invalid upload ID should return error response",
			false,
			[]request{
				{
					"?part=1",
					[]byte("hello world"),
				},
			},
			[]int{http.StatusNotFound},
			[]errorResp{{}},
		},
		{
			"Uploading duplicate part IDs should result in an error response",
			true,
			[]request{
				{
					"?part=1",
					[]byte("hello world"),
				},
				{
					"?part=1",
					[]byte("test"),
				},
			},
			[]int{http.StatusOK, http.StatusBadRequest},
			[]errorResp{
				{},
				{
					"part 1 already exists for upload with ID .*",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			t.Cleanup(cancel)

			var stdout bytes.Buffer
			args, err := NewArgs("memory")
			if err != nil {
				t.Fatal(err)
			}
			go Run(ctx, args, nil, &stdout, &stdout)

			err = waitForReady(ctx, time.Duration(time.Second*5), "http://localhost:8090/healthz")
			if err != nil {
				t.Fatal(err)
			}

			client := resty.New()

			var resp *resty.Response
			uploadId := "testuploadid"

			if tt.createUpload {
				resp, err = client.R().
					Post("http://localhost:8090/uploads")

				if resp.StatusCode() != http.StatusOK {
					t.Fatal(err)
				}

				type uploadResp struct {
					ID string `json:"id"`
				}
				var upload uploadResp
				err = json.Unmarshal(resp.Body(), &upload)
				if err != nil {
					t.Fatal(err)
				}

				uploadId = upload.ID
			}

			for i, req := range tt.requests {
				endpoint := "http://localhost:8090/uploads/" + uploadId + "/parts" + req.queryParam
				resp, err = client.R().
					EnableTrace().
					SetBody(req.body).
					Post(endpoint)

				if err != nil {
					t.Fatal(err)
				}

				var respBody errorResp

				if len(resp.Body()) > 0 {
					err = json.Unmarshal(resp.Body(), &respBody)
					if err != nil {
						t.Fatal(err)
					}
				}

				fmt.Println("i:", i)
				fmt.Println("want:", tt.wantRespBodies[i])
				fmt.Println("got:", respBody)

				require.Equal(t, tt.wantStatusCodes[i], resp.StatusCode())
				require.Regexp(t, tt.wantRespBodies[i], respBody)
			}
		})
	}
}
