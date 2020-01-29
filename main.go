package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/markbates/pkger"
)

func main() {
	info, err := pkger.Stat("/image.tar")
	if err != nil {
		panic(err)
	}
	fmt.Printf("total bytes: %v\n", info.Size())

	file, err := pkger.Open("/image.tar")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := loadImage(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("loaded bytes: %v\n", bytes)
}

func loadImage(input io.Reader) (int64, error) {
	cmd := exec.Command("docker", "load")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return 0, err
	}

	err = cmd.Start()
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(stdin, input)
	if err != nil {
		return bytes, err
	}

	err = stdin.Close()
	if err != nil {
		return bytes, err
	}

	return bytes, cmd.Wait()
}
