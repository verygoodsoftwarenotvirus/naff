package items

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	for _, typ := range types {
		pn := typ.Name.PluralRouteName()
		for path, file := range map[string]*jen.File{
			fmt.Sprintf("services/v1/%s/middleware.go", pn):          middlewareDotGo(typ),
			fmt.Sprintf("services/v1/%s/middleware_test.go", pn):     middlewareTestDotGo(typ),
			fmt.Sprintf("services/v1/%s/wire.go", pn):                wireDotGo(typ),
			fmt.Sprintf("services/v1/%s/doc.go", pn):                 docDotGo(typ),
			fmt.Sprintf("services/v1/%s/http_routes.go", pn):         httpRoutesDotGo(typ),
			fmt.Sprintf("services/v1/%s/http_routes_test.go", pn):    httpRoutesTestDotGo(typ),
			fmt.Sprintf("services/v1/%s/%s_service_test.go", pn, pn): iterableServiceTestDotGo(typ),
			fmt.Sprintf("services/v1/%s/%s_service.go", pn, pn):      iterableServiceDotGo(typ),
		} {
			if err := utils.RenderFile(path, file); err != nil {
				return err
			}
		}
	}

	return nil
}
