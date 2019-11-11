package users

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/users/middleware.go":         middlewareDotGo(pkg),
		"services/v1/users/middleware_test.go":    middlewareTestDotGo(pkg),
		"services/v1/users/users_service.go":      usersServiceDotGo(pkg),
		"services/v1/users/users_service_test.go": usersServiceTestDotGo(pkg),
		"services/v1/users/wire.go":               wireDotGo(pkg),
		"services/v1/users/doc.go":                docDotGo(),
		"services/v1/users/http_routes.go":        httpRoutesDotGo(pkg),
		"services/v1/users/http_routes_test.go":   httpRoutesTestDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
