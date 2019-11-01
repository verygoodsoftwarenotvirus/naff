package oauth2clients

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/oauth2clients/wire.go":                       wireDotGo(pkgRoot),
		"services/v1/oauth2clients/http_routes.go":                httpRoutesDotGo(pkgRoot),
		"services/v1/oauth2clients/implementation_test.go":        implementationTestDotGo(pkgRoot),
		"services/v1/oauth2clients/middleware.go":                 middlewareDotGo(pkgRoot),
		"services/v1/oauth2clients/oauth2_handler_mock_test.go":   oauth2HandlerMockTestDotGo(),
		"services/v1/oauth2clients/oauth2clients_service.go":      oauth2ClientsServiceDotGo(pkgRoot),
		"services/v1/oauth2clients/oauth2clients_service_test.go": oauth2ClientsServiceTestDotGo(pkgRoot),
		"services/v1/oauth2clients/doc.go":                        docDotGo(),
		"services/v1/oauth2clients/http_routes_test.go":           httpRoutesTestDotGo(pkgRoot),
		"services/v1/oauth2clients/implementation.go":             implementationDotGo(pkgRoot),
		"services/v1/oauth2clients/middleware_test.go":            middlewareTestDotGo(pkgRoot),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
