package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func docDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	code.PackageCommentf("Package %s provides a series of HTTP handlers for managing %s in a compatible database.\n", typ.Name.PackageName(), typ.Name.PluralCommonName())

	utils.AddImports(proj, code, false)

	return code
}
