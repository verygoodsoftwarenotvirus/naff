package webhooks

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) error {
	files := map[string]*jen.File{
		"services/v1/webhooks/wire.go":                  wireDotGo(),
		"services/v1/webhooks/doc.go":                   docDotGo(),
		"services/v1/webhooks/http_routes.go":           httpRoutesDotGo(),
		"services/v1/webhooks/http_routes_test.go":      httpRoutesTestDotGo(),
		"services/v1/webhooks/middleware.go":            middlewareDotGo(),
		"services/v1/webhooks/middleware_test.go":       middlewareTestDotGo(),
		"services/v1/webhooks/webhooks_service.go":      webhooksServiceDotGo(),
		"services/v1/webhooks/webhooks_service_test.go": webhooksServiceTestDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderFile(path, file); err != nil {
			return err
		}
	}

	return nil
}
