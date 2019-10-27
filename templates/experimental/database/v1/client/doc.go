package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	ret.PackageComment(`Package dbclient provides an abstraction around database queriers. The primary
purpose of this abstraction is to contain all the necessary logging and tracing
steps in a single place, so that the actual package that is responsible for
executing queries and loading their return values into structs isn't burdened
with inconsistent logging`)

	return ret
}
