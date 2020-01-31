package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*"))

type Options struct {
	Module  string
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

	if err := initializeModule(tmp, opts.Module); err != nil {
		return err
	}

	if err := copyTemplates(tmp, opts); err != nil {
		return err
	}

	if opts.Embed {
		if err := prepareEmbed(tmp, opts); err != nil {
			return err
		}
	}

	if err := installDependencies(tmp); err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	output := path.Join(cwd, opts.Name)
	if err := compileModule(tmp, output); err != nil {
		return err
	}

	fmt.Println(output)
	return os.RemoveAll(tmp)
}

func compileModule(cwd string, output string) error {
	fmt.Println("==>> Building binary...")
	cmd := exec.Command("go", "build", "-o", output)
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func initializeModule(cwd string, moduleName string) error {
	fmt.Println("==>> Initializing module...")
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installDependencies(cwd string) error {
	fmt.Println("==>> Installing dependencies...")
	cmd := exec.Command("go", "get", "-v")
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyTemplates(dest string, opts Options) error {
	fmt.Println("==>> Running the code generator...")
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

func prepareEmbed(cwd string, opts Options) error {
	assetsDir := path.Join(cwd, "assets")
	imageTarball := path.Join(assetsDir, "image.tar.gz")

	if err := os.Mkdir(assetsDir, 0755); err != nil {
		return err
	}

	if err := saveImage(opts.Image, imageTarball); err != nil {
		return err
	}

	if err := createStatik(assetsDir, cwd); err != nil {
		return err
	}

	return nil
}

func createStatik(src string, dest string) error {
	fmt.Println("==>> Embedding the docker image into the binary...")
	cmd := exec.Command("statik", "-src="+src, "-dest="+dest, "-Z")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func saveImage(image string, tarball string) error {
	fmt.Println("==>> Exporting the docker image as a tarball...")
	cmdStr := fmt.Sprintf("docker save \"%s\" | gzip > \"%s\"", image, tarball)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
