package binny

import (
	"io"
	"os"
	"os/exec"

	"github.com/markbates/pkger"
	"github.com/pkg/errors"
)

func Load(file io.Reader) error {
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

func LoadEmbedded(filename string) error {
	file, err := pkger.Open(filename)
	if err != nil {
		return errors.Wrap(err, "open image failed")
	}
	defer file.Close()

	if err := Load(file); err != nil {
		return errors.Wrap(err, "load failed")
	}
	return nil
}
