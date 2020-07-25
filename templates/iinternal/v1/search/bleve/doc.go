package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf("Package %s provides an interface-compatible wrapper around the bleve indexer\n", packageName)

	return code
}
