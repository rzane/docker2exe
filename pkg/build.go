package binny

import "os/exec"

func Build(config Config) error {
	cmd := exec.Command("docker", "build", "-t", config.Image, config.Build)
	return cmd.Run()
}
