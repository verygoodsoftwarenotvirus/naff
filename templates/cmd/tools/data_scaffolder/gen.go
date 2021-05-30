package data_scaffolder

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "data_scaffolder"

	basePackagePath = "cmd/tools/data_scaffolder"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"main.go":        mainDotGoString(proj),
		"exiter.go":      exiterDotGoString(proj),
		"exiter_test.go": exiterTestDotGoString(proj),
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

//go:embed exiter.gotpl
var exiterTemplate string

func exiterDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, exiterTemplate, nil)
}

//go:embed exiter_test.gotpl
var exiterTestTemplate string

func exiterTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, exiterTestTemplate, nil)
}
