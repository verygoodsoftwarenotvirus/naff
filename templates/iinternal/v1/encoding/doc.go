package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("encoding")

	code.PackageComment("Package encoding provides HTTP response encoding abstractions\n")

	return code
}
