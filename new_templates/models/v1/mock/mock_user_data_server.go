package mockmodels

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mockUserDataServerDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"UserDataServer",
	).Op("=").Parens(jen.Op("*").ID("UserDataServer")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("UserDataServer").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
