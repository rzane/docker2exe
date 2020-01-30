package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/urfave/cli/v2"
)

type options struct {
	Image   string
	Embed   bool
	Build   string
	Workdir string
	Env     []string
	Volumes []string
}

func main() {
	app := &cli.App{
		Name:   "binny",
		Usage:  "create an executable from a docker image",
		Action: generate,
		Flags: []cli.Flag{
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
		},
	}

	app.Run(os.Args)
}

func generate(c *cli.Context) error {
	opts := parseOptions(c)
	out := mainTemplate()

	tmpl := template.New("main")
	tmpl, err := tmpl.Parse(out)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", opts)
	tmpl.Execute(os.Stdout, opts)
	return nil
}

func parseOptions(c *cli.Context) options {
	return options{
		Image:   c.String("image"),
		Build:   c.String("build"),
		Embed:   c.Bool("embed"),
		Workdir: c.String("workdir"),
		Env:     c.StringSlice("env"),
		Volumes: c.StringSlice("volume"),
	}
}

func mainTemplate() string {
	return `package main

import (
  "os"
  "fmt"
  binny "github.com/rzane/binny/pkg"
)

func main() {
  shim := binny.Shim{
    Image: "{{.Image}}",
    Workdir: "{{.Workdir}}",
  }

  if !shim.Exists() {
    err := shim.Build("{{.Build}}")
    if err != nil {
      fmt.Fprintln(err)
      os.Exit(1)
    }
  }

  err := shim.Exec(os.Args[1:])
  if err != nil {
    fmt.Fprintln(err)
    os.Exit(1)
  }
}
`
}
