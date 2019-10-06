package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func httpRoutesDotGo() *jen.File {
	ret := jen.NewFile("items")
	utils.AddImports(ret)
	return ret
}
