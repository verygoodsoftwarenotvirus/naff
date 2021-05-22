package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"path/filepath"
)

const (
	packageName = "auth"

	basePackagePrefix = "internal/services/authentication"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":               docDotGo(),
		"http_routes_test.go":  httpRoutesTestDotGo(proj),
		"auth_service.go":      authServiceDotGo(proj),
		"auth_service_test.go": authServiceTestDotGo(proj),
		"middleware_test.go":   middlewareTestDotGo(proj),
		"mock_test.go":         mockTestDotGo(proj),
		"wire.go":              wireDotGo(proj),
		"wire_test.go":         wireTestDotGo(proj),
		"http_routes.go":       httpRoutesDotGo(proj),
		"middleware.go":        middlewareDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePrefix, path), file); err != nil {
			return err
		}
	}

	return nil
}
