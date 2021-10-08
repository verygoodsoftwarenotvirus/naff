package websockets

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/services/websockets"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"config.go": configDotGo(proj),
		"config_test.go": configTestDotGo(proj),
		"data_changes.go": dataChangesDotGo(proj),
		"data_changes_test.go": dataChangesTestDotGo(proj),
		"doc.go": docDotGo(proj),
		"http_routes.go": httpRoutesDotGo(proj),
		"http_routes_test.go": httpRoutesTestDotGo(proj),
		"http_helpers_test.go": httpHelpersTestDotGo(proj),
		"service.go": serviceDotGo(proj),
		"service_test.go": serviceTestDotGo(proj),
		"wire.go": wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
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

//go:embed data_changes.gotpl
var dataChangesTemplate string
func dataChangesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, dataChangesTemplate, nil)
}

//go:embed data_changes_test.gotpl
var dataChangesTestTemplate string
func dataChangesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, dataChangesTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string
func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
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

//go:embed http_helpers_test.gotpl
var httpHelpersTestTemplate string
func httpHelpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpHelpersTestTemplate, nil)
}

//go:embed service.gotpl
var serviceTemplate string
func serviceDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serviceTemplate, nil)
}

//go:embed service_test.gotpl
var serviceTestTemplate string
func serviceTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serviceTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string
func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

