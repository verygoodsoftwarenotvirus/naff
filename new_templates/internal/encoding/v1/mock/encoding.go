package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func encodingDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("encoding").Dot(
		"EncoderDecoder",
	).Op("=").Parens(jen.Op("*").ID("EncoderDecoder")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("EncoderDecoder").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
