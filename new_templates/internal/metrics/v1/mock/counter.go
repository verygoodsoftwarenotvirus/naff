package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func counterDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("metrics").Dot(
		"UnitCounter",
	).Op("=").Parens(jen.Op("*").ID("UnitCounter")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("UnitCounter").Struct(
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
	return ret
}
