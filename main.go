package main

import "fmt"
import "github.com/markbates/pkger"

func main() {
	info, err := pkger.Stat("/image.tar.gz")
	if err != nil {
		panic(err)
	}

	fmt.Println("Name: ", info.Name())
	fmt.Println("Size: ", info.Size())
	fmt.Println("Mode: ", info.Mode())
	fmt.Println("ModTime: ", info.ModTime())
}
