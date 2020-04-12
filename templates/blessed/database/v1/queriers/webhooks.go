package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	ret := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("eventsSeparator").Equals().RawString(`,`),
			jen.ID("typesSeparator").Equals().RawString(`,`),
			jen.ID("topicsSeparator").Equals().RawString(`,`),
			jen.Line(),
			jen.ID("webhooksTableName").Equals().Lit("webhooks"),
			jen.ID("webhooksTableOwnershipColumn").Equals().Lit("belongs_to_user"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("webhooksTableColumns").Equals().Index().String().Valuesln(
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.name"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.content_type"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.url"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.method"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.events"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.data_types"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.topics"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.created_on"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.updated_on"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("webhooksTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.ID("webhooksTableName"), jen.ID("webhooksTableOwnershipColumn")),
			),
		),
		jen.Line(),
	)

	ret.Add(buildScanWebhook(proj, dbvendor)...)
	ret.Add(buildScanWebhooks(proj, dbvendor)...)
	ret.Add(buildBuildGetWebhookQuery(proj, dbvendor)...)
	ret.Add(build_GetWebhook(proj, dbvendor)...)
	ret.Add(buildBuildGetAllWebhooksCountQuery(proj, dbvendor)...)
	ret.Add(build_GetAllWebhooksCount(proj, dbvendor)...)
	ret.Add(buildBuildGetAllWebhooksQuery(proj, dbvendor)...)
	ret.Add(build_GetAllWebhooks(proj, dbvendor)...)
	ret.Add(buildBuildGetWebhooksQuery(proj, dbvendor)...)
	ret.Add(buildGetWebhooks(proj, dbvendor)...)
	ret.Add(buildBuildWebhookCreationQuery(proj, dbvendor)...)
	ret.Add(buildCreateWebhook(proj, dbvendor)...)
	ret.Add(buildBuildUpdateWebhookQuery(proj, dbvendor)...)
	ret.Add(buildUpdateWebhook(proj, dbvendor)...)
	ret.Add(buildBuildArchiveWebhookQuery(proj, dbvendor)...)
	ret.Add(buildArchiveWebhook(proj, dbvendor)...)

	return ret
}

func buildScanWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("scanWebhook is a consistent way to turn a *sql.Row into a webhook struct"),
		jen.Line(),
		jen.Func().ID("scanWebhook").Params(
			jen.ID("scan").Qual(proj.DatabaseV1Package(), "Scanner"),
			jen.ID("includeCount").Bool(),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Uint64(),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("x").Equals().AddressOf().Qual(proj.ModelsV1Package(), "Webhook").Values(),
				jen.ID("count").Uint64(),
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
				jen.AddressOf().ID("x").Dot("UpdatedOn"),
				jen.AddressOf().ID("x").Dot("ArchivedOn"),
				jen.AddressOf().ID("x").Dot("BelongsToUser"),
			),
			jen.Line(),
			jen.If(jen.ID("includeCount")).Block(
				utils.AppendItemsToList(jen.ID("targetVars"), jen.AddressOf().ID("count")),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Spread()), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.ID("events").Assign().Qual("strings", "Split").Call(jen.ID("eventsStr"), jen.ID("eventsSeparator")), jen.ID("len").Call(jen.ID("events")).Op(">=").One().And().ID("events").Index(jen.Zero()).DoesNotEqual().EmptyString()).Block(
				jen.ID("x").Dot("Events").Equals().ID("events"),
			),
			jen.If(jen.ID("dataTypes").Assign().Qual("strings", "Split").Call(jen.ID("dataTypesStr"), jen.ID("typesSeparator")), jen.ID("len").Call(jen.ID("dataTypes")).Op(">=").One().And().ID("dataTypes").Index(jen.Zero()).DoesNotEqual().EmptyString()).Block(
				jen.ID("x").Dot("DataTypes").Equals().ID("dataTypes"),
			),
			jen.If(jen.ID("topics").Assign().Qual("strings", "Split").Call(jen.ID("topicsStr"), jen.ID("topicsSeparator")), jen.ID("len").Call(jen.ID("topics")).Op(">=").One().And().ID("topics").Index(jen.Zero()).DoesNotEqual().EmptyString()).Block(
				jen.ID("x").Dot("Topics").Equals().ID("topics"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("count"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildScanWebhooks(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks"),
		jen.Line(),
		jen.Func().ID("scanWebhooks").Params(
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("rows").PointerTo().Qual("database/sql", "Rows"),
		).Params(
			jen.Index().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Uint64(),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("list").Index().Qual(proj.ModelsV1Package(), "Webhook"),
				jen.ID("count").Uint64(),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("webhook"), jen.ID("c"), jen.Err()).Assign().ID("scanWebhook").Call(jen.ID("rows"), jen.True()),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Zero(), jen.Err()),
				),
				jen.Line(),
				jen.If(jen.ID("count").IsEqualTo().Zero()).Block(
					jen.ID("count").Equals().ID("c"),
				),
				jen.Line(),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.PointerTo().ID("webhook")),
			),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("count"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetWebhookQuery").Params(jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("webhooksTableColumns").Spread()).
				Dotln("From").Call(jen.ID("webhooksTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.ID("webhooksTableName")).MapAssign().ID("webhookID"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("webhooksTableName"),
					jen.ID("webhooksTableOwnershipColumn"),
				).MapAssign().ID("userID"),
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
		jen.Comment("GetWebhook fetches a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhookQuery").Call(jen.ID("webhookID"), jen.ID("userID")),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("webhook"), jen.Underscore(), jen.Err()).Assign().ID("scanWebhook").Call(jen.ID("row"), jen.False()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Lit("querying for webhook"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("webhook"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetAllWebhooksCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
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
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllWebhooksCountQuery").Params().Params(jen.String()).Block(
			jen.ID("getAllWebhooksCountQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("getAllWebhooksCountQuery"), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(jen.ID("countQuery"), jen.ID("webhooksTableName"))).
					Dotln("From").Call(jen.ID("webhooksTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("webhooksTableName")).MapAssign().ID("nil"))).
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

func build_GetAllWebhooksCount(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetAllWebhooksCount will fetch the count of every active webhook in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllWebhooksCount").Params(utils.CtxParam()).Params(jen.ID("count").Uint64(), jen.Err().Error()).Block(
			jen.Err().Equals().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID(dbfl).Dot("buildGetAllWebhooksCountQuery").Call()).Dot("Scan").Call(jen.AddressOf().ID("count")),
			jen.Return().List(jen.ID("count"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetAllWebhooksQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("getAllWebhooksQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllWebhooksQuery").String(),
		),
		jen.Line(),
		jen.Line(),
		jen.Comment("buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllWebhooksQuery").Params().Params(jen.String()).Block(
			jen.ID("getAllWebhooksQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("getAllWebhooksQuery"), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("webhooksTableColumns").Spread()).
					Dotln("From").Call(jen.ID("webhooksTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("webhooksTableName")).MapAssign().ID("nil"))).
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
		jen.Comment("GetAllWebhooks fetches a list of all webhooks from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllWebhooks").Params(utils.CtxParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Block(
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID(dbfl).Dot("buildGetAllWebhooksQuery").Call()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for webhooks: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("count"), jen.Err()).Assign().ID("scanWebhooks").Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "WebhookList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().One(),
					jen.ID("TotalCount").MapAssign().ID("count"),
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

func build_GetAllWebhooksForUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetAllWebhooksForUser fetches a list of all webhooks from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllWebhooksForUser").Params(utils.CtxParam(), jen.ID("userID").Uint64()).Params(jen.Index().Qual(proj.ModelsV1Package(), "Webhook"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhooksQuery").Call(jen.ID("userID"), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for webhooks: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().ID("scanWebhooks").Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
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
			jen.ID("userID").Uint64(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.Append(
				jen.ID("webhooksTableColumns"),
				jen.Qual("fmt", "Sprintf").Call(jen.ID("countQuery"), jen.ID("webhooksTableName")),
			).Spread()).
				Dotln("From").Call(jen.ID("webhooksTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("webhooksTableName"),
					jen.ID("webhooksTableOwnershipColumn"),
				).MapAssign().ID("userID"),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("webhooksTableName")).MapAssign().ID("nil")),
			).
				Dotln("GroupBy").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.ID("webhooksTableName"))),
			jen.Line(),
			jen.If(jen.ID(utils.FilterVarName).DoesNotEqual().ID("nil")).Block(
				jen.ID("builder").Equals().ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder"), jen.ID("webhooksTableName")),
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
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetWebhooks").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetWebhooksQuery").Call(jen.ID("userID"), jen.ID(utils.FilterVarName)),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("count"), jen.Err()).Assign().ID("scanWebhooks").Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "WebhookList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID("filter").Dot("Page"),
					jen.ID("TotalCount").MapAssign().ID("count"),
					jen.ID("Limit").MapAssign().ID("filter").Dot("Limit"),
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
			jen.Lit("name"),
			jen.Lit("content_type"),
			jen.Lit("url"),
			jen.Lit("method"),
			jen.Lit("events"),
			jen.Lit("data_types"),
			jen.Lit("topics"),
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
			jen.ID("x").Dot("BelongsToUser"),
		}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Insert").Call(jen.ID("webhooksTableName")).
			Dotln("Columns").Callln(cols...).
			Dotln("Values").Callln(vals...)

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	lines := []jen.Code{
		jen.Comment("buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildWebhookCreationQuery").Params(jen.ID("x").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
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
				jen.ID("BelongsToUser").MapAssign().ID("input").Dot("BelongsToUser"),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildWebhookCreationQuery").Call(jen.ID("x")),
		}

		if isPostgres(dbvendor) {
			out = append(out,
				jen.If(jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("ID"), jen.AddressOf().ID("x").Dot("CreatedOn")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing webhook creation query: %w"), jen.Err())),
				),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			out = append(out,
				jen.List(jen.ID("res"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return(jen.List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing webhook creation query: %w"), jen.Err()))),
				),
				jen.Line(),
				jen.Comment("fetch the last inserted ID"),
				jen.List(jen.ID("id"), jen.ID("err")).Assign().ID("res").Dot("LastInsertId").Call(),
				jen.ID(dbfl).Dot("logIDRetrievalError").Call(jen.Err()),
				jen.ID("x").Dot("ID").Equals().Uint64().Call(jen.ID("id")),
				jen.Line(),
				jen.Comment("this won't be completely accurate, but it will suffice"),
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
		jen.Comment("CreateWebhook creates a webhook in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("CreateWebhook").Params(utils.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"), jen.Error()).Block(
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
			Dotln("Set").Call(jen.Lit("name"), jen.ID("input").Dot("Name")).
			Dotln("Set").Call(jen.Lit("content_type"), jen.ID("input").Dot("ContentType")).
			Dotln("Set").Call(jen.Lit("url"), jen.ID("input").Dot("URL")).
			Dotln("Set").Call(jen.Lit("method"), jen.ID("input").Dot("Method")).
			Dotln("Set").Call(jen.Lit("events"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Events"), jen.ID("topicsSeparator"))).
			Dotln("Set").Call(jen.Lit("data_types"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("DataTypes"), jen.ID("typesSeparator"))).
			Dotln("Set").Call(jen.Lit("topics"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Topics"), jen.ID("topicsSeparator"))).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").MapAssign().ID("input").Dot("ID"),
			jen.ID("webhooksTableOwnershipColumn").MapAssign().ID("input").Dot("BelongsToUser")),
		)

		if isPostgres(dbvendor) {
			q.Dot("Suffix").Call(jen.Lit("RETURNING updated_on"))
		}

		q.Dotln("ToSql").Call()

		return q
	}

	lines := []jen.Code{
		jen.Comment("buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildUpdateWebhookQuery").Params(jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
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
				jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("input").Dot("UpdatedOn")),
			}
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateWebhookQuery").Call(jen.ID("input")),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.Return(jen.Err()),
			}
		}
		return nil
	}

	lines := []jen.Code{
		jen.Comment("UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("UpdateWebhook").Params(utils.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.Error()).Block(
			buildUpdateWebhookBody()...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveWebhookQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildArchiveWebhookQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("webhooksTableName")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").MapAssign().ID("webhookID"),
			jen.ID("webhooksTableOwnershipColumn").MapAssign().ID("userID"),
			jen.Lit("archived_on").MapAssign().ID("nil"),
		))

		if isPostgres(dbvendor) {
			q.Dot("Suffix").Call(jen.Lit("RETURNING archived_on"))
		}

		q.Dotln("ToSql").Call()

		return q
	}

	lines := []jen.Code{
		jen.Comment("buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildArchiveWebhookQuery").Params(jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
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

func buildArchiveWebhook(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("ArchiveWebhook archives a webhook from the database by its ID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("ArchiveWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveWebhookQuery").Call(jen.ID("webhookID"), jen.ID("userID")),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return().Err(),
		),
		jen.Line(),
	}

	return lines
}
