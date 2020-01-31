package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
)

var templates = template.Must(template.ParseGlob("gen/templates/*"))

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

	if err := writeMain(tmp, opts); err != nil {
		return err
	}

	if opts.Embed {
		if err := prepareEmbed(tmp, opts); err != nil {
			return err
		}
	}

	fmt.Println(tmp)

	return nil
}

func writeMain(dest string, opts Options) error {
	return writeTemplate("main.go.tmpl", path.Join(dest, "main.go"), opts)
}

func writeTemplate(name string, dest string, opts Options) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	return templates.ExecuteTemplate(file, name, opts)
}

func prepareEmbed(tmp string, opts Options) error {
	assetsDir := path.Join(tmp, "assets")
	imageTarball := path.Join(assetsDir, "image.tar.gz")

	if err := os.Mkdir(assetsDir, 0755); err != nil {
		return err
	}

	if err := saveImage(opts.Image, imageTarball); err != nil {
		return err
	}

	if err := createStatik(tmp, assetsDir); err != nil {
		return err
	}

	return nil
}

func createStatik(tmp string, assetsDir string) error {
	cmd := exec.Command("statik", "-src="+assetsDir, "-dest="+tmp, "-Z")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func saveImage(image string, tarball string) error {
	cmdStr := fmt.Sprintf("docker save \"%s\" | gzip > \"%s\"", image, tarball)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
