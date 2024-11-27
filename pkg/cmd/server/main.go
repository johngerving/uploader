package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/johngerving/uploader/pkg/server"
)

func main() {
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
