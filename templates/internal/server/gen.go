package server

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "server"

	basePackagePath = "internal/server"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":                      docDotGo(),
		"middleware_test.go":          middlewareTestDotGo(proj),
		"routes.go":                   routesDotGo(proj),
		"server.go":                   serverDotGo(proj),
		"server_test.go":              serverTestDotGo(proj),
		"middleware.go":               middlewareDotGo(proj),
		"wire.go":                     wireDotGo(proj),
		"wire_param_fetchers.go":      wireParamFetchersDotGo(proj),
		"wire_param_fetchers_test.go": wireParamFetchersTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
