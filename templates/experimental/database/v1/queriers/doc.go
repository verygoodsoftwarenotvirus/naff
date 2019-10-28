package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo(vendor, dbDesc string) *jen.File {
	ret := jen.NewFile(vendor)

	ret.PackageCommentf("Package %s provides a Database implementation that is compatible with %s\n", vendor, dbDesc)

	return ret
}
