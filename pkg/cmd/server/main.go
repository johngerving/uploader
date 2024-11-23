package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/johngerving/uploader/pkg/server"
)

func main() {

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

	dbArg := flag.String("db", "disk", "location of the database, 'disk' or 'memory'")

	// Create args from flags
	args, err := server.NewArgs(*dbArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := server.Run(ctx, args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
