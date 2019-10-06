package mockmodels

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mockWebhookDataManagerDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"WebhookDataManager",
	).Op("=").Parens(jen.Op("*").ID("WebhookDataManager")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("WebhookDataManager").Struct(
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
	return ret
}
