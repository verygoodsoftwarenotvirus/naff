package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntryDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AuditLogEntryDataManager").Op("=").Parens(jen.Op("*").ID("AuditLogEntryDataManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AuditLogEntryDataManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogUserBanEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogUserBanEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("banGiver"), jen.ID("banReceiver")).ID("uint64"), jen.ID("reason").ID("string")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("banGiver"),
				jen.ID("banReceiver"),
				jen.ID("reason"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogAccountTerminationEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogAccountTerminationEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("adminID"), jen.ID("accountID")).ID("uint64"), jen.ID("reason").ID("string")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("adminID"),
				jen.ID("accountID"),
				jen.ID("reason"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntry is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("GetAuditLogEntry").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("entryID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("AuditLogEntry")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAuditLogEntriesCount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("GetAllAuditLogEntriesCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAuditLogEntries is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("GetAllAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")).Body(
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
		jen.Comment("GetAuditLogEntries is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("GetAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("AuditLogEntryList"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("AuditLogEntryList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
