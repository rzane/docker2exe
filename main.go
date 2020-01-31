package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/rzane/binny/gen"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "binny",
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
				Name:  "build",
				Usage: "build a docker image upon install",
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
			&cli.StringFlag{
				Name:  "module",
				Usage: "name of generated golang module",
			},
		},
	}

	app.Run(os.Args)
}

func generate(c *cli.Context) error {
	opts := gen.Options{
		Output:  c.String("output"),
		Module:  c.String("module"),
		Name:    c.String("name"),
		Image:   c.String("image"),
		Build:   c.String("build"),
		Embed:   c.Bool("embed"),
		Workdir: c.String("workdir"),
		Env:     c.StringSlice("env"),
		Volumes: c.StringSlice("volume"),
	}

	if opts.Output == "" {
		cwd, _ := os.Getwd()
		opts.Output = cwd
	}

	if opts.Module == "" {
		user, _ := user.Current()
		opts.Module = fmt.Sprintf("github.com/%s/%s", user.Username, opts.Name)
	}

	return gen.Generate(opts)
}
