package bleve

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "bleve"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/search/bleve/bleve.go":      bleveDotGo(proj),
		"internal/v1/search/bleve/bleve_test.go": bleveTestDotGo(proj),
		"internal/v1/search/bleve/doc.go":        docDotGo(),
		"internal/v1/search/bleve/mappings.go":   mappingsDotGo(proj),
		"internal/v1/search/bleve/utils.go":      utilsDotGo(proj),
		"internal/v1/search/bleve/utils_test.go": utilsTestDotGo(proj),
		"internal/v1/search/bleve/wire.go":       wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
