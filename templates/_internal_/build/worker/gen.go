package server

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "worker"

	basePackagePrefix = "internal/build/worker"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"build.go": buildDotGo(proj),
		"doc.go":   docDotGo(proj),
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

//go:embed build.gotpl
var buildTemplate string

func buildDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, buildTemplate, nil)
}
