package oauth2clients

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/oauth2clients/wire.go":                       wireDotGo(),
		"services/v1/oauth2clients/http_routes.go":                httpRoutesDotGo(),
		"services/v1/oauth2clients/implementation_test.go":        implementationTestDotGo(),
		"services/v1/oauth2clients/middleware.go":                 middlewareDotGo(),
		"services/v1/oauth2clients/oauth2_handler_mock_test.go":   oauth2HandlerMockTestDotGo(),
		"services/v1/oauth2clients/oauth2clients_service.go":      oauth2ClientsServiceDotGo(),
		"services/v1/oauth2clients/oauth2clients_service_test.go": oauth2ClientsServiceTestDotGo(),
		"services/v1/oauth2clients/doc.go":                        docDotGo(),
		"services/v1/oauth2clients/http_routes_test.go":           httpRoutesTestDotGo(),
		"services/v1/oauth2clients/implementation.go":             implementationDotGo(),
		"services/v1/oauth2clients/middleware_test.go":            middlewareTestDotGo(),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
