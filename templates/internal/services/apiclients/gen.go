package apiclients

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "apiclients"

	basePackagePath = "internal/services/apiclients"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"wire.go":                       wireDotGo(proj),
		"http_routes.go":                httpRoutesDotGo(proj),
		"implementation_test.go":        implementationTestDotGo(proj),
		"middleware.go":                 middlewareDotGo(proj),
		"oauth2_handler_mock_test.go":   oauth2HandlerMockTestDotGo(proj),
		"oauth2clients_service.go":      oauth2ClientsServiceDotGo(proj),
		"oauth2clients_service_test.go": oauth2ClientsServiceTestDotGo(proj),
		"doc.go":                        docDotGo(),
		"http_routes_test.go":           httpRoutesTestDotGo(proj),
		"implementation.go":             implementationDotGo(proj),
		"middleware_test.go":            middlewareTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
