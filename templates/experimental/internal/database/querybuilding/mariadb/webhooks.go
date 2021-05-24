package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("querybuilding").Dot("WebhookSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("MariaDB")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("BuildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("WebhooksTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("WebhooksTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("WebhooksTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("webhookID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("WebhooksTableName"),
					jen.ID("querybuilding").Dot("WebhooksTableOwnershipColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("WebhooksTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetAllWebhooksCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.ID("querybuilding").Dot("WebhooksTableName"),
				)).Dot("From").Call(jen.ID("querybuilding").Dot("WebhooksTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("WebhooksTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfWebhooksQuery returns a query that fetches every item in the database within a bucketed range."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetBatchOfWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("WebhooksTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("WebhooksTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("WebhooksTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("beginID"))).Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("WebhooksTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("BuildGetWebhooksQuery returns a SQL query (and arguments) that would return a query and arguments to retrieve a list of webhooks.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				)),
			jen.Return().ID("b").Dot("buildListQuery").Call(
				jen.ID("ctx"),
				jen.ID("querybuilding").Dot("WebhooksTableName"),
				jen.ID("querybuilding").Dot("WebhooksTableOwnershipColumn"),
				jen.ID("querybuilding").Dot("WebhooksTableColumns"),
				jen.ID("accountID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("BuildCreateWebhookQuery returns a SQL query (and arguments) that would create a given webhook.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildCreateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("x").Op("*").ID("types").Dot("WebhookCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("WebhooksTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("ExternalIDColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableNameColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableContentTypeColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableURLColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableMethodColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableEventsColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableDataTypesColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableTopicsColumn"),
					jen.ID("querybuilding").Dot("WebhooksTableOwnershipColumn"),
				).Dot("Values").Call(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("x").Dot("Name"),
					jen.ID("x").Dot("ContentType"),
					jen.ID("x").Dot("URL"),
					jen.ID("x").Dot("Method"),
					jen.Qual("strings", "Join").Call(
						jen.ID("x").Dot("Events"),
						jen.ID("querybuilding").Dot("WebhooksTableEventsSeparator"),
					),
					jen.Qual("strings", "Join").Call(
						jen.ID("x").Dot("DataTypes"),
						jen.ID("querybuilding").Dot("WebhooksTableDataTypesSeparator"),
					),
					jen.Qual("strings", "Join").Call(
						jen.ID("x").Dot("Topics"),
						jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
					),
					jen.ID("x").Dot("BelongsToAccount"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateWebhookQuery takes a given webhook and returns a SQL query to update."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildUpdateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Webhook")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("BelongsToAccount"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("WebhooksTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("WebhooksTableNameColumn"),
					jen.ID("input").Dot("Name"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("WebhooksTableContentTypeColumn"),
					jen.ID("input").Dot("ContentType"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("WebhooksTableURLColumn"),
					jen.ID("input").Dot("URL"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("WebhooksTableMethodColumn"),
					jen.ID("input").Dot("Method"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("WebhooksTableEventsColumn"),
					jen.Qual("strings", "Join").Call(
						jen.ID("input").Dot("Events"),
						jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
					),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("WebhooksTableDataTypesColumn"),
					jen.Qual("strings", "Join").Call(
						jen.ID("input").Dot("DataTypes"),
						jen.ID("querybuilding").Dot("WebhooksTableDataTypesSeparator"),
					),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("WebhooksTableTopicsColumn"),
					jen.Qual("strings", "Join").Call(
						jen.ID("input").Dot("Topics"),
						jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
					),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("input").Dot("ID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"), jen.ID("querybuilding").Dot("WebhooksTableOwnershipColumn").Op(":").ID("input").Dot("BelongsToAccount"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("BuildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildArchiveWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("WebhooksTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("webhookID"), jen.ID("querybuilding").Dot("WebhooksTableOwnershipColumn").Op(":").ID("accountID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesForWebhookQuery constructs a SQL query for fetching an audit log entry with a given ID belong to a user with a given ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetAuditLogEntriesForWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Expr").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("jsonPluckQuery"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableContextColumn"),
					jen.ID("webhookID"),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookAssignmentKey"),
				))).Dot("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("CreatedOnColumn"),
				)),
			),
		),
		jen.Line(),
	)

	return code
}
