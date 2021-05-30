package server

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "server"

	basePackagePath = "internal/server"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"wire.go":             wireDotGo(proj),
		"config.go":           configDotGo(proj),
		"config_test.go":      configTestDotGo(proj),
		"doc.go":              docDotGo(proj),
		"http_routes.go":      routesDotGo(proj),
		"http_server.go":      httpServerDotGo(proj),
		"http_server_test.go": httpServerTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
