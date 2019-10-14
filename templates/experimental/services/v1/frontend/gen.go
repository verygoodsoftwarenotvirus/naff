package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/frontend/wire.go":                  wireDotGo(),
		"services/v1/frontend/doc.go":                   docDotGo(),
		"services/v1/frontend/frontend_service.go":      frontendServiceDotGo(),
		"services/v1/frontend/frontend_service_test.go": frontendServiceTestDotGo(),
		"services/v1/frontend/http_routes.go":           httpRoutesDotGo(),
		"services/v1/frontend/http_routes_test.go":      httpRoutesTestDotGo(),
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
