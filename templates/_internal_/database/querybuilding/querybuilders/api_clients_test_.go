package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestSqlite_BuildGetBatchOfAPIClientsQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAPIClientQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAllAPIClientsCountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAPIClientsQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAPIClientByDatabaseIDQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildCreateAPIClientQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildUpdateAPIClientQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildArchiveAPIClientQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAuditLogEntriesForAPIClientQuery(proj, dbvendor)...)

	return code
}

func buildTestSqlite_BuildGetBatchOfAPIClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetBatchOfAPIClientsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("beginID"), jen.ID("endID")).Op(":=").List(jen.ID("uint64").Call(jen.Lit(1)), jen.ID("uint64").Call(jen.Lit(1000))),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT api_clients.id, api_clients.external_id, api_clients.name, api_clients.client_id, api_clients.secret_key, api_clients.created_on, api_clients.last_updated_on, api_clients.archived_on, api_clients.belongs_to_user FROM api_clients WHERE api_clients.id > ? AND api_clients.id < ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetBatchOfAPIClientsQuery").Call(
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
	}

	return lines
}

func buildTestSqlite_BuildGetAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetAPIClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT api_clients.id, api_clients.external_id, api_clients.name, api_clients.client_id, api_clients.secret_key, api_clients.created_on, api_clients.last_updated_on, api_clients.archived_on, api_clients.belongs_to_user FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.client_id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ClientID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAPIClientByClientIDQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ClientID"),
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
	}

	return lines
}

func buildTestSqlite_BuildGetAllAPIClientsCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetAllAPIClientsCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(api_clients.id) FROM api_clients WHERE api_clients.archived_on IS NULL"),
					jen.ID("actualQuery").Op(":=").ID("q").Dot("BuildGetAllAPIClientsCountQuery").Call(jen.ID("ctx")),
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
	}

	return lines
}

func buildTestSqlite_BuildGetAPIClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetAPIClientsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("filter").Op(":=").ID("fakes").Dot("BuildFleshedOutQueryFilter").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT api_clients.id, api_clients.external_id, api_clients.name, api_clients.client_id, api_clients.secret_key, api_clients.created_on, api_clients.last_updated_on, api_clients.archived_on, api_clients.belongs_to_user, (SELECT COUNT(api_clients.id) FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.belongs_to_user = ?) as total_count, (SELECT COUNT(api_clients.id) FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.belongs_to_user = ? AND api_clients.created_on > ? AND api_clients.created_on < ? AND api_clients.last_updated_on > ? AND api_clients.last_updated_on < ?) as filtered_count FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.belongs_to_user = ? AND api_clients.created_on > ? AND api_clients.created_on < ? AND api_clients.last_updated_on > ? AND api_clients.last_updated_on < ? GROUP BY api_clients.id LIMIT 20 OFFSET 180"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAPIClientsQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
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
	}

	return lines
}

func buildTestSqlite_BuildGetAPIClientByDatabaseIDQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetAPIClientByDatabaseIDQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT api_clients.id, api_clients.external_id, api_clients.name, api_clients.client_id, api_clients.secret_key, api_clients.created_on, api_clients.last_updated_on, api_clients.archived_on, api_clients.belongs_to_user FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.belongs_to_user = ? AND api_clients.id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAPIClient").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAPIClientByDatabaseIDQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
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
	}

	return lines
}

func buildTestSqlite_BuildCreateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildCreateAPIClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("exampleAPIClientInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
					jen.ID("exIDGen").Op(":=").Op("&").ID("querybuilding").Dot("MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleAPIClient").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO api_clients (external_id,name,client_id,secret_key,belongs_to_user) VALUES (?,?,?,?,?)"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ExternalID"), jen.ID("exampleAPIClient").Dot("Name"), jen.ID("exampleAPIClient").Dot("ClientID"), jen.ID("exampleAPIClient").Dot("ClientSecret"), jen.ID("exampleAPIClient").Dot("BelongsToUser")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildCreateAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClientInput"),
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
	}

	return lines
}

func buildTestSqlite_BuildUpdateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildUpdateAPIClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE api_clients SET client_id = ?, last_updated_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ClientID"), jen.ID("exampleAPIClient").Dot("BelongsToUser"), jen.ID("exampleAPIClient").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUpdateAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient"),
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
	}

	return lines
}

func buildTestSqlite_BuildArchiveAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildArchiveAPIClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE api_clients SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAPIClient").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildArchiveAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
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
	}

	return lines
}

func buildTestSqlite_BuildGetAuditLogEntriesForAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetAuditLogEntriesForAPIClientQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE json_extract(audit_log.context, '$.api_client_id') = ? ORDER BY audit_log.created_on"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAuditLogEntriesForAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ID"),
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
	}

	return lines
}
