package database

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/database"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"database.go":      databaseDotGo(proj),
		"database_mock.go": databaseMockDotGo(proj),
		"doc.go":           docDotGo(proj),
		"wire.go":          wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed database.gotpl
var databaseTemplate string

func databaseDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, databaseTemplate, nil)
}

//go:embed database_mock.gotpl
var databaseMockTemplate string

func databaseMockDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, databaseMockTemplate, nil)
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
