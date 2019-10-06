package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockItemDataManagerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"ItemDataManager",
	).Op("=").Parens(jen.Op("*").ID("ItemDataManager")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("ItemDataManager").Struct(
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
	return ret
}
