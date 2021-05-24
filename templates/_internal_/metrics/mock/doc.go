package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf(`Package %s provides mockable implementations of every interface
defined in the outer metrics package.`, packageName)

	return code
}
