package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntriesTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestSqlite_BuildGetAuditLogEntryQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAuditLogEntry").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntry").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAuditLogEntry").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAuditLogEntryQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAuditLogEntry").Dot("ID"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestSqlite_BuildGetAllAuditLogEntriesCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(audit_log.id) FROM audit_log"),
					jen.ID("actualQuery").Op(":=").ID("q").Dot("BuildGetAllAuditLogEntriesCountQuery").Call(jen.ID("ctx")),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestSqlite_BuildGetBatchOfAuditLogEntriesQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("beginID"), jen.ID("endID")).Op(":=").List(jen.ID("uint64").Call(jen.Lit(1)), jen.ID("uint64").Call(jen.Lit(1000))),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.id > ? AND audit_log.id < ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetBatchOfAuditLogEntriesQuery").Call(
						jen.ID("ctx"),
						jen.ID("beginID"),
						jen.ID("endID"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestSqlite_BuildGetAuditLogEntriesQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("filter").Op(":=").ID("fakes").Dot("BuildFleshedOutQueryFilter").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on, (SELECT COUNT(*) FROM audit_log) FROM audit_log WHERE audit_log.created_on > ? AND audit_log.created_on < ? AND audit_log.last_updated_on > ? AND audit_log.last_updated_on < ? ORDER BY audit_log.created_on LIMIT 20 OFFSET 180"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAuditLogEntriesQuery").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestSqlite_BuildCreateAuditLogEntryQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAuditLogEntry").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntry").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryCreationInputFromAuditLogEntry").Call(jen.ID("exampleAuditLogEntry")),
					jen.ID("exIDGen").Op(":=").Op("&").ID("querybuilding").Dot("MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleAuditLogEntry").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO audit_log (external_id,event_type,context) VALUES (?,?,?)"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAuditLogEntry").Dot("ExternalID"), jen.ID("exampleAuditLogEntry").Dot("EventType"), jen.ID("exampleAuditLogEntry").Dot("Context")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildCreateAuditLogEntryQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("exIDGen"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}

func buildTestSqlite_BuildGetAuditLogEntryQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		//
	}

	return lines
}

func buildTestSqlite_BuildGetAllAuditLogEntriesCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		//
	}

	return lines
}

func buildTestSqlite_BuildGetBatchOfAuditLogEntriesQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		//
	}

	return lines
}

func buildTestSqlite_BuildGetAuditLogEntriesQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		//
	}

	return lines
}

func buildTestSqlite_BuildCreateAuditLogEntryQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		//
	}

	return lines
}
