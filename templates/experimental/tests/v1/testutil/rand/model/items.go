package model

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsDotGo() *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// RandomItemCreationInput creates a random ItemInput").ID("RandomItemCreationInput").Params().Params(jen.Op("*").ID("models").Dot(
		"ItemCreationInput",
	)).Block(
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"ItemCreationInput",
		).Valuesln(jen.ID("Name").Op(":").ID("fake").Dot(
			"Word",
		).Call(), jen.ID("Details").Op(":").ID("fake").Dot(
			"Sentence",
		).Call()),
		jen.Return().ID("x"),
	),

		jen.Line(),
	)
	return ret
}
