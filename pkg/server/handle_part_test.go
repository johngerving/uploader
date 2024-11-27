package server

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestHandlePart(t *testing.T) {
	type request struct {
		suffix string
		body   []byte
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
					"/1",
					[]byte("hello world"),
				},
				{
					"/2",
					[]byte("testing"),
				},
			},
			[]int{http.StatusOK, http.StatusOK},
			[]errorResp{
				{
					"^$",
				},
				{
					"^$",
				},
			},
		},
		{
			"Invalid part ID should result in a response error",
			true,
			[]request{
				{
					"/bad",
					[]byte("hello world"),
				},
			},
			[]int{http.StatusBadRequest},
			[]errorResp{
				{
					Message: "invalid 'part' path value",
				},
			},
		},
		{
			"Invalid upload ID should return error response",
			false,
			[]request{
				{
					"/1",
					[]byte("hello world"),
				},
			},
			[]int{http.StatusNotFound},
			[]errorResp{
				{
					"^$",
				},
			},
		},
		{
			"Uploading duplicate part IDs should result in an error response",
			true,
			[]request{
				{
					"/1",
					[]byte("hello world"),
				},
				{
					"/1",
					[]byte("test"),
				},
			},
			[]int{http.StatusOK, http.StatusBadRequest},
			[]errorResp{
				{
					"^$",
				},
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
				hash := md5.Sum(req.body)

				endpoint := "http://localhost:8090/uploads/" + uploadId + "/parts" + req.suffix
				resp, err = client.R().
					EnableTrace().
					SetBody(req.body).
					SetHeader("Content-MD5", hex.EncodeToString(hash[:])).
					Post(endpoint)

				if err != nil {
					t.Fatal(err)
				}

				var respBody errorResp
				json.Unmarshal(resp.Body(), &respBody)
				if tt.wantRespBodies[i].Message != "" {
					require.Regexp(t, tt.wantRespBodies[i].Message, respBody.Message)
				}

				require.Equal(t, tt.wantStatusCodes[i], resp.StatusCode())
			}
		})
	}
}

func TestHandlePartHeaders(t *testing.T) {
	tests := []struct {
		description    string
		body           []byte
		md5Header      string
		wantStatusCode int
		wantRespBody   string
	}{
		{
			"Normal request should return successfully",
			[]byte("hello world"),
			"5eb63bbbe01eeed093cb22bb8f5acdc3",
			http.StatusOK,
			"^$",
		},
		{
			"Missing Content-MD5 header should result in an error response",
			[]byte("hello world"),
			"",
			http.StatusBadRequest,
			`{"message":"missing Content-MD5 header"}`,
		},
		{
			"Incorrect Content-MD5 header should result in an error response",
			[]byte("hello world"),
			"5eb63bbbe01eeed093cb22bb8f5acdc2",
			http.StatusBadRequest,
			`{"message":"Content-MD5 header does not match body MD5"}`,
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

			uploadId := upload.ID

			endpoint := "http://localhost:8090/uploads/" + uploadId + "/parts/1"
			resp, err = client.R().
				EnableTrace().
				SetBody(tt.body).
				SetHeader("Content-MD5", tt.md5Header).
				Post(endpoint)

			if err != nil {
				t.Fatal(err)
			}

			require.Equal(t, tt.wantStatusCode, resp.StatusCode())
			require.Regexp(t, tt.wantRespBody, string(resp.Body()))
		})
	}
}
