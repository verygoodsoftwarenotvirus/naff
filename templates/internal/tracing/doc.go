package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageCommentf("Package %s provides common functions for attaching values to trace spans\n", packageName)

	return code
}
