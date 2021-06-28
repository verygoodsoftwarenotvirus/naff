package panicking

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/panicking"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"mock.go":          mockDotGo(proj),
		"panicker.go":      panickerDotGo(proj),
		"standard.go":      standardDotGo(proj),
		"standard_test.go": standardTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed mock.gotpl
var mockTemplate string

func mockDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockTemplate, nil)
}

//go:embed panicker.gotpl
var panickerTemplate string

func panickerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, panickerTemplate, nil)
}

//go:embed standard.gotpl
var standardTemplate string

func standardDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, standardTemplate, nil)
}

//go:embed standard_test.gotpl
var standardTestTemplate string

func standardTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, standardTestTemplate, nil)
}
