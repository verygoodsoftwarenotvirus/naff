package oauth2clients

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "oauth2clients"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/oauth2clients/wire.go":                       wireDotGo(proj),
		"services/v1/oauth2clients/http_routes.go":                httpRoutesDotGo(proj),
		"services/v1/oauth2clients/implementation_test.go":        implementationTestDotGo(proj),
		"services/v1/oauth2clients/middleware.go":                 middlewareDotGo(proj),
		"services/v1/oauth2clients/oauth2_handler_mock_test.go":   oauth2HandlerMockTestDotGo(proj),
		"services/v1/oauth2clients/oauth2clients_service.go":      oauth2ClientsServiceDotGo(proj),
		"services/v1/oauth2clients/oauth2clients_service_test.go": oauth2ClientsServiceTestDotGo(proj),
		"services/v1/oauth2clients/doc.go":                        docDotGo(),
		"services/v1/oauth2clients/http_routes_test.go":           httpRoutesTestDotGo(proj),
		"services/v1/oauth2clients/implementation.go":             implementationDotGo(proj),
		"services/v1/oauth2clients/middleware_test.go":            middlewareTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
