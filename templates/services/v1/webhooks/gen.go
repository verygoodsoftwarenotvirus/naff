package webhooks

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "webhooks"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/webhooks/doc.go":                   docDotGo(),
		"services/v1/webhooks/wire.go":                  wireDotGo(proj),
		"services/v1/webhooks/http_routes.go":           httpRoutesDotGo(proj),
		"services/v1/webhooks/http_routes_test.go":      httpRoutesTestDotGo(proj),
		"services/v1/webhooks/middleware.go":            middlewareDotGo(proj),
		"services/v1/webhooks/middleware_test.go":       middlewareTestDotGo(proj),
		"services/v1/webhooks/webhooks_service.go":      webhooksServiceDotGo(proj),
		"services/v1/webhooks/webhooks_service_test.go": webhooksServiceTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
