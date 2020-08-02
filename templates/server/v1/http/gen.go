package httpserver

import (
	"log"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "httpserver"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"server/v1/http/doc.go":                      docDotGo(),
		"server/v1/http/middleware_test.go":          middlewareTestDotGo(proj),
		"server/v1/http/routes.go":                   routesDotGo(proj),
		"server/v1/http/server.go":                   serverDotGo(proj),
		"server/v1/http/server_test.go":              serverTestDotGo(proj),
		"server/v1/http/middleware.go":               middlewareDotGo(proj),
		"server/v1/http/wire.go":                     wireDotGo(proj),
		"server/v1/http/wire_param_fetchers.go":      wireParamFetchersDotGo(proj),
		"server/v1/http/wire_param_fetchers_test.go": wireParamFetchersTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			log.Printf("error rendering %q: %v\n", path, err)
			return err
		}
	}

	return nil
}
