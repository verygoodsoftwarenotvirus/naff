package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("encoding")

	ret.PackageComment("Package encoding provides HTTP response encoding abstractions\n")

	return ret
}
