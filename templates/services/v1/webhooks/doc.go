package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("webhooks")

	code.PackageComment("Package webhooks provides a series of HTTP handlers for managing webhooks in a compatible database.\n")

	return code
}
