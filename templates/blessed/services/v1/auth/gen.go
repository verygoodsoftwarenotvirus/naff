package auth

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/auth/doc.go":               docDotGo(),
		"services/v1/auth/http_routes_test.go":  httpRoutesTestDotGo(pkgRoot),
		"services/v1/auth/auth_service.go":      authServiceDotGo(pkgRoot),
		"services/v1/auth/auth_service_test.go": authServiceTestDotGo(pkgRoot),
		"services/v1/auth/middleware_test.go":   middlewareTestDotGo(pkgRoot),
		"services/v1/auth/mock_test.go":         mockTestDotGo(pkgRoot),
		"services/v1/auth/wire.go":              wireDotGo(pkgRoot),
		"services/v1/auth/http_routes.go":       httpRoutesDotGo(pkgRoot),
		"services/v1/auth/middleware.go":        middlewareDotGo(pkgRoot),
	}

	for path, file := range files {
		if err := utils.RenderFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
