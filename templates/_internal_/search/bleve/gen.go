package bleve

import (
	_ "embed"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "bleve"

	basePackagePath = "internal/search/bleve"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"doc.go":        docDotGo(proj),
		"utils.go":      utilsDotGo(proj),
		"utils_test.go": utilsTestDotGo(proj),
		"wire.go":       wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{
		"bleve.go":      bleveDotGo(proj),
		"bleve_test.go": bleveTestDotGo(proj),
		"mappings.go":   mappingsDotGo(proj),
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed doc.gotpl
var Template string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, Template, nil)
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
