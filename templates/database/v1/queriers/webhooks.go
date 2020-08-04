package queriers

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	code.Add(buildWebhooksConstDeclarations()...)
	code.Add(buildWebhooksVarDeclarations()...)
	code.Add(buildScanWebhook(proj, dbvendor)...)
	code.Add(buildScanWebhooks(proj, dbvendor)...)
	code.Add(buildBuildGetWebhookQuery(dbvendor)...)
	code.Add(build_GetWebhook(proj, dbvendor)...)
	code.Add(buildBuildGetAllWebhooksCountQuery(dbvendor)...)
	code.Add(build_GetAllWebhooksCount(dbvendor)...)
	code.Add(buildBuildGetAllWebhooksQuery(dbvendor)...)
	code.Add(build_GetAllWebhooks(proj, dbvendor)...)
	code.Add(buildBuildGetWebhooksQuery(proj, dbvendor)...)
	code.Add(buildGetWebhooks(proj, dbvendor)...)
	code.Add(buildBuildWebhookCreationQuery(proj, dbvendor)...)
	code.Add(buildCreateWebhook(proj, dbvendor)...)
	code.Add(buildBuildUpdateWebhookQuery(proj, dbvendor)...)
	code.Add(buildUpdateWebhook(proj, dbvendor)...)
	code.Add(buildBuildArchiveWebhookQuery(dbvendor)...)
	code.Add(buildArchiveWebhook(dbvendor)...)

	return code
}

func buildWebhooksConstDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("commaSeparator").Equals().Lit(","),
			jen.Line(),
			jen.ID("eventsSeparator").Equals().ID("commaSeparator"),
			jen.ID("typesSeparator").Equals().ID("commaSeparator"),
			jen.ID("topicsSeparator").Equals().ID("commaSeparator"),
			jen.Line(),
			jen.ID("webhooksTableName").Equals().Lit("webhooks"),
			jen.ID("webhooksTableNameColumn").Equals().Lit("name"),
			jen.ID("webhooksTableContentTypeColumn").Equals().Lit("content_type"),
			jen.ID("webhooksTableURLColumn").Equals().Lit("url"),
			jen.ID("webhooksTableMethodColumn").Equals().Lit("method"),
			jen.ID("webhooksTableEventsColumn").Equals().Lit("events"),
			jen.ID("webhooksTableDataTypesColumn").Equals().Lit("data_types"),
			jen.ID("webhooksTableTopicsColumn").Equals().Lit("topics"),
			jen.ID("webhooksTableOwnershipColumn").Equals().Lit("belongs_to_user"),
		),
		jen.Line(),
	}

	return lines
}

func buildWebhooksVarDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("webhooksTableColumns").Equals().Index().String().Valuesln(
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("idColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableNameColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableContentTypeColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableURLColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableMethodColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableEventsColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableDataTypesColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableTopicsColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("createdOnColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("lastUpdatedOnColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("archivedOnColumn")),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("webhooksTableOwnershipColumn")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildScanWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("scanWebhook is a consistent way to turn a *sql.Row into a webhook struct."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("scanWebhook").Params(
			jen.ID("scan").Qual(proj.DatabaseV1Package(), "Scanner"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error(),
		).Body(
			jen.Var().Defs(
				jen.ID("x").Equals().AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Values(),
				jen.Listln(
					jen.ID("eventsStr"),
					jen.ID("dataTypesStr"),
					jen.ID("topicsStr").String(),
				),
			),
			jen.Line(),
			jen.ID("targetVars").Assign().Index().Interface().Valuesln(
				jen.AddressOf().ID("x").Dot("ID"),
				jen.AddressOf().ID("x").Dot("Name"),
				jen.AddressOf().ID("x").Dot("ContentType"),
				jen.AddressOf().ID("x").Dot("URL"),
				jen.AddressOf().ID("x").Dot("Method"),
				jen.AddressOf().ID("eventsStr"),
				jen.AddressOf().ID("dataTypesStr"),
				jen.AddressOf().ID("topicsStr"),
				jen.AddressOf().ID("x").Dot("CreatedOn"),
				jen.AddressOf().ID("x").Dot("LastUpdatedOn"),
				jen.AddressOf().ID("x").Dot("ArchivedOn"),
				jen.AddressOf().ID("x").Dot(constants.UserOwnershipFieldName),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Spread()), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.ID("events").Assign().Qual("strings", "Split").Call(jen.ID("eventsStr"), jen.ID("eventsSeparator")), jen.Len(jen.ID("events")).Op(">=").One().And().ID("events").Index(jen.Zero()).DoesNotEqual().EmptyString()).Body(
				jen.ID("x").Dot("Events").Equals().ID("events"),
			),
			jen.If(jen.ID("dataTypes").Assign().Qual("strings", "Split").Call(jen.ID("dataTypesStr"), jen.ID("typesSeparator")), jen.Len(jen.ID("dataTypes")).Op(">=").One().And().ID("dataTypes").Index(jen.Zero()).DoesNotEqual().EmptyString()).Body(
				jen.ID("x").Dot("DataTypes").Equals().ID("dataTypes"),
			),
			jen.If(jen.ID("topics").Assign().Qual("strings", "Split").Call(jen.ID("topicsStr"), jen.ID("topicsSeparator")), jen.Len(jen.ID("topics")).Op(">=").One().And().ID("topics").Index(jen.Zero()).DoesNotEqual().EmptyString()).Body(
				jen.ID("x").Dot("Topics").Equals().ID("topics"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildScanWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("scanWebhooks").Params(
			jen.ID("rows").Qual(proj.DatabaseV1Package(), "ResultIterator"),
		).Params(
			jen.Index().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error(),
		).Body(
			jen.Var().Defs(
				jen.ID("list").Index().Qual(proj.ModelsV1Package(), "Webhook"),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("webhook"), jen.Err()).Assign().ID(dbfl).Dot("scanWebhook").Call(jen.ID("rows")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Line(),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.PointerTo().ID("webhook")),
			),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetWebhookQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetWebhookQuery").Params(jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("webhooksTableColumns").Spread()).
				Dotln("From").Call(jen.ID("webhooksTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("idColumn")).MapAssign().ID("webhookID"),
				utils.FormatString("%s.%s",
					jen.ID("webhooksTableName"),
					jen.ID("webhooksTableOwnershipColumn"),
				).MapAssign().ID(constants.UserIDVarName),
			)).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func build_GetWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetWebhook fetches a webhook from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhookQuery").Call(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("webhook"), jen.Err()).Assign().ID(dbfl).Dot("scanWebhook").Call(jen.ID("row")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Lit("querying for webhook"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("webhook"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetAllWebhooksCountQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("getAllWebhooksCountQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllWebhooksCountQuery").String(),
		),
		jen.Line(),
		jen.Line(),
		jen.Comment("buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllWebhooksCountQuery").Params().Params(jen.String()).Body(
			jen.ID("getAllWebhooksCountQueryBuilder").Dot("Do").Call(jen.Func().Params().Body(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("getAllWebhooksCountQuery"), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(utils.FormatStringWithArg(jen.ID("countQuery"), jen.ID("webhooksTableName"))).
					Dotln("From").Call(jen.ID("webhooksTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("archivedOnColumn")).MapAssign().ID("nil"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			)),
			jen.Line(),
			jen.Return().ID("getAllWebhooksCountQuery"),
		),
		jen.Line(),
	}

	return lines
}

func build_GetAllWebhooksCount(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetAllWebhooksCount will fetch the count of every active webhook in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllWebhooksCount").Params(constants.CtxParam()).Params(jen.ID("count").Uint64(), jen.Err().Error()).Body(
			jen.Err().Equals().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID(dbfl).Dot("buildGetAllWebhooksCountQuery").Call()).Dot("Scan").Call(jen.AddressOf().ID("count")),
			jen.Return().List(jen.ID("count"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetAllWebhooksQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("getAllWebhooksQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllWebhooksQuery").String(),
		),
		jen.Line(),
		jen.Line(),
		jen.Comment("buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllWebhooksQuery").Params().Params(jen.String()).Body(
			jen.ID("getAllWebhooksQueryBuilder").Dot("Do").Call(jen.Func().Params().Body(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("getAllWebhooksQuery"), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("webhooksTableColumns").Spread()).
					Dotln("From").Call(jen.ID("webhooksTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("archivedOnColumn")).MapAssign().ID("nil"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			)),
			jen.Line(),
			jen.Return().ID("getAllWebhooksQuery"),
		),
		jen.Line(),
	}

	return lines
}

func build_GetAllWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetAllWebhooks fetches a list of all webhooks from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllWebhooks").Params(constants.CtxParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Body(
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID(dbfl).Dot("buildGetAllWebhooksQuery").Call()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Body(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for webhooks: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().ID(dbfl).Dot("scanWebhooks").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "WebhookList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().One(),
				),
				jen.ID("Webhooks").MapAssign().ID("list"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetWebhooksQuery returns a SQL query (and arguments) that would return a"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetWebhooksQuery").Params(
			constants.UserIDParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("webhooksTableColumns").Spread()).
				Dotln("From").Call(jen.ID("webhooksTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				utils.FormatString("%s.%s",
					jen.ID("webhooksTableName"),
					jen.ID("webhooksTableOwnershipColumn"),
				).MapAssign().ID(constants.UserIDVarName),
				utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("archivedOnColumn")).MapAssign().ID("nil")),
			).
				Dotln("OrderBy").Call(utils.FormatString("%s.%s", jen.ID("webhooksTableName"), jen.ID("idColumn"))),
			jen.Line(),
			jen.If(jen.ID(constants.FilterVarName).DoesNotEqual().ID("nil")).Body(
				jen.ID("builder").Equals().ID(constants.FilterVarName).Dot("ApplyToQueryBuilder").Call(jen.ID("builder"), jen.ID("webhooksTableName")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildGetWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetWebhooks").Params(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhooksQuery").Call(jen.ID(constants.UserIDVarName), jen.ID(constants.FilterVarName)),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Body(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().ID(dbfl).Dot("scanWebhooks").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "WebhookList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID(constants.FilterVarName).Dot("Page"),
					jen.ID("Limit").MapAssign().ID(constants.FilterVarName).Dot("Limit"),
				),
				jen.ID("Webhooks").MapAssign().ID("list"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildWebhookCreationQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildWebhookCreationQueryQuery := func() jen.Code {
		cols := []jen.Code{
			jen.ID("webhooksTableNameColumn"),
			jen.ID("webhooksTableContentTypeColumn"),
			jen.ID("webhooksTableURLColumn"),
			jen.ID("webhooksTableMethodColumn"),
			jen.ID("webhooksTableEventsColumn"),
			jen.ID("webhooksTableDataTypesColumn"),
			jen.ID("webhooksTableTopicsColumn"),
			jen.ID("webhooksTableOwnershipColumn"),
		}
		vals := []jen.Code{
			jen.ID("x").Dot("Name"),
			jen.ID("x").Dot("ContentType"),
			jen.ID("x").Dot("URL"),
			jen.ID("x").Dot("Method"),
			jen.Qual("strings", "Join").Call(jen.ID("x").Dot("Events"), jen.ID("eventsSeparator")),
			jen.Qual("strings", "Join").Call(jen.ID("x").Dot("DataTypes"), jen.ID("typesSeparator")),
			jen.Qual("strings", "Join").Call(jen.ID("x").Dot("Topics"), jen.ID("topicsSeparator")),
			jen.ID("x").Dot(constants.UserOwnershipFieldName),
		}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Insert").Call(jen.ID("webhooksTableName")).
			Dotln("Columns").Callln(cols...).
			Dotln("Values").Callln(vals...)

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s, %s"), jen.ID("idColumn"), jen.ID("createdOnColumn")))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	lines := []jen.Code{
		jen.Comment("buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildWebhookCreationQuery").Params(jen.ID("x").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			buildWebhookCreationQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildCreateWebhookQuery := func() []jen.Code {
		out := []jen.Code{
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Valuesln(
				jen.ID("Name").MapAssign().ID("input").Dot("Name"),
				jen.ID("ContentType").MapAssign().ID("input").Dot("ContentType"),
				jen.ID("URL").MapAssign().ID("input").Dot("URL"),
				jen.ID("Method").MapAssign().ID("input").Dot("Method"),
				jen.ID("Events").MapAssign().ID("input").Dot("Events"),
				jen.ID("DataTypes").MapAssign().ID("input").Dot("DataTypes"),
				jen.ID("Topics").MapAssign().ID("input").Dot("Topics"),
				jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("input").Dot(constants.UserOwnershipFieldName),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildWebhookCreationQuery").Call(jen.ID("x")),
		}

		if isPostgres(dbvendor) {
			out = append(out,
				jen.If(jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("ID"), jen.AddressOf().ID("x").Dot("CreatedOn")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing webhook creation query: %w"), jen.Err())),
				),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			out = append(out,
				jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return(jen.List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing webhook creation query: %w"), jen.Err()))),
				),
				jen.Line(),
				jen.Comment("fetch the last inserted ID."),
				jen.List(jen.ID("id"), jen.ID("err")).Assign().ID(constants.ResponseVarName).Dot("LastInsertId").Call(),
				jen.ID(dbfl).Dot("logIDRetrievalError").Call(jen.Err()),
				jen.ID("x").Dot("ID").Equals().Uint64().Call(jen.ID("id")),
				jen.Line(),
				jen.Comment("this won't be completely accurate, but it will suffice."),
				jen.ID("x").Dot("CreatedOn").Equals().ID(dbfl).Dot("timeTeller").Dot("Now").Call(),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		)

		return out
	}

	lines := []jen.Code{
		jen.Comment("CreateWebhook creates a webhook in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("CreateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"), jen.Error()).Body(
			buildCreateWebhookQuery()...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildUpdateWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildUpdateWebhookQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("webhooksTableName")).
			Dotln("Set").Call(jen.ID("webhooksTableNameColumn"), jen.ID("input").Dot("Name")).
			Dotln("Set").Call(jen.ID("webhooksTableContentTypeColumn"), jen.ID("input").Dot("ContentType")).
			Dotln("Set").Call(jen.ID("webhooksTableURLColumn"), jen.ID("input").Dot("URL")).
			Dotln("Set").Call(jen.ID("webhooksTableMethodColumn"), jen.ID("input").Dot("Method")).
			Dotln("Set").Call(jen.ID("webhooksTableEventsColumn"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Events"), jen.ID("topicsSeparator"))).
			Dotln("Set").Call(jen.ID("webhooksTableDataTypesColumn"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("DataTypes"), jen.ID("typesSeparator"))).
			Dotln("Set").Call(jen.ID("webhooksTableTopicsColumn"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Topics"), jen.ID("topicsSeparator"))).
			Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.ID("idColumn").MapAssign().ID("input").Dot("ID"),
			jen.ID("webhooksTableOwnershipColumn").MapAssign().ID("input").Dot(constants.UserOwnershipFieldName)),
		)

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.ID("lastUpdatedOnColumn")))
		}

		q.Dotln("ToSql").Call()

		return q
	}

	lines := []jen.Code{
		jen.Comment("buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildUpdateWebhookQuery").Params(jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			buildUpdateWebhookQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildUpdateWebhookBody := func() []jen.Code {
		if isPostgres(dbvendor) {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateWebhookQuery").Call(jen.ID("input")),
				jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("input").Dot("LastUpdatedOn")),
			}
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateWebhookQuery").Call(jen.ID("input")),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.Return(jen.Err()),
			}
		}
		panic(fmt.Sprintf("invalid dbvendor: %q", dbvendor))
	}

	lines := []jen.Code{
		jen.Comment("UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("UpdateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.Error()).Body(
			buildUpdateWebhookBody()...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveWebhookQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildArchiveWebhookQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("webhooksTableName")).
			Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Set").Call(jen.ID("archivedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.ID("idColumn").MapAssign().ID("webhookID"),
			jen.ID("webhooksTableOwnershipColumn").MapAssign().ID(constants.UserIDVarName),
			jen.ID("archivedOnColumn").MapAssign().ID("nil"),
		))

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.ID("archivedOnColumn")))
		}

		q.Dotln("ToSql").Call()

		return q
	}

	lines := []jen.Code{
		jen.Comment("buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildArchiveWebhookQuery").Params(jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			buildArchiveWebhookQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveWebhook(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("ArchiveWebhook archives a webhook from the database by its ID."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("ArchiveWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveWebhookQuery").Call(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return().Err(),
		),
		jen.Line(),
	}

	return lines
}
