package iterables

import (
	"fmt"
	"io/ioutil"
	"strings"

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

			if path == fmt.Sprintf("services/v1/%s/http_routes_test.go", pn) {
				p := utils.BuildTemplatePath(proj.OutputPath, path)

				fileBytes, err := ioutil.ReadFile(p)
				if err != nil {
					return fmt.Errorf("error reading recently written file: %w", err)
				}

				newFile := strings.Replace(string(fileBytes), `
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	mocksearch "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"`, `
	mocksearch "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"`, 1)

				if err := ioutil.WriteFile(p, []byte(newFile), 0644); err != nil {
					return fmt.Errorf("error correcting file: %w", err)
				}
			}
		}
	}

	return nil
}
