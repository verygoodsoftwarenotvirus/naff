package v1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"cmd/server/v1/coverage_test.go": coverageTestDotGo(proj),
		"cmd/server/v1/doc.go":           docDotGo(),
		"cmd/server/v1/main.go":          mainDotGo(proj),
		"cmd/server/v1/wire.go":          wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
