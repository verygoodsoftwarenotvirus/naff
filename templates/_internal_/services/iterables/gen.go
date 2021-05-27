package iterables

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/services"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, typ := range proj.DataTypes {
		pn := typ.Name.PackageName()

		files := map[string]*jen.File{
			fmt.Sprintf("%s/wire.go", pn):              wireDotGo(proj, typ),
			fmt.Sprintf("%s/doc.go", pn):               docDotGo(typ),
			fmt.Sprintf("%s/http_routes.go", pn):       httpRoutesDotGo(proj, typ),
			fmt.Sprintf("%s/http_routes_test.go", pn):  httpRoutesTestDotGo(proj, typ),
			fmt.Sprintf("%s/http_helpers_test.go", pn): httpHelpersTestDotGo(proj, typ),
			fmt.Sprintf("%s/service.go", pn):           iterableServiceDotGo(proj, typ),
			fmt.Sprintf("%s/service_test.go", pn):      iterableServiceTestDotGo(proj, typ),
		}

		for path, file := range files {
			if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
				return err
			}
		}
	}

	return nil
}
