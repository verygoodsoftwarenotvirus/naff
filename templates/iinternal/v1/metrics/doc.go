package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("metrics")

	code.PackageComment("Package metrics provides metrics collection functions and structs\n")

	return code
}
