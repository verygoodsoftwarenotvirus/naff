package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/frontend/wire.go":                  wireDotGo(proj),
		"services/v1/frontend/doc.go":                   docDotGo(),
		"services/v1/frontend/frontend_service.go":      frontendServiceDotGo(proj),
		"services/v1/frontend/frontend_service_test.go": frontendServiceTestDotGo(proj),
		"services/v1/frontend/http_routes.go":           httpRoutesDotGo(proj),
		"services/v1/frontend/http_routes_test.go":      httpRoutesTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
