package converters

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "pkg/types/converters"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"converters_test.go": convertersTestDotGo(proj),
		"converters.go":      convertersDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed converters_test.gotpl
var convertersTestTemplate string

func convertersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, convertersTestTemplate, nil)
}

//go:embed converters.gotpl
var convertersTemplate string

func convertersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, convertersTemplate, nil)
}
