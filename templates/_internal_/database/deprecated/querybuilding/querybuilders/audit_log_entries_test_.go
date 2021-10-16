package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntriesTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestSqlite_BuildGetAuditLogEntryQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAllAuditLogEntriesCountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetBatchOfAuditLogEntriesQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAuditLogEntriesQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildCreateAuditLogEntryQuery(proj, dbvendor)...)

	return code
}

func buildTestSqlite_BuildGetAuditLogEntryQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAuditLogEntryQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAuditLogEntry").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntry").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.id = %s", getIncIndex(dbvendor, 0)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(
						jen.ID("exampleAuditLogEntry").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAuditLogEntryQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAuditLogEntry").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAllAuditLogEntriesCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAllAuditLogEntriesCountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(audit_log.id) FROM audit_log"),
					jen.ID("actualQuery").Assign().ID("q").Dot("BuildGetAllAuditLogEntriesCountQuery").Call(jen.ID("ctx")),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetBatchOfAuditLogEntriesQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetBatchOfAuditLogEntriesQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("beginID"), jen.ID("endID")).Assign().List(jen.Uint64().Call(jen.Lit(1)), jen.Uint64().Call(jen.Lit(1000))),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.id > %s AND audit_log.id < %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(
						jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetBatchOfAuditLogEntriesQuery").Call(
						jen.ID("ctx"),
						jen.ID("beginID"),
						jen.ID("endID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAuditLogEntriesQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAuditLogEntriesQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("filter").Assign().Qual(proj.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on, (SELECT COUNT(*) FROM audit_log) FROM audit_log WHERE audit_log.created_on > %s AND audit_log.created_on < %s AND audit_log.last_updated_on > %s AND audit_log.last_updated_on < %s ORDER BY audit_log.created_on LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(
						jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAuditLogEntriesQuery").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildCreateAuditLogEntryQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var querySuffix string
	if dbvendor.SingularPackageName() == "postgres" {
		querySuffix = " RETURNING id"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateAuditLogEntryQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAuditLogEntry").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntry").Call(),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryCreationInputFromAuditLogEntry").Call(jen.ID("exampleAuditLogEntry")),
					jen.Newline(),
					jen.ID("exIDGen").Assign().AddressOf().Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleAuditLogEntry").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Equals().ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("INSERT INTO audit_log (external_id,event_type,context) VALUES (%s,%s,%s)%s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), querySuffix),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(
						jen.ID("exampleAuditLogEntry").Dot("ExternalID"), jen.ID("exampleAuditLogEntry").Dot("EventType"), jen.ID("exampleAuditLogEntry").Dot("Context")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildCreateAuditLogEntryQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("exIDGen"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}
