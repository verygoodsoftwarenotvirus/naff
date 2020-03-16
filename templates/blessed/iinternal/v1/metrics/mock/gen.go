package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkg *models.Project) error {
	files := map[string]*jen.File{
		"internal/v1/metrics/mock/counter.go": counterDotGo(pkg),
		"internal/v1/metrics/mock/doc.go":     docDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkg, path, file); err != nil {
			return err
		}
	}

	return nil
}
