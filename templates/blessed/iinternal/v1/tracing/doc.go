package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile(packageName)

	ret.PackageCommentf("Package %s provides common functions for attaching values to trace spans\n", packageName)

	return ret
}
