package main

import (
	"fmt"
	"os"

	"github.com/rakyll/statik/fs"
	binny "github.com/rzane/binny/pkg"
	_ "github.com/rzane/binny/statik"
)

func main() {
	statik, err := fs.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	config := binny.Config{
		Image:      "binny",
		Load:       "/image.tar.gz",
		FileSystem: statik,
		Workdir:    "/workdir",
		Args:       os.Args[1:],
		Env:        []string{},
	}

	if err := binny.Shim(config); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
