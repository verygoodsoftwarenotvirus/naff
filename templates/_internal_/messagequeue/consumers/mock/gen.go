package mockconsumers

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "mockconsumers"

	basePackagePath = "internal/messagequeue/consumers/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"mock.go": mockDotGo(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
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
