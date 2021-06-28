package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(proj.QuerybuildingPackage(), "WebhookSQLQueryBuilder").Equals().Parens(jen.PointerTo().ID(dbvendor.Singular())).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(buildBuildGetWebhookQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAllWebhooksCountQuery(proj, dbvendor)...)
	code.Add(buildBuildGetBatchOfWebhooksQuery(proj, dbvendor)...)
	code.Add(buildBuildGetWebhooksQuery(proj, dbvendor)...)
	code.Add(buildBuildCreateWebhookQuery(proj, dbvendor)...)
	code.Add(buildBuildUpdateWebhookQuery(proj, dbvendor)...)
	code.Add(buildBuildArchiveWebhookQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAuditLogEntriesForWebhookQuery(proj, dbvendor)...)

	return code
}

func buildBuildGetWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("webhookID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableOwnershipColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllWebhooksCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAllWebhooksCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetBatchOfWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetBatchOfWebhooksQuery returns a query that fetches every webhook in the database within a bucketed range."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetBatchOfWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("beginID"))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetWebhooksQuery returns a SQL query (and arguments) that would return a list of webhooks."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetWebhooksQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").Uint64(), jen.ID("filter").PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.String().Call(jen.ID("filter").Dot("SortBy")),
				)),
			jen.Return().ID("b").Dot("buildListQuery").Callln(
				jen.ID("ctx"),
				jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName"),
				jen.ID("nil"),
				jen.ID("nil"),
				jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableOwnershipColumn"),
				jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableColumns"),
				jen.ID("accountID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildCreateWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildCreateWebhookQuery returns a SQL query (and arguments) that would create a given webhook."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildCreateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("x").PointerTo().Qual(proj.TypesPackage(), "WebhookCreationInput")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableNameColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableContentTypeColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableURLColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableMethodColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableOwnershipColumn"),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("x").Dot("Name"),
					jen.ID("x").Dot("ContentType"),
					jen.ID("x").Dot("URL"),
					jen.ID("x").Dot("Method"),
					jen.Qual("strings", "Join").Call(
						jen.ID("x").Dot("Events"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsSeparator"),
					),
					jen.Qual("strings", "Join").Call(
						jen.ID("x").Dot("DataTypes"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesSeparator"),
					),
					jen.Qual("strings", "Join").Call(
						jen.ID("x").Dot("Topics"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsSeparator"),
					),
					jen.ID("x").Dot("BelongsToAccount"),
				).Add(utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildingPackage(), "IDColumn"))))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUpdateWebhookQuery takes a given webhook and returns a SQL query to update."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildUpdateWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "Webhook")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("BelongsToAccount"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableNameColumn"),
					jen.ID("input").Dot("Name"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableContentTypeColumn"),
					jen.ID("input").Dot("ContentType"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableURLColumn"),
					jen.ID("input").Dot("URL"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableMethodColumn"),
					jen.ID("input").Dot("Method"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsColumn"),
					jen.Qual("strings", "Join").Call(
						jen.ID("input").Dot("Events"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableEventsSeparator"),
					),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesColumn"),
					jen.Qual("strings", "Join").Call(
						jen.ID("input").Dot("DataTypes"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableDataTypesSeparator"),
					),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsColumn"),
					jen.Qual("strings", "Join").Call(
						jen.ID("input").Dot("Topics"),
						jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableTopicsSeparator"),
					),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("input").Dot("ID"), jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableOwnershipColumn").Op(":").ID("input").Dot("BelongsToAccount"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildArchiveWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildArchiveWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("webhookID"), jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableOwnershipColumn").Op(":").ID("accountID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesForWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAuditLogEntriesForWebhookQuery constructs a SQL query for fetching audit log entries belong to a webhook with a given ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAuditLogEntriesForWebhookQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.Newline(),
			jen.ID("webhookIDKey").Assign().Qual("fmt", "Sprintf").Callln(
				jen.ID("jsonPluckQuery"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
				utils.ConditionalCode(dbvendor.SingularPackageName() == "mariadb", jen.ID("webhookID")),
				jen.Qual(proj.InternalAuditPackage(), "WebhookAssignmentKey"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				func() jen.Code {
					if dbvendor.SingularPackageName() == "mariadb" {
						return jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
							Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
							Dotln("Where").Call(jen.ID("squirrel").Dot("Expr").Call(jen.ID("webhookIDKey"))).
							Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
							jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
						))
					}
					return jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
						Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
						Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Values(jen.ID("webhookIDKey").Op(":").ID("webhookID"))).
						Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
					))
				}(),
			),
		),
		jen.Newline(),
	}

	return lines
}
