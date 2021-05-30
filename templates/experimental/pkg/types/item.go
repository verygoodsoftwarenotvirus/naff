package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ItemsSearchIndexName").ID("search").Dot("IndexName").Op("=").Lit("items"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Item").Struct(
				jen.ID("ArchivedOn").Op("*").ID("uint64"),
				jen.ID("LastUpdatedOn").Op("*").ID("uint64"),
				jen.ID("ExternalID").ID("string"),
				jen.ID("Name").ID("string"),
				jen.ID("Details").ID("string"),
				jen.ID("CreatedOn").ID("uint64"),
				jen.ID("ID").ID("uint64"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("ItemList").Struct(
				jen.ID("Items").Index().Op("*").ID("Item"),
				jen.ID("Pagination"),
			),
			jen.ID("ItemCreationInput").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("Details").ID("string"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("ItemUpdateInput").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("Details").ID("string"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("ItemDataManager").Interface(
				jen.ID("ItemExists").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("bool"), jen.ID("error")),
				jen.ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.Op("*").ID("Item"), jen.ID("error")),
				jen.ID("GetAllItemsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetAllItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("resultChannel").Chan().Index().Op("*").ID("Item"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")),
				jen.ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("ItemList"), jen.ID("error")),
				jen.ID("GetItemsWithIDs").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("uint64")).Params(jen.Index().Op("*").ID("Item"), jen.ID("error")),
				jen.ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("ItemCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("Item"), jen.ID("error")),
				jen.ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("Item"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("FieldChangeSummary")).Params(jen.ID("error")),
				jen.ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("belongsToAccount"), jen.ID("archivedByUserID")).ID("uint64")).Params(jen.ID("error")),
				jen.ID("GetAuditLogEntriesForItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.Index().Op("*").ID("AuditLogEntry"), jen.ID("error")),
			),
			jen.ID("ItemDataService").Interface(
				jen.ID("SearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ExistenceHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("UpdateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Update merges an ItemUpdateInput with an item."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("Item")).ID("Update").Params(jen.ID("input").Op("*").ID("ItemUpdateInput")).Params(jen.Index().Op("*").ID("FieldChangeSummary")).Body(
			jen.Var().Defs(
				jen.ID("out").Index().Op("*").ID("FieldChangeSummary"),
			),
			jen.If(jen.ID("input").Dot("Name").Op("!=").Lit("").Op("&&").ID("input").Dot("Name").Op("!=").ID("x").Dot("Name")).Body(
				jen.ID("out").Op("=").ID("append").Call(
					jen.ID("out"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("Name"), jen.ID("OldValue").Op(":").ID("x").Dot("Name"), jen.ID("NewValue").Op(":").ID("input").Dot("Name")),
				),
				jen.ID("x").Dot("Name").Op("=").ID("input").Dot("Name"),
			),
			jen.If(jen.ID("input").Dot("Details").Op("!=").Lit("").Op("&&").ID("input").Dot("Details").Op("!=").ID("x").Dot("Details")).Body(
				jen.ID("out").Op("=").ID("append").Call(
					jen.ID("out"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("Details"), jen.ID("OldValue").Op(":").ID("x").Dot("Details"), jen.ID("NewValue").Op(":").ID("input").Dot("Details")),
				),
				jen.ID("x").Dot("Details").Op("=").ID("input").Dot("Details"),
			),
			jen.Return().ID("out"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("ItemCreationInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a ItemUpdateInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("ItemUpdateInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Name"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("ItemUpdateInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a ItemUpdateInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("ItemUpdateInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Name"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	return code
}
