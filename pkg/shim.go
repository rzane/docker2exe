package binny

import "os/exec"

func Shim(config Config) error {
	if !isExistingImage(config.Image) {
		if config.Load != "" {
			err := Load(config.Load)
			if err != nil {
				return err
			}
		}

		if config.Build != "" {
			err := Build(config)
			if err != nil {
				return err
			}
		}
	}

	return Exec(config)
}

func isExistingImage(image string) bool {
	return exec.Command("docker", "inspect", image).Run() == nil
}
