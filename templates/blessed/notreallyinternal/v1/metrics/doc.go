package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	ret.PackageComment("Package metrics provides metrics collection functions and structs\n")

	return ret
}
