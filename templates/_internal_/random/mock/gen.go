package mockrandom

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/random/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"mock_rand.go": mockRandDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed mock_rand.gotpl
var mockRandTemplate string

func mockRandDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockRandTemplate, nil)
}
