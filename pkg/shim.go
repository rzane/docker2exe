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
		return errors.Wrap(LoadEmbedded(config.Load), "docker load failed")
	}

	if config.Build != "" {
		return errors.Wrap(Build(config), "docker build failed")
	}

	return errors.Wrap(Pull(config.Image), "docker pull failed")
}
