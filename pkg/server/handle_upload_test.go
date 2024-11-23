package server

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

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

	client := http.Client{}
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8090/uploads",
		nil,
	)

	if err != nil {
		t.Fatal("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %s\n", err.Error())
	}
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
}
