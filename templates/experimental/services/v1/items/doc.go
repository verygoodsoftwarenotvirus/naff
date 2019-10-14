package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func docDotGo() *jen.File {
	ret := jen.NewFile("items")

	utils.AddImports(ret)

	return ret
}
