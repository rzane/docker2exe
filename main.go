package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/markbates/pkger"
	"github.com/rzane/binny/pkg"
)

const (
	ImageTarball = "/image.tar.gz"
	WorkingDir   = "/workdir"
)

func main() {
	info, err := pkger.Stat(ImageTarball)
	if err != nil {
		panic(err)
	}
	fmt.Printf("total bytes: %v\n", info.Size())

	file, err := pkger.Open(ImageTarball)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := LoadImage(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("loaded bytes: %v\n", bytes)

	opts := binny.ExecOptions{
		Image:   "binny",
		Workdir: "/workdir",
		Args:    os.Args[1:],
	}

	err = binny.Exec(opts)
	if err != nil {
		panic(err)
	}
}

func LoadImage(input io.Reader) (int64, error) {
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
