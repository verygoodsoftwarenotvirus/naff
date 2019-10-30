package webhooks

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/webhooks/wire.go":                  wireDotGo(pkgRoot),
		"services/v1/webhooks/doc.go":                   docDotGo(),
		"services/v1/webhooks/http_routes.go":           httpRoutesDotGo(pkgRoot),
		"services/v1/webhooks/http_routes_test.go":      httpRoutesTestDotGo(pkgRoot),
		"services/v1/webhooks/middleware.go":            middlewareDotGo(pkgRoot),
		"services/v1/webhooks/middleware_test.go":       middlewareTestDotGo(pkgRoot),
		"services/v1/webhooks/webhooks_service.go":      webhooksServiceDotGo(pkgRoot),
		"services/v1/webhooks/webhooks_service_test.go": webhooksServiceTestDotGo(pkgRoot),
	}

	for path, file := range files {
		if err := utils.RenderFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
