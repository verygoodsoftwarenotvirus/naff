package fake

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)
	utils.AddImports(proj, ret)

	ret.PackageComment("Package fake provides fake model builders\n")

	return ret
}
