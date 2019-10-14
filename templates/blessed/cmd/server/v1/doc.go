package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("main")

	ret.PackageComment(`Command server is the main compilable application that runs an instance of the todo service
`)

	return ret
}
