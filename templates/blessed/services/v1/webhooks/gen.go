package webhooks

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"services/v1/webhooks/doc.go":                   docDotGo(),
		"services/v1/webhooks/wire.go":                  wireDotGo(pkg),
		"services/v1/webhooks/http_routes.go":           httpRoutesDotGo(pkg),
		"services/v1/webhooks/http_routes_test.go":      httpRoutesTestDotGo(pkg),
		"services/v1/webhooks/middleware.go":            middlewareDotGo(pkg),
		"services/v1/webhooks/middleware_test.go":       middlewareTestDotGo(pkg),
		"services/v1/webhooks/webhooks_service.go":      webhooksServiceDotGo(pkg),
		"services/v1/webhooks/webhooks_service_test.go": webhooksServiceTestDotGo(pkg),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
