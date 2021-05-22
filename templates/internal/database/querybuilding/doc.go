package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo(vendor, dbDesc string) *jen.File {
	code := jen.NewFile(vendor)

	code.PackageCommentf("Package %s provides a Database implementation that is compatible with %s\n", vendor, dbDesc)

	return code
}
