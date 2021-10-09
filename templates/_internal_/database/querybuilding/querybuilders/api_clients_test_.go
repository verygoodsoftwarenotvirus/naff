package querybuilders

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select(
			"api_clients.id",
			"api_clients.external_id",
			"api_clients.name",
			"api_clients.client_id",
			"api_clients.secret_key",
			"api_clients.created_on",
			"api_clients.last_updated_on",
			"api_clients.archived_on",
			"api_clients.belongs_to_user",
		).
			From("api_clients").
			Where(squirrel.Gt{"api_clients.id": whateverValue}).
			Where(squirrel.Lt{"api_clients.id": whateverValue}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetBatchOfAPIClientsQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetBatchOfAPIClientsQuery").Call(
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

func buildTestSqlite_BuildGetAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select(
			"api_clients.id",
			"api_clients.external_id",
			"api_clients.name",
			"api_clients.client_id",
			"api_clients.secret_key",
			"api_clients.created_on",
			"api_clients.last_updated_on",
			"api_clients.archived_on",
			"api_clients.belongs_to_user",
		).
			From("api_clients").
			Where(squirrel.Eq{
				"api_clients.client_id":   whateverValue,
				"api_clients.archived_on": nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAPIClientQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleAPIClient").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAPIClient").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ClientID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAPIClientByClientIDQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ClientID"),
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

func buildTestSqlite_BuildGetAllAPIClientsCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(columnCountQueryTemplate, "api_clients")).
		From("api_clients").
		Where(squirrel.Eq{
			"api_clients.archived_on": nil,
		}))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAllAPIClientsCountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("actualQuery").Assign().ID("q").Dot("BuildGetAllAPIClientsCountQuery").Call(jen.ID("ctx")),
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

func buildTestSqlite_BuildGetAPIClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAPIClientsQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("filter").Assign().Qual(proj.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("SELECT api_clients.id, api_clients.external_id, api_clients.name, api_clients.client_id, api_clients.secret_key, api_clients.created_on, api_clients.last_updated_on, api_clients.archived_on, api_clients.belongs_to_user, (SELECT COUNT(api_clients.id) FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.belongs_to_user = %s) as total_count, (SELECT COUNT(api_clients.id) FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.belongs_to_user = %s AND api_clients.created_on > %s AND api_clients.created_on < %s AND api_clients.last_updated_on > %s AND api_clients.last_updated_on < %s) as filtered_count FROM api_clients WHERE api_clients.archived_on IS NULL AND api_clients.belongs_to_user = %s AND api_clients.created_on > %s AND api_clients.created_on < %s AND api_clients.last_updated_on > %s AND api_clients.last_updated_on < %s GROUP BY api_clients.id LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7), getIncIndex(dbvendor, 8), getIncIndex(dbvendor, 9), getIncIndex(dbvendor, 10)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAPIClientsQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
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

func buildTestSqlite_BuildGetAPIClientByDatabaseIDQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select(
			"api_clients.id",
			"api_clients.external_id",
			"api_clients.name",
			"api_clients.client_id",
			"api_clients.secret_key",
			"api_clients.created_on",
			"api_clients.last_updated_on",
			"api_clients.archived_on",
			"api_clients.belongs_to_user",
		).
			From("api_clients").
			Where(squirrel.Eq{
				"api_clients.belongs_to_user": whateverValue,
				"api_clients.id":              whateverValue,
				"api_clients.archived_on":     nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAPIClientByDatabaseIDQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAPIClient").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAPIClient").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAPIClient").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAPIClientByDatabaseIDQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
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

func buildTestSqlite_BuildCreateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).Insert("api_clients").
		Columns(
			"external_id",
			"name",
			"client_id",
			"secret_key",
			"belongs_to_user",
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
		)

	if dbvendor.SingularPackageName() == "postgres" {
		qb = qb.Suffix("RETURNING id")
	}

	expectedQuery, _ := buildQuery(qb)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateAPIClientQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleAPIClient").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAPIClient").Call(),
					jen.ID("exampleAPIClientInput").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
					jen.Newline(),
					jen.ID("exIDGen").Assign().AddressOf().Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleAPIClient").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Equals().ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ExternalID"), jen.ID("exampleAPIClient").Dot("Name"), jen.ID("exampleAPIClient").Dot("ClientID"), jen.ID("exampleAPIClient").Dot("ClientSecret"), jen.ID("exampleAPIClient").Dot("BelongsToUser")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildCreateAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClientInput"),
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

func buildTestSqlite_BuildUpdateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("api_clients").
			Set("client_id", whateverValue).
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":              whateverValue,
				"belongs_to_user": whateverValue,
				"archived_on":     nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUpdateAPIClientQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleAPIClient").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAPIClient").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ClientID"), jen.ID("exampleAPIClient").Dot("BelongsToUser"), jen.ID("exampleAPIClient").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUpdateAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient"),
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

func buildTestSqlite_BuildArchiveAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("api_clients").
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":              whateverValue,
				"archived_on":     nil,
				"belongs_to_user": whateverValue,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildArchiveAPIClientQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAPIClient").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAPIClient").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAPIClient").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildArchiveAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
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

func buildTestSqlite_BuildGetAuditLogEntriesForAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var apiClientIDKey string

	switch dbvendor.LowercaseAbbreviation() {
	case "m":
		apiClientIDKey = fmt.Sprintf(`JSON_CONTAINS(%s.%s, '%s', '$.%s')`, "audit_log", "context", "%d", "api_client_id")
	case "p":
		apiClientIDKey = fmt.Sprintf(`%s.%s->'%s'`, "audit_log", "context", "api_client_id")
	case "s":
		apiClientIDKey = fmt.Sprintf(`json_extract(%s.%s, '$.%s')`, "audit_log", "context", "api_client_id")
	}

	queryBuilder := queryBuilderForDatabase(dbvendor).Select(
		"audit_log.id",
		"audit_log.external_id",
		"audit_log.event_type",
		"audit_log.context",
		"audit_log.created_on",
	).
		From("audit_log")

	if dbvendor.SingularPackageName() == "mariadb" {
		queryBuilder = queryBuilder.Where(squirrel.Expr(apiClientIDKey))
	} else {
		queryBuilder = queryBuilder.Where(squirrel.Eq{apiClientIDKey: whateverValue})
	}

	queryBuilder = queryBuilder.OrderBy("audit_log.created_on")

	expectedQuery, _ := buildQuery(queryBuilder)

	expectedQueryDecl := jen.ID("expectedQuery").Assign().Lit(expectedQuery)
	expectedArgs := jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAPIClient").Dot("ID"))

	if dbvendor.SingularPackageName() == "mariadb" {
		expectedQueryDecl = jen.ID("expectedQuery").Assign().Qual("fmt", "Sprintf").Call(jen.Lit(expectedQuery), jen.ID("exampleAPIClient").Dot("ID"))
		expectedArgs = jen.ID("expectedArgs").Assign().Index().Interface().Call(jen.Nil())
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAuditLogEntriesForAPIClientQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleAPIClient").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAPIClient").Call(),
					jen.Newline(),
					expectedQueryDecl,
					expectedArgs,
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAuditLogEntriesForAPIClientQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleAPIClient").Dot("ID"),
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
