package mocksearch

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf("Package %s provides an interface-compatible search index mock\n", packageName)

	return code
}
