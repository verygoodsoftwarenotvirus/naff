package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockAuditLogEntrySQLQueryBuilderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("AuditLogEntrySQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("AuditLogEntrySQLQueryBuilder")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AuditLogEntrySQLQueryBuilder").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntryQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntrySQLQueryBuilder")).ID("BuildGetAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("entryID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllAuditLogEntriesCountQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntrySQLQueryBuilder")).ID("BuildGetAllAuditLogEntriesCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().ID("returnArgs").Dot("String").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfAuditLogEntriesQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntrySQLQueryBuilder")).ID("BuildGetBatchOfAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("beginID"),
				jen.ID("endID"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntrySQLQueryBuilder")).ID("BuildGetAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateAuditLogEntryQuery implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntrySQLQueryBuilder")).ID("BuildCreateAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Index().Interface())),
		),
		jen.Line(),
	)

	return code
}
