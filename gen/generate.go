package gen

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*"))

type Options struct {
	Name    string
	Output  string
	Targets []string
	Module  string
	Image   string
	Embed   bool
	Build   string
	Workdir string
	Env     []string
	Volumes []string
}

func Generate(opts Options) error {
	tmp, err := ioutil.TempDir("", opts.Name)
	if err != nil {
		return err
	}

	if err := copyTemplates(tmp, opts); err != nil {
		return err
	}

	if err := make(tmp); err != nil {
		return err
	}

	return os.RemoveAll(tmp)
}

func make(cwd string) error {
	cmd := exec.Command("make")
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyTemplates(dest string, opts Options) error {
	for _, template := range templates.Templates() {
		filename := strings.TrimRight(template.Name(), ".tmpl")
		filepath := path.Join(dest, filename)

		file, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = template.Execute(file, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
