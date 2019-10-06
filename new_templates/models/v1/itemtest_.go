package models

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemTestDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)


	ret.Add(jen.Func().ID("TestItem_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
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
