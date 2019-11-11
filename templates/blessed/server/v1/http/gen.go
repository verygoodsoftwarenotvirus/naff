package httpserver

import (
	"log"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {

	files := map[string]*jen.File{
		"server/v1/http/doc.go":                      docDotGo(),
		"server/v1/http/middleware_test.go":          middlewareTestDotGo(pkg),
		"server/v1/http/routes.go":                   routesDotGo(pkg),
		"server/v1/http/server.go":                   serverDotGo(pkg),
		"server/v1/http/server_test.go":              serverTestDotGo(pkg),
		"server/v1/http/middleware.go":               middlewareDotGo(pkg),
		"server/v1/http/wire.go":                     wireDotGo(pkg),
		"server/v1/http/wire_param_fetchers.go":      wireParamFetchersDotGo(pkg),
		"server/v1/http/wire_param_fetchers_test.go": wireParamFetchersTestDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			log.Printf("error rendering %q: %v\n", path, err)
			return err
		}
	}

	return nil
}
