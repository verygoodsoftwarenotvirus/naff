package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("webhooks")
	utils.AddImports(ret)

	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
