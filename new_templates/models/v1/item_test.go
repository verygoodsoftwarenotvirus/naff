package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func itemTestDotGo() *jen.File {
	ret := jen.NewFile("models")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("TestItem_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("i").Op(":=").Op("&").ID("Item").Valuesln(),
			jen.ID("expected").Op(":=").Op("&").ID("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").Lit("expected name"), jen.ID("Details").Op(":").Lit("expected details")),
			jen.ID("i").Dot(
				"Update",
			).Call(jen.ID("expected")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"Name",
			), jen.ID("i").Dot(
				"Name",
			)),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"Details",
			), jen.ID("i").Dot(
				"Details",
			)),
		)),
	),
	)
	return ret
}
