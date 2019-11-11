package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/auth/mock/mock.go": mockDotGo(pkg),
		"internal/v1/auth/mock/doc.go":  docDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg.OutputPath, path, file); err != nil {
			return err
		}
	}

	return nil
}
