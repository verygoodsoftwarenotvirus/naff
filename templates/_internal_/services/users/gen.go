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
		"service.go":           usersServiceDotGo(proj),
		"service_test.go":      usersServiceTestDotGo(proj),
		"wire.go":              wireDotGo(proj),
		"doc.go":               docDotGo(),
		"http_routes.go":       httpRoutesDotGo(proj),
		"http_routes_test.go":  httpRoutesTestDotGo(proj),
		"http_helpers_test.go": httpHelpersTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
