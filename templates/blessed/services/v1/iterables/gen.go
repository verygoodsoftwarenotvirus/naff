package iterables

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, typ := range proj.DataTypes {
		pn := typ.Name.PackageName()
		for path, file := range map[string]*jen.File{
			fmt.Sprintf("services/v1/%s/middleware.go", pn):          middlewareDotGo(proj, typ),
			fmt.Sprintf("services/v1/%s/middleware_test.go", pn):     middlewareTestDotGo(proj, typ),
			fmt.Sprintf("services/v1/%s/wire.go", pn):                wireDotGo(proj, typ),
			fmt.Sprintf("services/v1/%s/wire_test.go", pn):           wireTestDotGo(proj, typ),
			fmt.Sprintf("services/v1/%s/doc.go", pn):                 docDotGo(typ),
			fmt.Sprintf("services/v1/%s/http_routes.go", pn):         httpRoutesDotGo(proj, typ),
			fmt.Sprintf("services/v1/%s/http_routes_test.go", pn):    httpRoutesTestDotGo(proj, typ),
			fmt.Sprintf("services/v1/%s/%s_service.go", pn, pn):      iterableServiceDotGo(proj, typ),
			fmt.Sprintf("services/v1/%s/%s_service_test.go", pn, pn): iterableServiceTestDotGo(proj, typ),
		} {
			if err := utils.RenderGoFile(proj, path, file); err != nil {
				return err
			}
		}
	}

	return nil
}
