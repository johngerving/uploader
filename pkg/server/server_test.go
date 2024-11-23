package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestRunShutdown(t *testing.T) {
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

	cancel()

	want := "Listening on localhost:8090\nShutting down http server\n"
	got := stdout.String()

	// Wait for the server to shut down for 10 seconds
	timeout := time.Duration(time.Second * 10)
	startTime := time.Now()
	for got != want {
		if time.Since(startTime) >= timeout {
			t.Log("Timeout reached while waiting for shutdown\n")
			break
		}
		time.Sleep(250 * time.Millisecond)

		got = stdout.String()
	}

	require.Equal(t, want, got)
}

func TestNewArgs(t *testing.T) {
	tests := []struct {
		description string
		database    string
		wantArgs    args
		wantErr     error
	}{
		{"Default arguments", "", args{database: "disk"}, nil},
		{"Specifying disk database", "disk", args{database: "disk"}, nil},
		{"Specifying in-memory database", "memory", args{database: "memory"}, nil},
		{"Invalid database argument", "test", args{}, errors.New("database type 'test' invalid - must be 'disk' or 'memory'")},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			gotArgs, gotErr := NewArgs(tt.database)

			require.Equal(t, tt.wantErr, gotErr)
			require.Equal(t, tt.wantArgs, gotArgs)
		})
	}
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
			continue
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
