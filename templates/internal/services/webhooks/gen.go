package webhooks

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "webhooks"

	basePackaegPath = "internal/services/webhooks"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":                   docDotGo(),
		"wire.go":                  wireDotGo(proj),
		"http_routes.go":           httpRoutesDotGo(proj),
		"http_routes_test.go":      httpRoutesTestDotGo(proj),
		"middleware.go":            middlewareDotGo(proj),
		"middleware_test.go":       middlewareTestDotGo(proj),
		"webhooks_service.go":      webhooksServiceDotGo(proj),
		"webhooks_service_test.go": webhooksServiceTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackaegPath, path), file); err != nil {
			return err
		}
	}

	return nil
}
