package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"internal/v1/auth/mock/mock.go": mockDotGo(pkgRoot, types),
		"internal/v1/auth/mock/doc.go":  docDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
