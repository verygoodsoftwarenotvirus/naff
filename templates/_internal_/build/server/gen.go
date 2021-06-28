package server

import (
	_ "embed"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "server"

	basePackagePrefix = "internal/build/server"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	jenFiles := map[string]*jen.File{
		"build.go": buildDotGo(proj),
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePrefix, path), file); err != nil {
			return err
		}
	}

	stringFiles := map[string]string{
		"doc.go": docDotGo(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePrefix, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}
