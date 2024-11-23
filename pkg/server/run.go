package server

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	embed "github.com/johngerving/uploader"
	"github.com/johngerving/uploader/repository"
)

type args struct {
	database string
}

// NewArgs returns an instance of an args struct.
// It returns an error if the args are invalid.
func NewArgs(database string) (args, error) {
	database = strings.ToLower(database)

	// Validate args
	switch {
	case database == "":
		database = "disk"
	case database != "disk" && database != "memory":
		return args{}, fmt.Errorf("database type '%v' invalid - must be 'disk' or 'memory'", database)
	}

	args := args{
		database: database,
	}

	return args, nil
}

func Run(
	ctx context.Context,
	args args,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := log.New(stdout, "", log.LstdFlags)

	var db *sql.DB
	var err error
	// If the args specified an in-memory database, create a
	// SQLite database in memory. Otherwise, create it on the
	// disk.
	if args.database == "memory" {
		db, err = sql.Open("sqlite3", ":memory:")
	} else {
		db, err = sql.Open("sqlite3", "./uploads.db")
	}

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	if _, err := db.ExecContext(context.Background(), embed.DBSchema); err != nil {
		return err
	}

	queries := repository.New(db)

	srv := NewServer(
		logger,
		queries,
	)

	host := "localhost"
	port := "8090"

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: srv,
	}

	go func() {
		fmt.Fprintf(stdout, "Listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "Error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		fmt.Fprintf(stdout, "Shutting down http server\n")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stdout, "Error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}
