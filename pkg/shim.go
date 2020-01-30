package binny

import (
	"os/exec"

	"github.com/pkg/errors"
)

func Shim(config Config) error {
	if err := ensureImageExists(config); err != nil {
		return err
	}

	return Exec(config)
}

func ensureImageExists(config Config) error {
	inspect := exec.Command("docker", "inspect", config.Image)
	if err := inspect.Run(); err == nil {
		return nil
	}

	if config.Load != "" {
		file, err := config.FileSystem.Open(config.Load)
		if err != nil {
			return errors.Wrap(err, "open file failed")
		}
		defer file.Close()

		return errors.Wrap(Load(file), "docker load failed")
	}

	if config.Build != "" {
		return errors.Wrap(Build(config), "docker build failed")
	}

	return errors.Wrap(Pull(config.Image), "docker pull failed")
}
