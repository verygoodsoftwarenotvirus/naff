package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
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

func buildTestSqlite_BuildGetWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleWebhook").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = ? AND webhooks.id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetWebhookQuery").Call(
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

func buildTestSqlite_BuildGetAllWebhooksCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAllWebhooksCountQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL"),
					jen.ID("actualQuery").Op(":=").ID("q").Dot("BuildGetAllWebhooksCountQuery").Call(jen.ID("ctx")),
					jen.Newline(),
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

func buildTestSqlite_BuildGetBatchOfWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetBatchOfWebhooksQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("beginID"), jen.ID("endID")).Op(":=").List(jen.ID("uint64").Call(jen.Lit(1)), jen.ID("uint64").Call(jen.Lit(1000))),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account FROM webhooks WHERE webhooks.id > ? AND webhooks.id < ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetBatchOfWebhooksQuery").Call(
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

func buildTestSqlite_BuildGetWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetWebhooksQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("filter").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account, (SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = ?) as total_count, (SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = ? AND webhooks.created_on > ? AND webhooks.created_on < ? AND webhooks.last_updated_on > ? AND webhooks.last_updated_on < ?) as filtered_count FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = ? AND webhooks.created_on > ? AND webhooks.created_on < ? AND webhooks.last_updated_on > ? AND webhooks.last_updated_on < ? GROUP BY webhooks.id LIMIT 20 OFFSET 180"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetWebhooksQuery").Call(
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

func buildTestSqlite_BuildCreateWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleWebhook").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.Newline(),
					jen.ID("exIDGen").Op(":=").Op("&").Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleWebhook").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO webhooks (external_id,name,content_type,url,method,events,data_types,topics,belongs_to_account) VALUES (?,?,?,?,?,?,?,?,?)"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("ExternalID"), jen.ID("exampleWebhook").Dot("Name"), jen.ID("exampleWebhook").Dot("ContentType"), jen.ID("exampleWebhook").Dot("URL"), jen.ID("exampleWebhook").Dot("Method"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Events"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("DataTypes"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Topics"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsSeparator"),
					), jen.ID("exampleWebhook").Dot("BelongsToAccount")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildCreateWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
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
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUpdateWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleWebhook").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET name = ?, content_type = ?, url = ?, method = ?, events = ?, data_types = ?, topics = ?, last_updated_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_account = ? AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("Name"), jen.ID("exampleWebhook").Dot("ContentType"), jen.ID("exampleWebhook").Dot("URL"), jen.ID("exampleWebhook").Dot("Method"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Events"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("DataTypes"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Topics"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsSeparator"),
					), jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUpdateWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook"),
					),
					jen.Newline(),
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

func buildTestSqlite_BuildArchiveWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildArchiveWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleWebhook").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_account = ? AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildArchiveWebhookQuery").Call(
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

func buildTestSqlite_BuildGetAuditLogEntriesForWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAuditLogEntriesForWebhookQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleWebhook").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE json_extract(audit_log.context, '$.webhook_id') = ? ORDER BY audit_log.created_on"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAuditLogEntriesForWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
					),
					jen.Newline(),
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
