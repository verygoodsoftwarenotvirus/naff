package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("postgres")

	ret.PackageComment("Package postgres provides a Database implementation that is compatible with PostgreSQL instances\n")

	return ret
}
