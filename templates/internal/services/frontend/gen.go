package frontend

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "frontend"

	basePackagePrefix = "internal/services/frontend"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"wire.go":                  wireDotGo(proj),
		"doc.go":                   docDotGo(),
		"frontend_service.go":      frontendServiceDotGo(proj),
		"frontend_service_test.go": frontendServiceTestDotGo(proj),
		"http_routes.go":           httpRoutesDotGo(proj),
		"http_routes_test.go":      httpRoutesTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePrefix, path), file); err != nil {
			return err
		}
	}

	return nil
}
