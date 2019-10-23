package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("frontend")

	ret.PackageComment("Package frontend implements a frontend, mostly-static file server\n")

	return ret
}
