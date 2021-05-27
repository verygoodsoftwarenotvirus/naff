package builders

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestPostgres_BuildGetWebhookQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = $1 AND webhooks.id = $2"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleWebhook").Dot("BelongsToAccount"),
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
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPostgres_BuildGetAllWebhooksCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL"),
					jen.ID("actualQuery").Op(":=").ID("q").Dot("BuildGetAllWebhooksCountQuery").Call(jen.ID("ctx")),
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
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPostgres_BuildGetBatchOfWebhooksQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("beginID"), jen.ID("endID")).Op(":=").List(jen.ID("uint64").Call(jen.Lit(1)), jen.ID("uint64").Call(jen.Lit(1000))),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account FROM webhooks WHERE webhooks.id > $1 AND webhooks.id < $2"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetBatchOfWebhooksQuery").Call(
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
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPostgres_BuildGetWebhooksQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("filter").Op(":=").ID("fakes").Dot("BuildFleshedOutQueryFilter").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT webhooks.id, webhooks.external_id, webhooks.name, webhooks.content_type, webhooks.url, webhooks.method, webhooks.events, webhooks.data_types, webhooks.topics, webhooks.created_on, webhooks.last_updated_on, webhooks.archived_on, webhooks.belongs_to_account, (SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = $1) as total_count, (SELECT COUNT(webhooks.id) FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = $2 AND webhooks.created_on > $3 AND webhooks.created_on < $4 AND webhooks.last_updated_on > $5 AND webhooks.last_updated_on < $6) as filtered_count FROM webhooks WHERE webhooks.archived_on IS NULL AND webhooks.belongs_to_account = $7 AND webhooks.created_on > $8 AND webhooks.created_on < $9 AND webhooks.last_updated_on > $10 AND webhooks.last_updated_on < $11 GROUP BY webhooks.id LIMIT 20 OFFSET 180"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetWebhooksQuery").Call(
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
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPostgres_BuildWebhookCreationQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.ID("exIDGen").Op(":=").Op("&").ID("querybuilding").Dot("MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleWebhook").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO webhooks (external_id,name,content_type,url,method,events,data_types,topics,belongs_to_account) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("ExternalID"), jen.ID("exampleWebhook").Dot("Name"), jen.ID("exampleWebhook").Dot("ContentType"), jen.ID("exampleWebhook").Dot("URL"), jen.ID("exampleWebhook").Dot("Method"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Events"),
						jen.ID("querybuilding").Dot("WebhooksTableEventsSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("DataTypes"),
						jen.ID("querybuilding").Dot("WebhooksTableDataTypesSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Topics"),
						jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
					), jen.ID("exampleWebhook").Dot("BelongsToAccount")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildCreateWebhookQuery").Call(
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
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPostgres_BuildUpdateWebhookQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET name = $1, content_type = $2, url = $3, method = $4, events = $5, data_types = $6, topics = $7, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $8 AND id = $9"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("Name"), jen.ID("exampleWebhook").Dot("ContentType"), jen.ID("exampleWebhook").Dot("URL"), jen.ID("exampleWebhook").Dot("Method"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Events"),
						jen.ID("querybuilding").Dot("WebhooksTableEventsSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("DataTypes"),
						jen.ID("querybuilding").Dot("WebhooksTableDataTypesSeparator"),
					), jen.Qual("strings", "Join").Call(
						jen.ID("exampleWebhook").Dot("Topics"),
						jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
					), jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUpdateWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook"),
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
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPostgres_BuildArchiveWebhookQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE webhooks SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_account = $1 AND id = $2"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("BelongsToAccount"), jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildArchiveWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleWebhook").Dot("BelongsToAccount"),
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
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPostgres_BuildGetAuditLogEntriesForWebhookQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE audit_log.context->'webhook_id' = $1 ORDER BY audit_log.created_on"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleWebhook").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAuditLogEntriesForWebhookQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
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
		jen.Line(),
	)

	return code
}
