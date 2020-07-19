package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("models")

	code.PackageComment("Package models provides the core data models for the service\n")

	return code
}
