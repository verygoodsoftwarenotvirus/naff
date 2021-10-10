package frontend

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "tests/frontend"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"auth_test.go": authTestDotGo(proj),
		"doc.go":       docDotGo(proj),
		"helpers.go":   helpersDotGo(proj),
		"init.go":      initDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed auth_test.gotpl
var authTestTemplate string

func authTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed helpers.gotpl
var helpersTemplate string

func helpersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTemplate, nil)
}

//go:embed init.gotpl
var initTemplate string

func initDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, initTemplate, nil)
}
