package random

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/random"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"rand.go":      randDotGo(proj),
		"rand_test.go": randTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed rand.gotpl
var randTemplate string

func randDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, randTemplate, nil)
}

//go:embed rand_test.gotpl
var randTestTemplate string

func randTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, randTestTemplate, nil)
}
