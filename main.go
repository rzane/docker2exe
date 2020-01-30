package main

import (
	"os"

	binny "github.com/rzane/binny/pkg"
)

func main() {
	config := binny.Config{
		Image:   "binny",
		Workdir: "/workdir",
		Args:    os.Args[1:],
		Env:     []string{},
	}

	err := binny.Load("/image.tar.gz")
	if err != nil {
		panic(err)
	}

	err = binny.Exec(config)
	if err != nil {
		panic(err)
	}
}
