package binny

import "os/exec"

func Shim(config Config) error {
	if !isLoaded(config.Image) {
		if err := Load(config.Load); err != nil {
			return err
		}
	}

	return Exec(config)
}

func isLoaded(image string) bool {
	return exec.Command("docker", "inspect", image).Run() == nil
}
