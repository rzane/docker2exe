package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/rzane/docker2exe/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "docker2exe",
		Usage:  "create an executable from a docker image",
		Action: generate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required: true,
				Name:     "name",
				Usage:    "name of your executable",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "image",
				Usage:    "name of your docker image",
			},
			&cli.BoolFlag{
				Name:  "embed",
				Usage: "embed a docker image in the binary",
			},
			&cli.StringFlag{
				Name:    "workdir",
				Aliases: []string{"w"},
				Usage:   "mount the user's current directory in the image",
			},
			&cli.StringSliceFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "whitelist environment variables",
			},
			&cli.StringSliceFlag{
				Name:    "volume",
				Aliases: []string{"v"},
				Usage:   "bind mount a volume",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "directory to output",
			},
			&cli.StringSliceFlag{
				Name:    "target",
				Aliases: []string{"t"},
				Usage:   "platforms and architectures to target",
			},
			&cli.StringFlag{
				Name:  "module",
				Usage: "name of generated golang module",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generate(c *cli.Context) error {
	generator := cmd.Generator{
		Output:  c.String("output"),
		Name:    c.String("name"),
		Targets: c.StringSlice("target"),
		Module:  c.String("module"),
		Image:   c.String("image"),
		Embed:   c.Bool("embed"),
		Workdir: c.String("workdir"),
		Env:     c.StringSlice("env"),
		Volumes: c.StringSlice("volume"),
	}

	if generator.Output == "" {
		cwd, _ := os.Getwd()
		generator.Output = path.Join(cwd, "dist")
	}

	if generator.Module == "" {
		user, _ := user.Current()
		generator.Module = fmt.Sprintf("github.com/%s/%s", user.Username, generator.Name)
	}

	if len(generator.Targets) == 0 {
		generator.Targets = []string{"darwin/amd64", "darwin/arm64", "linux/amd64", "windows/amd64"}
	}

	return generator.Run()
}
