package iterables

import (
	_ "embed"
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
			fmt.Sprintf("%s/http_helpers_test.go", pn): httpHelpersTestDotGo(proj, typ),
			fmt.Sprintf("%s/service.go", pn):           serviceDotGo(proj, typ),
			fmt.Sprintf("%s/service_test.go", pn):      serviceTestDotGo(proj, typ),
			fmt.Sprintf("%s/http_routes.go", pn):       httpRoutesDotGo(proj, typ),
			fmt.Sprintf("%s/http_routes_test.go", pn):  httpRoutesTestDotGo(proj, typ),
		}

		for path, file := range files {
			if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
				return err
			}
		}

		strFiles := map[string]string{
			fmt.Sprintf("%s/config.go", pn): configDotGo(proj),
			fmt.Sprintf("%s/wire.go", pn):   wireDotGo(proj),
			fmt.Sprintf("%s/doc.go", pn):    docDotGo(proj),
		}

		for path, file := range strFiles {
			if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
				return err
			}
		}
	}

	return nil
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}
