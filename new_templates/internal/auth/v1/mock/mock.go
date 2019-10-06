package authmock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mockDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("_").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1", "Authenticator").Op("=").Parens(jen.Op("*").ID("Authenticator")).Call(jen.ID("nil")))
	ret.Add(jen.Null().Type().ID("Authenticator").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
