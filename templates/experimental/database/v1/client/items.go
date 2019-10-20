package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").ID("models").Dot(
			"ItemDataManager",
		).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachItemIDToSpan provides a consistent way to attach an item's ID to a span"),
		jen.Line(),
		jen.Func().ID("attachItemIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("itemID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("item_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("itemID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetItem fetches an item from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("models").Dot("Item"),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetItem")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachItemIDToSpan").Call(jen.ID("span"), jen.ID("itemID")),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("item_id").Op(":").ID("itemID"), jen.Lit("user_id").Op(":").ID("userID"))).Dot("Debug").Call(jen.Lit("GetItem called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetItem",
			).Call(jen.ID("ctx"), jen.ID("itemID"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetItemCount fetches the count of items from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetItemCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot("QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetItemCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetItemCount called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetItemCount",
			).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllItemsCount fetches the count of items from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllItemsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllItemsCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllItemsCount called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"GetAllItemsCount",
			).Call(jen.ID("ctx")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetItems fetches a list of items from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot("QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot(
			"ItemList",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetItems")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetItems called")),
			jen.List(jen.ID("itemList"), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dot(
				"GetItems",
			).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Return().List(jen.ID("itemList"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllItemsForUser fetches a list of items from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllItemsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("models").Dot("Item"),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllItemsForUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetAllItemsForUser called")),
			jen.List(jen.ID("itemList"), jen.ID("err")).Op(":=").ID("c").Dot("querier").Dot(
				"GetAllItemsForUser",
			).Call(jen.ID("ctx"), jen.ID("userID")),
			jen.Return().List(jen.ID("itemList"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateItem creates an item in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
			"ItemCreationInput",
		)).Params(jen.Op("*").ID("models").Dot("Item"),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("CreateItem")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")).Dot("Debug").Call(jen.Lit("CreateItem called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"CreateItem",
			).Call(jen.ID("ctx"), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateItem updates a particular item. Note that UpdateItem expects the"),
		jen.Line(),
		jen.Comment("provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot("Item")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("UpdateItem")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachItemIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("ID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("item_id"), jen.ID("input").Dot("ID")).Dot("Debug").Call(jen.Lit("UpdateItem called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"UpdateItem",
			).Call(jen.ID("ctx"), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveItem archives an item from the database by its ID"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("ArchiveItem")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachItemIDToSpan").Call(jen.ID("span"), jen.ID("itemID")),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("item_id").Op(":").ID("itemID"), jen.Lit("user_id").Op(":").ID("userID"))).Dot("Debug").Call(jen.Lit("ArchiveItem called")),
			jen.Return().ID("c").Dot("querier").Dot(
				"ArchiveItem",
			).Call(jen.ID("ctx"), jen.ID("itemID"), jen.ID("userID")),
		),
		jen.Line(),
	)
	return ret
}
