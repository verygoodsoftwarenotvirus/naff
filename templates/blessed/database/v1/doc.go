package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("database")

	code.PackageComment("Package database provides interface abstractions for interacting with relational data stores\n")

	return code
}
