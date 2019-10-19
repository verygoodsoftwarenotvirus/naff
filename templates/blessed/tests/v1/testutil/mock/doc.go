package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("mock")

	ret.PackageComment("Package mock contains mock structures that are useful for unit/integration/load tests\n")

	return ret
}
