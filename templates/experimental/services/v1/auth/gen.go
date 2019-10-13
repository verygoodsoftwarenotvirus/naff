package auth

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/auth/auth_service.go":      authServiceDotGo(),
		"services/v1/auth/auth_service_test.go": authServiceTestDotGo(),
		"services/v1/auth/middleware.go":        middlewareDotGo(),
		"services/v1/auth/mock_test.go":         mockTestDotGo(),
		"services/v1/auth/wire.go":              wireDotGo(),
		"services/v1/auth/doc.go":               docDotGo(),
		"services/v1/auth/http_routes.go":       httpRoutesDotGo(),
		"services/v1/auth/http_routes_test.go":  httpRoutesTestDotGo(),
		"services/v1/auth/middleware_test.go":   middlewareTestDotGo(),
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
