package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("database")

	ret.PackageComment("Package database provides interface abstractions for interacting with relational data stores\n")

	return ret
}
