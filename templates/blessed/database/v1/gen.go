package v1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"database/v1/database.go":      databaseDotGo(pkgRoot, types),
		"database/v1/database_mock.go": databaseMockDotGo(pkgRoot, types),
		"database/v1/doc.go":           docDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
