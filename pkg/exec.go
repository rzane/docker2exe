package binny

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/mattn/go-isatty"
)

type ExecOptions struct {
	Image   string
	Args    []string
	Env     []string
	Workdir string
}

func Exec(opts ExecOptions) error {
	cmd, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	args, err := assembleExecArgs(cmd, opts)
	if err != nil {
		return err
	}

	return syscall.Exec(cmd, args, os.Environ())
}

func assembleExecArgs(cmd string, opts ExecOptions) ([]string, error) {
	args := []string{cmd, "run", "--rm"}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		args = append(args, "-it")
	}

	for _, env := range opts.Env {
		args = append(args, "-e", env)
	}

	if opts.Workdir != "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		args = append(args, "-w", opts.Workdir)
		args = append(args, "-v", fmt.Sprintf("%s:%s", cwd, opts.Workdir))
	}

	args = append(args, opts.Image)
	args = append(args, opts.Args...)

	return args, nil
}
