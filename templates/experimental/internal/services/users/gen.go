package users

import (
	_ "embed"
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
		"http_routes_test.go":  httpRoutesTestDotGo(proj),
		"service.go":           serviceDotGo(proj),
		"service_test.go":      serviceTestDotGo(proj),
		"wire.go":              wireDotGo(proj),
		"doc.go":               docDotGo(proj),
		"http_helpers_test.go": httpHelpersTestDotGo(proj),
		"http_routes.go":       httpRoutesDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
