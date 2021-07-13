package querybuilders

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

func webhooksTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestSqlite_BuildGetWebhookQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAllWebhooksCountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetBatchOfWebhooksQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetWebhooksQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildCreateWebhookQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildUpdateWebhookQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildArchiveWebhookQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAuditLogEntriesForWebhookQuery(proj, dbvendor)...)

	return code
}

var webhooksTableColumns = []string{
	"webhooks.id",
	"webhooks.external_id",
	"webhooks.name",
	"webhooks.content_type",
	"webhooks.url",
	"webhooks.method",
	"webhooks.events",
	"webhooks.data_types",
	"webhooks.topics",
	"webhooks.created_on",
	"webhooks.last_updated_on",
	"webhooks.archived_on",
	"webhooks.belongs_to_account",
}

func buildTestSqlite_BuildGetWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select(webhooksTableColumns...).
			From("webhooks").
			Where(squirrel.Eq{
				"webhooks.id":                 whateverValue,
				"webhooks.belongs_to_account": whateverValue,
				"webhooks.archived_on":        nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleWebhook").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleWebhook").Dot("BelongsToAccount"),
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

func buildTestSqlite_BuildGetAllWebhooksCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAllWebhooksCountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL"),
					jen.ID("actualQuery").Assign().ID("q").Dot("BuildGetAllWebhooksCountQuery").Call(jen.ID("ctx")),
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

func buildTestSqlite_BuildGetBatchOfWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetBatchOfWebhooksQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Litf("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account FROM webhooks WHERE webhooks.id > %s AND webhooks.id < %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetBatchOfWebhooksQuery").Call(
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

func buildTestSqlite_BuildGetWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetWebhooksQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Litf("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account, (SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = %s) as total_count, (SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = %s AND webhooks.created_on > %s AND webhooks.created_on < %s AND webhooks.last_updated_on > %s AND webhooks.last_updated_on < %s) as filtered_count FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = %s AND webhooks.created_on > %s AND webhooks.created_on < %s AND webhooks.last_updated_on > %s AND webhooks.last_updated_on < %s GROUP BY webhooks.id LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7), getIncIndex(dbvendor, 8), getIncIndex(dbvendor, 9), getIncIndex(dbvendor, 10)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetWebhooksQuery").Call(
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

func buildTestSqlite_BuildCreateWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var querySuffix string
	if dbvendor.SingularPackageName() == "postgres" {
		querySuffix = " RETURNING id"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleWebhook").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.Newline(),
					jen.ID("exIDGen").Assign().AddressOf().Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleWebhook").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Equals().ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("INSERT INTO webhooks (external_id,name,content_type,url,method,events,data_types,topics,belongs_to_account) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s)%s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7), getIncIndex(dbvendor, 8), querySuffix),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("ExternalID"), jen.ID("exampleWebhook").Dot("Name"), jen.ID("exampleWebhook").Dot("ContentType"), jen.ID("exampleWebhook").Dot("URL"), jen.ID("exampleWebhook").Dot("Method"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Events"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("DataTypes"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Topics"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsSeparator"),
					), jen.ID("exampleWebhook").Dot("BelongsToAccount")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildCreateWebhookQuery").Call(
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

func buildTestSqlite_BuildUpdateWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("webhooks").
			Set("name", whateverValue).
			Set("content_type", whateverValue).
			Set("url", whateverValue).
			Set("method", whateverValue).
			Set("events", whateverValue).
			Set("data_types", whateverValue).
			Set("topics", whateverValue).
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":                 whateverValue,
				"belongs_to_account": whateverValue,
				"archived_on":        nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUpdateWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleWebhook").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("Name"), jen.ID("exampleWebhook").Dot("ContentType"), jen.ID("exampleWebhook").Dot("URL"), jen.ID("exampleWebhook").Dot("Method"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Events"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("DataTypes"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Topics"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsSeparator"),
					), jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUpdateWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook"),
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

func buildTestSqlite_BuildArchiveWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("webhooks").
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":                 whateverValue,
				"belongs_to_account": whateverValue,
				"archived_on":        nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildArchiveWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleWebhook").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildArchiveWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleWebhook").Dot("BelongsToAccount"),
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

func buildTestSqlite_BuildGetAuditLogEntriesForWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var webhookKey string

	switch dbvendor.LowercaseAbbreviation() {
	case "m":
		webhookKey = fmt.Sprintf(`JSON_CONTAINS(%s.%s, '%s', '$.%s')`, "audit_log", "context", "%d", "webhook_id")
	case "p":
		webhookKey = fmt.Sprintf(`%s.%s->'%s'`, "audit_log", "context", "webhook_id")
	case "s":
		webhookKey = fmt.Sprintf(`json_extract(%s.%s, '$.%s')`, "audit_log", "context", "webhook_id")
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
		queryBuilder = queryBuilder.Where(squirrel.Expr(webhookKey))
	} else {
		queryBuilder = queryBuilder.Where(squirrel.Eq{webhookKey: whateverValue})
	}

	queryBuilder = queryBuilder.OrderBy("audit_log.created_on")

	expectedQuery, _ := buildQuery(queryBuilder)

	expectedQueryDecl := jen.ID("expectedQuery").Assign().Lit(expectedQuery)
	expectedArgs := jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("ID"))

	if dbvendor.SingularPackageName() == "mariadb" {
		expectedQueryDecl = jen.ID("expectedQuery").Assign().Qual("fmt", "Sprintf").Call(jen.Lit(expectedQuery), jen.ID("exampleWebhook").Dot("ID"))
		expectedArgs = jen.ID("expectedArgs").Assign().Index().Interface().Call(jen.Nil())
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAuditLogEntriesForWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleWebhook").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					expectedQueryDecl,
					expectedArgs,
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAuditLogEntriesForWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
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
