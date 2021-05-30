package webhooks

import (
	_ "embed"
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
		"config.go":            configDotGo(proj),
		"config_test.go":       configTestDotGo(proj),
		"doc.go":               docDotGo(),
		"wire.go":              wireDotGo(proj),
		"http_helpers_test.go": httpHelpersTestDotGo(proj),
		"http_routes.go":       httpRoutesDotGo(proj),
		"http_routes_test.go":  httpRoutesTestDotGo(proj),
		"service.go":           webhooksServiceDotGo(proj),
		"service_test.go":      webhooksServiceTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackaegPath, path), file); err != nil {
			return err
		}
	}

	return nil
}
