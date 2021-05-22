package server

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageComment(`Command server is the main compilable application that runs an instance of the todo service
`)

	return code
}
