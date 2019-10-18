package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("client")

	ret.PackageComment(`Package client provides an HTTP client that can communicate with and interpret the responses
of an instance of the todo service.`)

	return ret
}
