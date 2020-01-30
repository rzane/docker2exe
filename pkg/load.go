package binny

import (
	"io"
	"os"
	"os/exec"

	"github.com/markbates/pkger"
)

func Load(filename string) error {
	file, err := pkger.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd := exec.Command("docker", "load")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
