// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

var assets = http.Dir("assets")

func main() {
	err := vfsgen.Generate(assets, vfsgen.Options{
		Filename: "assets.go",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
