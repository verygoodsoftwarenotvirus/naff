package users

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/services/users"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"service.go":           usersServiceDotGo(proj),
		"service_test.go":      usersServiceTestDotGo(proj),
		"wire.go":              wireDotGo(proj),
		"doc.go":               docDotGo(proj),
		"http_routes.go":       httpRoutesDotGo(proj),
		"http_routes_test.go":  httpRoutesTestDotGo(proj),
		"http_helpers_test.go": httpHelpersTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed service.gotpl
var usersServiceTemplate string

func usersServiceDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersServiceTemplate, nil)
}

//go:embed service_test.gotpl
var usersServiceTestTemplate string

func usersServiceTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersServiceTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
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
