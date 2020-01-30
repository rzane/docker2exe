//go:generate go run assets_generate.go

package main

import (
	"fmt"
	"os"

	_ "github.com/shurcooL/vfsgen"
	binny "github.com/rzane/binny/pkg"
)

func main() {
	config := binny.Config{
		Image:      "binny",
		Load:       "image.tar.gz",
		FileSystem: assets,
		Workdir:    "/workdir",
		Args:       os.Args[1:],
		Env:        []string{},
	}

	if err := binny.Shim(config); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
