package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rakyll/statik/fs"
	_ "github.com/rzane/binny/statik"
	binny "github.com/rzane/binny/pkg"
)

func main() {
	config := binny.Config{
		Image:   "binny",
		Load:    true,
		Workdir: "/workdir",
		Args:    os.Args[1:],
		Env:     []string{},
		Open: func() (http.File, error) {
			statik, err := fs.New()
			if err != nil {
				return nil, err
			}

			return statik.Open("/image.tar.gz")
		},
	}

	if err := binny.Shim(config); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
