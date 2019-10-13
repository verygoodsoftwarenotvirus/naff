package http

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"server/v1/http/wire_param_fetchers_test.go": wireParamFetchersTestDotGo(),
		"server/v1/http/middleware.go":               middlewareDotGo(),
		"server/v1/http/middleware_test.go":          middlewareTestDotGo(),
		"server/v1/http/routes.go":                   routesDotGo(),
		"server/v1/http/wire.go":                     wireDotGo(),
		"server/v1/http/doc.go":                      docDotGo(),
		"server/v1/http/server.go":                   serverDotGo(),
		"server/v1/http/server_test.go":              serverTestDotGo(),
		"server/v1/http/wire_param_fetchers.go":      wireParamFetchersDotGo(),
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
