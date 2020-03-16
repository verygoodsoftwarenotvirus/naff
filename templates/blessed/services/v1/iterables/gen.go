package iterables

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	for _, typ := range pkg.DataTypes {
		pn := typ.Name.PackageName()
		for path, file := range map[string]*jen.File{
			fmt.Sprintf("services/v1/%s/middleware.go", pn):          middlewareDotGo(pkg, typ),
			fmt.Sprintf("services/v1/%s/middleware_test.go", pn):     middlewareTestDotGo(pkg, typ),
			fmt.Sprintf("services/v1/%s/wire.go", pn):                wireDotGo(pkg, typ),
			fmt.Sprintf("services/v1/%s/doc.go", pn):                 docDotGo(typ),
			fmt.Sprintf("services/v1/%s/http_routes.go", pn):         httpRoutesDotGo(pkg, typ),
			fmt.Sprintf("services/v1/%s/http_routes_test.go", pn):    httpRoutesTestDotGo(pkg, typ),
			fmt.Sprintf("services/v1/%s/%s_service.go", pn, pn):      iterableServiceDotGo(pkg, typ),
			fmt.Sprintf("services/v1/%s/%s_service_test.go", pn, pn): iterableServiceTestDotGo(pkg, typ),
		} {
			if err := utils.RenderGoFile(pkg, path, file); err != nil {
				return err
			}
		}
	}

	return nil
}
