package users

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "users"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/users/middleware.go":         middlewareDotGo(proj),
		"services/v1/users/middleware_test.go":    middlewareTestDotGo(proj),
		"services/v1/users/users_service.go":      usersServiceDotGo(proj),
		"services/v1/users/users_service_test.go": usersServiceTestDotGo(proj),
		"services/v1/users/wire.go":               wireDotGo(proj),
		"services/v1/users/doc.go":                docDotGo(),
		"services/v1/users/rand.go":               randDotGo(),
		"services/v1/users/rand_test.go":          randTestDotGo(),
		"services/v1/users/http_routes.go":        httpRoutesDotGo(proj),
		"services/v1/users/http_routes_test.go":   httpRoutesTestDotGo(proj),
		"services/v1/users/wire_test.go":          wireTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
