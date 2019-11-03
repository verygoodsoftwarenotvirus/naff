package users

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/users/middleware.go":         middlewareDotGo(pkgRoot, types),
		"services/v1/users/middleware_test.go":    middlewareTestDotGo(pkgRoot, types),
		"services/v1/users/users_service.go":      usersServiceDotGo(pkgRoot, types),
		"services/v1/users/users_service_test.go": usersServiceTestDotGo(pkgRoot, types),
		"services/v1/users/wire.go":               wireDotGo(pkgRoot, types),
		"services/v1/users/doc.go":                docDotGo(),
		"services/v1/users/http_routes.go":        httpRoutesDotGo(pkgRoot, types),
		"services/v1/users/http_routes_test.go":   httpRoutesTestDotGo(pkgRoot, types),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
