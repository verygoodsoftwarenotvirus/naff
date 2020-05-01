package v1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"server/v1/doc.go":    docDotGo(),
		"server/v1/server.go": serverDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
