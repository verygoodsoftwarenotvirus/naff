package testutil

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("testutil")

	ret.PackageComment("Package testutil contains common functions for integration/load tests\n")

	return ret
}
