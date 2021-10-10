package mock

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/routing/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"route_param_manager.go": routeParamManagerDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed route_param_manager.gotpl
var routeParamManagerTemplate string

func routeParamManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, routeParamManagerTemplate, nil)
}
