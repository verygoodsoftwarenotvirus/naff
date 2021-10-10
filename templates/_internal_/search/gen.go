package search

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/search"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"doc.go":         docDotGo(proj),
		"config.go":      configDotGo(proj),
		"config_test.go": configTestDotGo(proj),
		"errors.go":      errorsDotGo(proj),
		"search.go":      searchDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
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

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed errors.gotpl
var errorsTemplate string

func errorsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, errorsTemplate, nil)
}

//go:embed search.gotpl
var searchTemplate string

func searchDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, searchTemplate, nil)
}
