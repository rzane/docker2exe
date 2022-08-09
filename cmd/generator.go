package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*"))

type Generator struct {
	Name    string
	Output  string
	Targets []string
	Module  string
	Image   string
	Embed   bool
	Workdir string
	Env     []string
	Volumes []string
}

func (gen *Generator) Run() error {
	tmp, err := ioutil.TempDir("", gen.Name)
	if err != nil {
		return err
	}

	if err := gen.copyTemplates(tmp); err != nil {
		return err
	}

	if err := make(tmp); err != nil {
		return err
	}

	return os.RemoveAll(tmp)
}

func (gen *Generator) copyTemplates(dest string) error {
	for _, template := range templates.Templates() {
		filename := strings.TrimRight(template.Name(), ".tmpl")
		filepath := path.Join(dest, filename)

		file, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = template.Execute(file, gen)
		if err != nil {
			return err
		}
	}

	return nil
}

func make(cwd string) error {
	cmd := exec.Command("make")
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
