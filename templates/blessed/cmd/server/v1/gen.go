package v1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"cmd/server/v1/coverage_test.go": coverageTestDotGo(pkgRoot, types),
		"cmd/server/v1/doc.go":           docDotGo(),
		"cmd/server/v1/main.go":          mainDotGo(pkgRoot, types),
		"cmd/server/v1/wire.go":          wireDotGo(pkgRoot, types),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
