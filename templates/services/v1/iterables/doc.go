package iterables

import (
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func docDotGo(typ models.DataType) *jen.File {
	pn := typ.Name.PackageName()
	cn := strings.ToLower(typ.Name.Plural())
	code := jen.NewFile(pn)

	code.PackageCommentf("Package %s provides a series of HTTP handlers for managing %s in a compatible database.\n", pn, cn)

	return code
}
