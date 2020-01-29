package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/markbates/pkger"
)

const (
	ImageName    = "binny"
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

	err = ExecDocker(ImageName, os.Args[1:])
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

func ExecDocker(imageName string, imageArgs []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	args := []string{
		"run",
		"--rm",
		"-it",
		"-v", fmt.Sprintf("%v:%v", cwd, WorkingDir),
		"-w", WorkingDir,
		imageName,
	}
	return Exec("docker", append(args, imageArgs...))
}

func Exec(name string, args []string) error {
	cmd, err := exec.LookPath(name)
	if err != nil {
		return err
	}

	env := os.Environ()
	args = append([]string{cmd}, args...)
	return syscall.Exec(cmd, args, env)
}
