package users

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "users"

	basePackagePath = "internal/services/users"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"middleware.go":         middlewareDotGo(proj),
		"middleware_test.go":    middlewareTestDotGo(proj),
		"users_service.go":      usersServiceDotGo(proj),
		"users_service_test.go": usersServiceTestDotGo(proj),
		"wire.go":               wireDotGo(proj),
		"doc.go":                docDotGo(),
		"rand.go":               randDotGo(),
		"rand_test.go":          randTestDotGo(),
		"http_routes.go":        httpRoutesDotGo(proj),
		"http_routes_test.go":   httpRoutesTestDotGo(proj),
		"wire_test.go":          wireTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
