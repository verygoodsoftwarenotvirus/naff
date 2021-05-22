package bleve

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "bleve"

	basePackagePath = "internal/search/bleve"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"bleve.go":      bleveDotGo(proj),
		"bleve_test.go": bleveTestDotGo(proj),
		"doc.go":        docDotGo(),
		"mappings.go":   mappingsDotGo(proj),
		"utils.go":      utilsDotGo(proj),
		"utils_test.go": utilsTestDotGo(proj),
		"wire.go":       wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
