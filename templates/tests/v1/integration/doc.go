package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("integration")

	code.PackageComment(`Package integration is a series of tests which utilize our HTTP client to talk to a running
HTTP server to validate behaviors, inputs, and outputs.`)

	return code
}
