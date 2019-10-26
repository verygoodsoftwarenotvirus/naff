package http

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("httpserver")

	ret.PackageComment("Package httpserver provides an HTTP server comprised of multiple HTTP services\n")

	return ret
}
