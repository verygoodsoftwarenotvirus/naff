package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockWebhookDataManagerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),
	)
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
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
