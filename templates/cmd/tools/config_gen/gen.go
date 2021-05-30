package config_gen

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "main"

	basePackagePath = "cmd/tools/config_gen"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"main.go": mainDotGoString(proj),
		"doc.go":  docDotGoString(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed main.gotpl
var mainTemplate string

func mainDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, mainTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}
