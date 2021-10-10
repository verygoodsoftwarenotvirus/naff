package webhooks

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackaegPath = "internal/services/webhooks"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"config.go":            configDotGo(proj),
		"config_test.go":       configTestDotGo(proj),
		"doc.go":               docDotGo(proj),
		"wire.go":              wireDotGo(proj),
		"http_helpers_test.go": httpHelpersTestDotGo(proj),
		"http_routes.go":       httpRoutesDotGo(proj),
		"http_routes_test.go":  httpRoutesTestDotGo(proj),
		"service.go":           webhooksServiceDotGo(proj),
		"service_test.go":      webhooksServiceTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackaegPath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed http_helpers_test.gotpl
var httpHelpersTestTemplate string

func httpHelpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpHelpersTestTemplate, nil)
}

//go:embed http_routes.gotpl
var httpRoutesTemplate string

func httpRoutesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpRoutesTemplate, nil)
}

//go:embed http_routes_test.gotpl
var httpRoutesTestTemplate string

func httpRoutesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpRoutesTestTemplate, nil)
}

//go:embed service.gotpl
var webhooksServiceTemplate string

func webhooksServiceDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksServiceTemplate, nil)
}

//go:embed service_test.gotpl
var webhooksServiceTestTemplate string

func webhooksServiceTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksServiceTestTemplate, nil)
}
