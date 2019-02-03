package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"gopkg.in/AlecAivazis/survey.v1"
)

// Project is our project type
type Project struct {
	Name        string `survey:"name"`
	BasePackage string `survey:"basePackage"`
}

// ValidateName validates a project name
func (p Project) ValidateName(input string) error {
	if len(input) == 0 {
		return errors.New("empty string")
	}
	return nil
}

// ValidatePath validates an output path
func (p Project) ValidatePath(ans interface{}) error {
	return survey.Required(ans)
}

// RootPath returns what it calculates as the root path for the project.
func (p *Project) RootPath() string {
	gp := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/%s", gp, p.BasePackage)
}

// EnsureRootPath ensures that our base path
func (p *Project) EnsureRootPath() error {
	gp := os.Getenv("GOPATH")
	if gp == "" {
		return errors.New("GOPATH is not set")
	}

	folder := fmt.Sprintf("%s/%s", gp, p.BasePackage)
	fi, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			if mkdirErr := os.MkdirAll(folder, 0777); mkdirErr != nil {
				return mkdirErr
			}
		}
	}

	if !fi.IsDir() {
		return errors.New("destination is not a directory")
	}

	return nil
}

// RenderTemplateToPath renders a template file to a specific path
func renderTemplateToPath(t *template.Template, data interface{}, path string) error {
	p := filepath.Dir(path)
	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return f.Close()
}

// RenderDirectory renders a directory full of templates with the project as the data
func (p *Project) RenderDirectory(directory *packr.Box) error {
	return directory.Walk(
		func(path string, f packd.File) error {
			if strings.HasSuffix(path, defaultFileExtension) {
				b, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}
				t := template.Must(template.New(path).Parse(string(b)))

				renderPath := filepath.Join(p.RootPath(), path)
				if err := renderTemplateToPath(t, p, renderPath); err != nil {
					return err
				}
			}
			return nil
		},
	)
}
