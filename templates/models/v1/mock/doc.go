package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	code := jen.NewFile("mock")

	code.PackageComment(`Package mock provides mockable implementations of every interface
defined in the outer models package`)

	return code
}
