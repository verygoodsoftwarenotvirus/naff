package oauth2clients

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/oauth2clients/wire.go":                       wireDotGo(pkg),
		"services/v1/oauth2clients/http_routes.go":                httpRoutesDotGo(pkg),
		"services/v1/oauth2clients/implementation_test.go":        implementationTestDotGo(pkg),
		"services/v1/oauth2clients/middleware.go":                 middlewareDotGo(pkg),
		"services/v1/oauth2clients/oauth2_handler_mock_test.go":   oauth2HandlerMockTestDotGo(pkg),
		"services/v1/oauth2clients/oauth2clients_service.go":      oauth2ClientsServiceDotGo(pkg),
		"services/v1/oauth2clients/oauth2clients_service_test.go": oauth2ClientsServiceTestDotGo(pkg),
		"services/v1/oauth2clients/doc.go":                        docDotGo(),
		"services/v1/oauth2clients/http_routes_test.go":           httpRoutesTestDotGo(pkg),
		"services/v1/oauth2clients/implementation.go":             implementationDotGo(pkg),
		"services/v1/oauth2clients/middleware_test.go":            middlewareTestDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg, path, file); err != nil {
			return err
		}
	}

	return nil
}
