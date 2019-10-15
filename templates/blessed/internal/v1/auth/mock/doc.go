package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("mock")

	ret.PackageComment(`Package mock provides mockable implementations of every interface
defined in the outer auth package.`)

	return ret
}
