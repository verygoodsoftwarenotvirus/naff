package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("tracing")

	ret.PackageComment("Package tracing provides common functions for attaching values to trace spans\n")

	return ret
}
