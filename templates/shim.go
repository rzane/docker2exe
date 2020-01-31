package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/mattn/go-isatty"
)

type Shim struct {
	Image   string
	Workdir string
	Env     []string
	Volumes []string
	Stdout  io.Writer
	Stderr  io.Writer
}

func (shim *Shim) Exists() bool {
	cmd := exec.Command("docker", "inspect", shim.Image)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (shim *Shim) Pull() error {
	cmd := exec.Command("docker", "pull", shim.Image)
	cmd.Stdout = shim.Stdout
	cmd.Stderr = shim.Stderr
	return cmd.Run()
}

func (shim *Shim) Build(context string) error {
	cmd := exec.Command("docker", "build", "-t", shim.Image, context)
	cmd.Stdout = shim.Stdout
	cmd.Stderr = shim.Stderr
	return cmd.Run()
}

func (shim *Shim) Load(file io.Reader) error {
	cmd := exec.Command("docker", "load")
	cmd.Stdout = shim.Stdout
	cmd.Stderr = shim.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	_, err = io.Copy(stdin, file)
	if err != nil {
		return err
	}

	err = stdin.Close()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (shim *Shim) Exec(containerArgs []string) error {
	cmd, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	args, err := shim.assembleRunArgs()
	if err != nil {
		return err
	}

	args = append([]string{cmd}, args...)
	args = append(args, containerArgs...)
	return syscall.Exec(cmd, args, os.Environ())
}

func (shim *Shim) assembleRunArgs() ([]string, error) {
	args := []string{"run", "--rm"}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		args = append(args, "-it")
	}

	for _, env := range shim.Env {
		args = append(args, "-e", env)
	}

	for _, volume := range shim.Volumes {
		args = append(args, "-v", volume)
	}

	if shim.Workdir != "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		args = append(args, "-w", shim.Workdir)
		args = append(args, "-v", fmt.Sprintf("%s:%s", cwd, shim.Workdir))
	}

	return append(args, shim.Image), nil
}
