package httpserver

import (
	"log"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {

	files := map[string]*jen.File{
		"server/v1/http/doc.go":                      docDotGo(),
		"server/v1/http/middleware_test.go":          middlewareTestDotGo(pkgRoot, types),
		"server/v1/http/routes.go":                   routesDotGo(pkgRoot, types),
		"server/v1/http/server.go":                   serverDotGo(pkgRoot, types),
		"server/v1/http/server_test.go":              serverTestDotGo(pkgRoot, types),
		"server/v1/http/middleware.go":               middlewareDotGo(pkgRoot, types),
		"server/v1/http/wire.go":                     wireDotGo(pkgRoot, types),
		"server/v1/http/wire_param_fetchers.go":      wireParamFetchersDotGo(pkgRoot, types),
		"server/v1/http/wire_param_fetchers_test.go": wireParamFetchersTestDotGo(pkgRoot, types),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			log.Printf("error rendering %q: %v\n", path, err)
			return err
		}
	}

	return nil
}
