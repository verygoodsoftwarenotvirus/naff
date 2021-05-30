package fakes

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeItem builds a faked item."),
		jen.Line(),
		jen.Func().ID("BuildFakeItem").Params().Params(jen.Op("*").ID("types").Dot("Item")).Body(
			jen.Return().Op("&").ID("types").Dot("Item").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("ExternalID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("Name").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call(), jen.ID("Details").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call(), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("BelongsToAccount").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeItemList builds a faked ItemList."),
		jen.Line(),
		jen.Func().ID("BuildFakeItemList").Params().Params(jen.Op("*").ID("types").Dot("ItemList")).Body(
			jen.Var().Defs(
				jen.ID("examples").Index().Op("*").ID("types").Dot("Item"),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("exampleQuantity"), jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.ID("BuildFakeItem").Call(),
				)),
			jen.Return().Op("&").ID("types").Dot("ItemList").Valuesln(jen.ID("Pagination").Op(":").ID("types").Dot("Pagination").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("FilteredCount").Op(":").ID("exampleQuantity").Op("/").Lit(2), jen.ID("TotalCount").Op(":").ID("exampleQuantity")), jen.ID("Items").Op(":").ID("examples")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeItemUpdateInput builds a faked ItemUpdateInput from an item."),
		jen.Line(),
		jen.Func().ID("BuildFakeItemUpdateInput").Params().Params(jen.Op("*").ID("types").Dot("ItemUpdateInput")).Body(
			jen.ID("item").Op(":=").ID("BuildFakeItem").Call(),
			jen.Return().Op("&").ID("types").Dot("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("item").Dot("Name"), jen.ID("Details").Op(":").ID("item").Dot("Details"), jen.ID("BelongsToAccount").Op(":").ID("item").Dot("BelongsToAccount")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeItemUpdateInputFromItem builds a faked ItemUpdateInput from an item."),
		jen.Line(),
		jen.Func().ID("BuildFakeItemUpdateInputFromItem").Params(jen.ID("item").Op("*").ID("types").Dot("Item")).Params(jen.Op("*").ID("types").Dot("ItemUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("item").Dot("Name"), jen.ID("Details").Op(":").ID("item").Dot("Details"), jen.ID("BelongsToAccount").Op(":").ID("item").Dot("BelongsToAccount"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeItemCreationInput builds a faked ItemCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeItemCreationInput").Params().Params(jen.Op("*").ID("types").Dot("ItemCreationInput")).Body(
			jen.ID("item").Op(":=").ID("BuildFakeItem").Call(),
			jen.Return().ID("BuildFakeItemCreationInputFromItem").Call(jen.ID("item")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeItemCreationInputFromItem builds a faked ItemCreationInput from an item."),
		jen.Line(),
		jen.Func().ID("BuildFakeItemCreationInputFromItem").Params(jen.ID("item").Op("*").ID("types").Dot("Item")).Params(jen.Op("*").ID("types").Dot("ItemCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("ItemCreationInput").Valuesln(jen.ID("Name").Op(":").ID("item").Dot("Name"), jen.ID("Details").Op(":").ID("item").Dot("Details"), jen.ID("BelongsToAccount").Op(":").ID("item").Dot("BelongsToAccount"))),
		jen.Line(),
	)

	return code
}
