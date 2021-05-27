package fakes

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntryDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeAuditLogEntry builds a faked item."),
		jen.Line(),
		jen.Func().ID("BuildFakeAuditLogEntry").Params().Params(jen.Op("*").ID("types").Dot("AuditLogEntry")).Body(
			jen.Return().Op("&").ID("types").Dot("AuditLogEntry").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("ExternalID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "SuccessfulLoginEvent"), jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(jen.Lit("fakes").Op(":").Lit("true")), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAuditLogEntryList builds a faked AuditLogEntryList."),
		jen.Line(),
		jen.Func().ID("BuildFakeAuditLogEntryList").Params().Params(jen.Op("*").ID("types").Dot("AuditLogEntryList")).Body(
			jen.Var().Defs(
				jen.ID("examples").Index().Op("*").ID("types").Dot("AuditLogEntry"),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("exampleQuantity"), jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.ID("BuildFakeAuditLogEntry").Call(),
				)),
			jen.Return().Op("&").ID("types").Dot("AuditLogEntryList").Valuesln(jen.ID("Pagination").Op(":").ID("types").Dot("Pagination").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").ID("exampleQuantity")), jen.ID("Entries").Op(":").ID("examples")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAuditLogEntryCreationInput builds a faked AuditLogEntryCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeAuditLogEntryCreationInput").Params().Params(jen.Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Body(
			jen.ID("item").Op(":=").ID("BuildFakeAuditLogEntry").Call(),
			jen.Return().ID("BuildFakeAuditLogEntryCreationInputFromAuditLogEntry").Call(jen.ID("item")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAuditLogEntryCreationInputFromAuditLogEntry builds a faked AuditLogEntryCreationInput from an item."),
		jen.Line(),
		jen.Func().ID("BuildFakeAuditLogEntryCreationInputFromAuditLogEntry").Params(jen.ID("entry").Op("*").ID("types").Dot("AuditLogEntry")).Params(jen.Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln(jen.ID("EventType").Op(":").ID("entry").Dot("EventType"), jen.ID("Context").Op(":").ID("entry").Dot("Context"))),
		jen.Line(),
	)

	return code
}
