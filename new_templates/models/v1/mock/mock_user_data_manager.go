package mockmodels

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mockUserDataManagerDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"UserDataManager",
	).Op("=").Parens(jen.Op("*").ID("UserDataManager")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("UserDataManager").Struct(
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
	return ret
}
