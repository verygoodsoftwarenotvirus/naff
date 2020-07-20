package fake

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func fakeDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)
	utils.AddImports(proj, code)

	code.Add(utils.FakeSeedFunc())

	return code
}
