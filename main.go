package main

import (
	"fmt"
	"os"

	"github.com/markbates/pkger"
	binny "github.com/rzane/binny/pkg"
)

func main() {
	config := binny.Config{
		Image:   "binny",
		Tarball: pkger.Include("/image.tar.gz"),
		Workdir: "/workdir",
		Args:    os.Args[1:],
		Env:     []string{},
	}

	if err := binny.Shim(config); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
