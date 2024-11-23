package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/johngerving/uploader/pkg/server"
	_ "github.com/mattn/go-sqlite3"
)

func run(
	ctx context.Context, 
	args []string,
	stdin io.Reader, 
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv := server.NewServer(nil)

	host := "localhost"
	port := "8090"

	httpServer := &http.Server{
		Addr: net.JoinHostPort(host, port),
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
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10 * time.Second)
		defer cancel()
		fmt.Fprintf(stdout, "Shutting down http server\n")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stdout, "Error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}


func main() {
	// db, err := sql.Open("sqlite3", "./uploads.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// if _, err := db.ExecContext(context.Background(), embed.DBSchema); err != nil {
	// 	log.Fatal(err)
	// }

	// queries := repository.New(db)

	// http.HandleFunc("POST /upload", func(w http.ResponseWriter, req *http.Request) {
	// 	id := uuid.New()

	// 	_, err := queries.CreateUpload(context.Background(), id.String())
	// 	if err != nil {
	// 		log.Printf("error creating upload: %v\n", err)
	// 		fmt.Fprintf(w, "error creating upload")
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "created upload with ID %v", id)
	// })

	// http.ListenAndServe(":8090", nil)

	ctx := context.Background()
	if err := run(ctx, os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}