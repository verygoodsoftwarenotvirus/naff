package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.PackageComment(`Package client provides an HTTP client that can communicate with and interpret the responses
of an instance of the todo service.`)

	return code
}
