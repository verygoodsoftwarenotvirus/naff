package server

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/server"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"wire.go":             wireDotGo(proj),
		"config.go":           configDotGo(proj),
		"config_test.go":      configTestDotGo(proj),
		"doc.go":              docDotGo(proj),
		"http_routes.go":      routesDotGo(proj),
		"http_server.go":      httpServerDotGo(proj),
		"http_server_test.go": httpServerTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
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

//go:embed http_routes.gotpl
var routesTemplate string

func routesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, routesTemplate, nil)
}

//go:embed http_server.gotpl
var httpServerTemplate string

func httpServerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpServerTemplate, nil)
}

//go:embed http_server_test.gotpl
var httpServerTestTemplate string

func httpServerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpServerTestTemplate, nil)
}
