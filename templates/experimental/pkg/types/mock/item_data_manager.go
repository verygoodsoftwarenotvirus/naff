package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("ItemDataManager").Op("=").Parens(jen.Op("*").ID("ItemDataManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("ItemDataManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ItemExists is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("ItemExists").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetItem is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("Item")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllItemsCount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("GetAllItemsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllItems is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("GetAllItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").ID("types").Dot("Item"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("results"),
				jen.ID("bucketSize"),
			),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetItems is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("ItemList"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("ItemList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetItemsWithIDs is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("GetItemsWithIDs").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("limit"),
				jen.ID("ids"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Op("*").ID("types").Dot("Item")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateItem is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("ItemCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("createdByUser"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("Item")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateItem is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("Item"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("types").Dot("FieldChangeSummary")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
				jen.ID("changedByUser"),
				jen.ID("changes"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveItem is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("belongsToAccount"), jen.ID("archivedBy")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("belongsToAccount"),
				jen.ID("archivedBy"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntriesForItem is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ItemDataManager")).ID("GetAuditLogEntriesForItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Op("*").ID("types").Dot("AuditLogEntry")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
