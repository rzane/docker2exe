package binny

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

func Exec(config Config) error {
	cmd, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	args, err := assembleExecArgs(cmd, config)
	if err != nil {
		return err
	}

	return syscall.Exec(cmd, args, os.Environ())
}

func assembleExecArgs(cmd string, config Config) ([]string, error) {
	args := []string{cmd, "run", "--rm"}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		args = append(args, "-it")
	}

	for _, env := range config.Env {
		args = append(args, "-e", env)
	}

	if config.Workdir != "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		args = append(args, "-w", config.Workdir)
		args = append(args, "-v", fmt.Sprintf("%s:%s", cwd, config.Workdir))
	}

	args = append(args, config.Image)
	args = append(args, config.Args...)

	return args, nil
}
