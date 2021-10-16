package workers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func docDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)
	code.PackageCommentf("Package %s provides async data writing workers.\n", packageName)

	utils.AddImports(proj, code, false)

	return code
}
