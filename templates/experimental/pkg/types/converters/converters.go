package converters

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func convertersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("ConvertAuditLogEntryCreationInputToEntry converts an AuditLogEntryCreationInput to an AuditLogEntry."),
		jen.Line(),
		jen.Func().ID("ConvertAuditLogEntryCreationInputToEntry").Params(jen.ID("e").Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Params(jen.Op("*").ID("types").Dot("AuditLogEntry")).Body(
			jen.Return().Op("&").ID("types").Dot("AuditLogEntry").Valuesln(jen.ID("EventType").Op(":").ID("e").Dot("EventType"), jen.ID("Context").Op(":").ID("e").Dot("Context"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ConvertAccountToAccountUpdateInput creates an AccountUpdateInput struct from an item."),
		jen.Line(),
		jen.Func().ID("ConvertAccountToAccountUpdateInput").Params(jen.ID("x").Op("*").ID("types").Dot("Account")).Params(jen.Op("*").ID("types").Dot("AccountUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AccountUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("x").Dot("Name"), jen.ID("BelongsToUser").Op(":").ID("x").Dot("BelongsToUser"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ConvertItemToItemUpdateInput creates an ItemUpdateInput struct from an item."),
		jen.Line(),
		jen.Func().ID("ConvertItemToItemUpdateInput").Params(jen.ID("x").Op("*").ID("types").Dot("Item")).Params(jen.Op("*").ID("types").Dot("ItemUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("x").Dot("Name"), jen.ID("Details").Op(":").ID("x").Dot("Details"))),
		jen.Line(),
	)

	return code
}
