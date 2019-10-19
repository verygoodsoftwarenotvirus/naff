package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsTestDotGo() *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("checkItemEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("models").Dot(
		"Item",
	)).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ID",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"Name",
	),
	jen.ID("actual").Dot(
			"Name",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"Details",
	),
	jen.ID("actual").Dot(
			"Details",
		)),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"CreatedOn",
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyItem").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("models").Dot(
		"Item",
	)).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"ItemCreationInput",
		).Valuesln(jen.ID("Name").Op(":").ID("fake").Dot(
			"Word",
		).Call(), jen.ID("Details").Op(":").ID("fake").Dot(
			"Sentence",
		).Call()),
		jen.List(jen.ID("y"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
			"CreateItem",
		).Call(jen.Qual("context", "Background").Call(), jen.ID("x")),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.Return().ID("y"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestItems").Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
		jen.ID("test").Dot(
			"Parallel",
		).Call(),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be createable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot(
					"Name",
				).Call()),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Item",
				).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details")),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateItem",
				).Call(jen.ID("ctx"), jen.Op("&").ID("models").Dot(
					"ItemCreationInput",
				).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
					"Name",
	),
	jen.ID("Details").Op(":").ID("expected").Dot(
					"Details",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.ID("checkItemEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("premade")),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveItem",
				).Call(jen.ID("ctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetItem",
				).Call(jen.ID("ctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkItemEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot(
					"NotZero",
				).Call(jen.ID("t"), jen.ID("actual").Dot(
					"ArchivedOn",
				)),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot(
					"Name",
				).Call()),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),

		jen.Var().ID("expected").Index().Op("*").ID("models").Dot(
					"Item",
				),
				jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
					jen.ID("expected").Op("=").ID("append").Call(jen.ID("expected"), jen.ID("buildDummyItem").Call(jen.ID("t"))),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetItems",
				).Call(jen.ID("ctx"), jen.ID("nil")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("assert").Dot(
					"True",
				).Call(jen.ID("t"), jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(
					"Items",
				)), jen.Lit("expected %d to be <= %d"), jen.ID("len").Call(jen.ID("expected")), jen.ID("len").Call(jen.ID("actual").Dot(
					"Items",
				))),
				jen.For(jen.List(jen.ID("_"), jen.ID("item")).Op(":=").Range().ID("actual").Dot(
					"Items",
				)).Block(
					jen.ID("err").Op("=").ID("todoClient").Dot(
						"ArchiveItem",
					).Call(jen.ID("ctx"), jen.ID("item").Dot(
						"ID",
					)),
					jen.ID("assert").Dot(
						"NoError",
					).Call(jen.ID("t"), jen.ID("err")),
				),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should return an error when trying to read something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot(
					"Name",
				).Call()),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetItem",
				).Call(jen.ID("ctx"), jen.ID("nonexistentID")),
				jen.ID("assert").Dot(
					"Error",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot(
					"Name",
				).Call()),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Item",
				).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details")),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateItem",
				).Call(jen.ID("ctx"), jen.Op("&").ID("models").Dot(
					"ItemCreationInput",
				).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
					"Name",
	),
	jen.ID("Details").Op(":").ID("expected").Dot(
					"Details",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetItem",
				).Call(jen.ID("ctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkItemEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveItem",
				).Call(jen.ID("ctx"), jen.ID("actual").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should return an error when trying to update something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot(
					"Name",
				).Call()),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.ID("err").Op(":=").ID("todoClient").Dot(
					"UpdateItem",
				).Call(jen.ID("ctx"), jen.Op("&").ID("models").Dot(
					"Item",
				).Valuesln(jen.ID("ID").Op(":").ID("nonexistentID"))),
				jen.ID("assert").Dot(
					"Error",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should be updatable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot(
					"Name",
				).Call()),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Item",
				).Valuesln(jen.ID("Name").Op(":").Lit("new name"), jen.ID("Details").Op(":").Lit("new details")),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateItem",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"ItemCreationInput",
				).Valuesln(jen.ID("Name").Op(":").Lit("old name"), jen.ID("Details").Op(":").Lit("old details"))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.List(jen.ID("premade").Dot(
					"Name",
	),
	jen.ID("premade").Dot(
					"Details",
				)).Op("=").List(jen.ID("expected").Dot(
					"Name",
	),
	jen.ID("expected").Dot(
					"Details",
				)),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"UpdateItem",
				).Call(jen.ID("ctx"), jen.ID("premade")),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetItem",
				).Call(jen.ID("ctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkItemEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot(
					"NotNil",
				).Call(jen.ID("t"), jen.ID("actual").Dot(
					"UpdatedOn",
				)),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveItem",
				).Call(jen.ID("ctx"), jen.ID("actual").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot(
					"Name",
				).Call()),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
					"Item",
				).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details")),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateItem",
				).Call(jen.ID("ctx"), jen.Op("&").ID("models").Dot(
					"ItemCreationInput",
				).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
					"Name",
	),
	jen.ID("Details").Op(":").ID("expected").Dot(
					"Details",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveItem",
				).Call(jen.ID("ctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
	),
	jen.Line(),
	)
	return ret
}
