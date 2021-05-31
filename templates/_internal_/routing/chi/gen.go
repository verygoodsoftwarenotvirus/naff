package chi

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/routing/chi"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"chi.go":              chiDotGo(proj),
		"chi_test.go":         chiTestDotGo(proj),
		"routeparams.go":      routeparamsDotGo(proj),
		"routeparams_test.go": routeparamsTestDotGo(proj),
		"wire.go":             wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed chi.gotpl
var chiTemplate string

func chiDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, chiTemplate, nil)
}

//go:embed chi_test.gotpl
var chiTestTemplate string

func chiTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, chiTestTemplate, nil)
}

//go:embed routeparams.gotpl
var routeparamsTemplate string

func routeparamsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, routeparamsTemplate, nil)
}

//go:embed routeparams_test.gotpl
var routeparamsTestTemplate string

func routeparamsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, routeparamsTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}
