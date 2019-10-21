package users

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/users/middleware.go":         middlewareDotGo(),
		"services/v1/users/middleware_test.go":    middlewareTestDotGo(),
		"services/v1/users/users_service.go":      usersServiceDotGo(),
		"services/v1/users/users_service_test.go": usersServiceTestDotGo(),
		"services/v1/users/wire.go":               wireDotGo(),
		"services/v1/users/doc.go":                docDotGo(),
		"services/v1/users/http_routes.go":        httpRoutesDotGo(),
		"services/v1/users/http_routes_test.go":   httpRoutesTestDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
