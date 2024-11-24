package server

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestHandleUpload(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	var stdout bytes.Buffer
	args, err := NewArgs("memory")
	if err != nil {
		t.Fatal(err)
	}
	go Run(ctx, args, nil, &stdout, &stdout)

	waitForReady(ctx, time.Duration(time.Second*5), "http://localhost:8090/healthz")

	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Post("http://localhost:8090/uploads")

	require.Equal(t, nil, err)
	require.Equal(t, 200, resp.StatusCode())

	responseStruct := struct {
		ID string `json:"id"`
	}{}
	json.Unmarshal(resp.Body(), &responseStruct)

	resp, err = client.R().
		EnableTrace().
		Get("http://localhost:8090/uploads/" + responseStruct.ID)

	require.Equal(t, nil, err)
	require.Equal(t, 200, resp.StatusCode())
}
