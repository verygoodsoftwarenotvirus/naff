package auth

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/auth/doc.go":               docDotGo(),
		"services/v1/auth/http_routes_test.go":  httpRoutesTestDotGo(pkg),
		"services/v1/auth/auth_service.go":      authServiceDotGo(pkg),
		"services/v1/auth/auth_service_test.go": authServiceTestDotGo(pkg),
		"services/v1/auth/middleware_test.go":   middlewareTestDotGo(pkg),
		"services/v1/auth/mock_test.go":         mockTestDotGo(pkg),
		"services/v1/auth/wire.go":              wireDotGo(pkg),
		"services/v1/auth/http_routes.go":       httpRoutesDotGo(pkg),
		"services/v1/auth/middleware.go":        middlewareDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
