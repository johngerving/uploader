package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	var stdout bytes.Buffer
	go run(ctx, []string{}, nil, &stdout, &stdout)

	waitForReady(ctx, time.Duration(time.Second * 5), "http://localhost:8090/healthz")

	cancel()

	want := "Listening on localhost:8090\nShutting down http server\n"
	got := stdout.String()

	// Wait for the server to shut down for 10 seconds
	timeout := time.Duration(time.Second * 10)
	startTime := time.Now()
	for got != want {
		if time.Since(startTime) >= timeout {
			t.Fatalf("Timeout reached while waiting for shutdown")
		}
		time.Sleep(250 * time.Millisecond)

		got = stdout.String()
	}

	
	require.Equal(t, want, got)
}

// waitForReady calls the specified endpoint until it gets a 200
// response or until the context is cancelled or the timeout is
// reached.
func waitForReady(
	ctx context.Context,
	timeout time.Duration, 
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)

		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
		}
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}