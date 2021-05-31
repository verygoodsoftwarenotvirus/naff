package bleve

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/search/bleve"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"bleve.go":      bleveDotGo(proj),
		"bleve_test.go": bleveTestDotGo(proj),
		"doc.go":        docDotGo(proj),
		"mappings.go":   mappingsDotGo(proj),
		"utils.go":      utilsDotGo(proj),
		"utils_test.go": utilsTestDotGo(proj),
		"wire.go":       wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed bleve.gotpl
var bleveTemplate string

func bleveDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bleveTemplate, nil)
}

//go:embed bleve_test.gotpl
var bleveTestTemplate string

func bleveTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bleveTestTemplate, nil)
}

//go:embed doc.gotpl
var Template string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, Template, nil)
}

//go:embed mappings.gotpl
var mappingsTemplate string

func mappingsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mappingsTemplate, nil)
}

//go:embed utils.gotpl
var utilsTemplate string

func utilsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, utilsTemplate, nil)
}

//go:embed utils_test.gotpl
var utilsTestTemplate string

func utilsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, utilsTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}
