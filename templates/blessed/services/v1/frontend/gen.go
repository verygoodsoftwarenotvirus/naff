package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/frontend/wire.go":                  wireDotGo(pkg),
		"services/v1/frontend/doc.go":                   docDotGo(),
		"services/v1/frontend/frontend_service.go":      frontendServiceDotGo(pkg),
		"services/v1/frontend/frontend_service_test.go": frontendServiceTestDotGo(pkg),
		"services/v1/frontend/http_routes.go":           httpRoutesDotGo(pkg),
		"services/v1/frontend/http_routes_test.go":      httpRoutesTestDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
