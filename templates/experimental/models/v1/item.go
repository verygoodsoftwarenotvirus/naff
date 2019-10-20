package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(
		jen.Type().ID("Item").Struct(jen.ID("ID").ID("uint64"), jen.ID("Name").ID("string"), jen.ID("Details").ID("string"), jen.ID("CreatedOn").ID("uint64"), jen.ID("UpdatedOn").Op("*").ID("uint64"), jen.ID("ArchivedOn").Op("*").ID("uint64"), jen.ID("BelongsTo").ID("uint64")).Type().ID("ItemList").Struct(jen.ID("Pagination"), jen.ID("Items").Index().ID("Item")).Type().ID("ItemCreationInput").Struct(jen.ID("Name").ID("string"), jen.ID("Details").ID("string"), jen.ID("BelongsTo").ID("uint64")).Type().ID("ItemUpdateInput").Struct(jen.ID("Name").ID("string"), jen.ID("Details").ID("string"), jen.ID("BelongsTo").ID("uint64")).Type().ID("ItemDataManager").Interface(jen.ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("Item"), jen.ID("error")), jen.ID("GetItemCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")), jen.ID("GetAllItemsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")), jen.ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("ItemList"), jen.ID("error")), jen.ID("GetAllItemsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("Item"), jen.ID("error")), jen.ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("ItemCreationInput")).Params(jen.Op("*").ID("Item"), jen.ID("error")), jen.ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("Item")).Params(jen.ID("error")), jen.ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("id"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error"))).Type().ID("ItemDataServer").Interface(jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc"))),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Update merges an ItemInput with an Item"),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("Item")).ID("Update").Params(jen.ID("input").Op("*").ID("ItemUpdateInput")).Block(
			jen.If(jen.ID("input").Dot("Name").Op("!=").Lit("").Op("||").ID("input").Dot("Name").Op("!=").ID("x").Dot("Name")).Block(
				jen.ID("x").Dot("Name").Op("=").ID("input").Dot("Name"),
			),
			jen.If(jen.ID("input").Dot(
				"Details",
			).Op("!=").Lit("").Op("||").ID("input").Dot(
				"Details",
			).Op("!=").ID("x").Dot(
				"Details",
			)).Block(
				jen.ID("x").Dot(
					"Details",
				).Op("=").ID("input").Dot(
					"Details",
				),
			),
		),
		jen.Line(),
	)
	return ret
}
