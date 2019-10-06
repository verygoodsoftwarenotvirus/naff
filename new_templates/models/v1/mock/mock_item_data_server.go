package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockItemDataServerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"ItemDataServer",
	).Op("=").Parens(jen.Op("*").ID("ItemDataServer")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("ItemDataServer").Struct(
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
	return ret
}
