package auth

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/auth/doc.go":               docDotGo(),
		"services/v1/auth/http_routes_test.go":  httpRoutesTestDotGo(proj),
		"services/v1/auth/auth_service.go":      authServiceDotGo(proj),
		"services/v1/auth/auth_service_test.go": authServiceTestDotGo(proj),
		"services/v1/auth/middleware_test.go":   middlewareTestDotGo(proj),
		"services/v1/auth/mock_test.go":         mockTestDotGo(proj),
		"services/v1/auth/wire.go":              wireDotGo(proj),
		"services/v1/auth/wire_test.go":         wireTestDotGo(proj),
		"services/v1/auth/http_routes.go":       httpRoutesDotGo(proj),
		"services/v1/auth/middleware.go":        middlewareDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
