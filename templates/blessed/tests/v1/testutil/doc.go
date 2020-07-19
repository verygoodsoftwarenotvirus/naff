package testutil

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("testutil")

	code.PackageComment("Package testutil contains common functions for integration/load tests\n")

	return code
}
