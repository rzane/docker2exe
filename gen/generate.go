package gen

import (
	"io/ioutil"
	"path"
)

type Options struct {
	Name    string
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

	mainFile := path.Join(tmp, "main.go")
	mainCode := renderMain(opts)
	err = ioutil.WriteFile(mainFile, []byte(mainCode), 0644)
	if err != nil {
		return err
	}

	return nil
}
