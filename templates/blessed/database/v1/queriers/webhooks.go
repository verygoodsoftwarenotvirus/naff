package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(pkgRoot string, types []models.DataType, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(pkgRoot, types, ret)
	sn := vendor.Singular()
	dbrn := vendor.RouteName()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Const().Defs(
			jen.ID("eventsSeparator").Op("=").RawString(`,`),
			jen.ID("typesSeparator").Op("=").RawString(`,`),
			jen.ID("topicsSeparator").Op("=").RawString(`,`),
			jen.Line(),
			jen.ID("webhooksTableName").Op("=").Lit("webhooks"),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("webhooksTableColumns").Op("=").Index().ID("string").Valuesln(
				jen.Lit("id"),
				jen.Lit("name"),
				jen.Lit("content_type"),
				jen.Lit("url"),
				jen.Lit("method"),
				jen.Lit("events"),
				jen.Lit("data_types"),
				jen.Lit("topics"),
				jen.Lit("created_on"),
				jen.Lit("updated_on"),
				jen.Lit("archived_on"),
				jen.Lit("belongs_to"),
			),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("scanWebhook is a consistent way to turn a *sql.Row into a webhook struct"),
		jen.Line(),
		jen.Func().ID("scanWebhook").Params(jen.ID("scan").Qual(filepath.Join(pkgRoot, "database/v1"), "Scanner")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook"), jen.ID("error")).Block(
			jen.Var().Defs(
				jen.ID("x").Op("=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook").Values(),
				jen.Listln(
					jen.ID("eventsStr"),
					jen.ID("dataTypesStr"),
					jen.ID("topicsStr").ID("string"),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("scan").Dot("Scan").Callln(
				jen.Op("&").ID("x").Dot("ID"),
				jen.Op("&").ID("x").Dot("Name"),
				jen.Op("&").ID("x").Dot("ContentType"),
				jen.Op("&").ID("x").Dot("URL"),
				jen.Op("&").ID("x").Dot("Method"),
				jen.Op("&").ID("eventsStr"),
				jen.Op("&").ID("dataTypesStr"),
				jen.Op("&").ID("topicsStr"),
				jen.Op("&").ID("x").Dot("CreatedOn"),
				jen.Op("&").ID("x").Dot("UpdatedOn"),
				jen.Op("&").ID("x").Dot("ArchivedOn"),
				jen.Op("&").ID("x").Dot("BelongsTo"),
			), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.If(jen.ID("events").Op(":=").Qual("strings", "Split").Call(jen.ID("eventsStr"), jen.ID("eventsSeparator")), jen.ID("len").Call(jen.ID("events")).Op(">=").Lit(1).Op("&&").ID("events").Index(jen.Lit(0)).Op("!=").Lit("")).Block(
				jen.ID("x").Dot("Events").Op("=").ID("events"),
			),
			jen.If(jen.ID("dataTypes").Op(":=").Qual("strings", "Split").Call(jen.ID("dataTypesStr"), jen.ID("typesSeparator")), jen.ID("len").Call(jen.ID("dataTypes")).Op(">=").Lit(1).Op("&&").ID("dataTypes").Index(jen.Lit(0)).Op("!=").Lit("")).Block(
				jen.ID("x").Dot("DataTypes").Op("=").ID("dataTypes"),
			),
			jen.If(jen.ID("topics").Op(":=").Qual("strings", "Split").Call(jen.ID("topicsStr"), jen.ID("topicsSeparator")), jen.ID("len").Call(jen.ID("topics")).Op(">=").Lit(1).Op("&&").ID("topics").Index(jen.Lit(0)).Op("!=").Lit("")).Block(
				jen.ID("x").Dot("Topics").Op("=").ID("topics"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks"),
		jen.Line(),
		jen.Func().ID("scanWebhooks").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook"), jen.ID("error")).Block(
			jen.Var().ID("list").Index().Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook"),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("scanWebhook").Call(jen.ID("rows")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.ID("list").Op("=").ID("append").Call(jen.ID("list"), jen.Op("*").ID("webhook")),
			),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Err").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetWebhookQuery").Params(jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("webhooksTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("webhooksTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("webhookID"),
				jen.Lit("belongs_to").Op(":").ID("userID"),
			)).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("GetWebhook fetches a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetWebhookQuery").Call(jen.ID("webhookID"), jen.ID("userID")),
			jen.ID("row").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Line(),
			jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("scanWebhook").Call(jen.ID("row")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("querying for webhook"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("webhook"), jen.ID("nil")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("buildGetWebhookCountQuery returns a SQL query (and arguments) that returns a list of webhooks"),
		jen.Line(),
		jen.Comment("meeting a given filter's criteria and belonging to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetWebhookCountQuery").Params(jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("CountQuery")).
				Dotln("From").Call(jen.ID("webhooksTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("belongs_to").Op(":").ID("userID"),
				jen.Lit("archived_on").Op(":").ID("nil"),
			)),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("GetWebhookCount will fetch the count of webhooks from the database that meet a particular filter,"),
		jen.Line(),
		jen.Comment("and belong to a particular user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetWebhookCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetWebhookCountQuery").Call(jen.ID("filter"), jen.ID("userID")),
			jen.ID("err").Op("=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllWebhooksCountQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllWebhooksCountQuery").ID("string"),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetAllWebhooksCountQuery").Params().Params(jen.ID("string")).Block(
			jen.ID("getAllWebhooksCountQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("getAllWebhooksCountQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("CountQuery")).
					Dotln("From").Call(jen.ID("webhooksTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			)),
			jen.Line(),
			jen.Return().ID("getAllWebhooksCountQuery"),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("GetAllWebhooksCount will fetch the count of every active webhook in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.ID("err").Op("=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID(dbfl).Dot("buildGetAllWebhooksCountQuery").Call()).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllWebhooksQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllWebhooksQuery").ID("string"),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetAllWebhooksQuery").Params().Params(jen.ID("string")).Block(
			jen.ID("getAllWebhooksQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("getAllWebhooksQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("webhooksTableColumns").Op("...")).
					Dotln("From").Call(jen.ID("webhooksTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			)),
			jen.Line(),
			jen.Return().ID("getAllWebhooksQuery"),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("GetAllWebhooks fetches a list of all webhooks from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "WebhookList"), jen.ID("error")).Block(
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID(dbfl).Dot("buildGetAllWebhooksQuery").Call()),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for webhooks: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("scanWebhooks").Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetAllWebhooksCount").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching webhook count: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "WebhookList").Valuesln(
				jen.ID("Pagination").Op(":").Qual(filepath.Join(pkgRoot, "models/v1"), "Pagination").Valuesln(
					jen.ID("Page").Op(":").Lit(1),
					jen.ID("TotalCount").Op(":").ID("count"),
				),
				jen.ID("Webhooks").Op(":").ID("list"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("err")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("GetAllWebhooksForUser fetches a list of all webhooks from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetAllWebhooksForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetWebhooksQuery").Call(jen.ID("nil"), jen.ID("userID")),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for webhooks: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("scanWebhooks").Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("buildGetWebhooksQuery returns a SQL query (and arguments) that would return a"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetWebhooksQuery").Params(jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("webhooksTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("webhooksTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("belongs_to").Op(":").ID("userID"),
				jen.Lit("archived_on").Op(":").ID("nil")),
			),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "WebhookList"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetWebhooksQuery").Call(jen.ID("filter"), jen.ID("userID")),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("scanWebhooks").Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetWebhookCount").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching count: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "WebhookList").Valuesln(
				jen.ID("Pagination").Op(":").Qual(filepath.Join(pkgRoot, "models/v1"), "Pagination").Valuesln(
					jen.ID("Page").Op(":").ID("filter").Dot("Page"),
					jen.ID("TotalCount").Op(":").ID("count"),
					jen.ID("Limit").Op(":").ID("filter").Dot("Limit"),
				),
				jen.ID("Webhooks").Op(":").ID("list"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("err")),
		),
		jen.Line(),
	)
	////////////

	buildWebhookCreationQueryQuery := func() jen.Code {
		cols := []jen.Code{
			jen.Lit("name"),
			jen.Lit("content_type"),
			jen.Lit("url"),
			jen.Lit("method"),
			jen.Lit("events"),
			jen.Lit("data_types"),
			jen.Lit("topics"),
			jen.Lit("belongs_to"),
		}
		vals := []jen.Code{
			jen.ID("x").Dot("Name"),
			jen.ID("x").Dot("ContentType"),
			jen.ID("x").Dot("URL"),
			jen.ID("x").Dot("Method"),
			jen.Qual("strings", "Join").Call(jen.ID("x").Dot("Events"), jen.ID("eventsSeparator")),
			jen.Qual("strings", "Join").Call(jen.ID("x").Dot("DataTypes"), jen.ID("typesSeparator")),
			jen.Qual("strings", "Join").Call(jen.ID("x").Dot("Topics"), jen.ID("topicsSeparator")),
			jen.ID("x").Dot("BelongsTo"),
		}

		if isMariaDB {
			cols = append(cols, jen.Lit("created_on"))
			vals = append(vals, jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery")))
		}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Insert").Call(jen.ID("webhooksTableName")).
			Dotln("Columns").Callln(cols...).
			Dotln("Values").Callln(vals...)

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	ret.Add(
		jen.Comment("buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildWebhookCreationQuery").Params(jen.ID("x").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildWebhookCreationQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	if isSqlite || isMariaDB {
		ret.Add(
			jen.Comment("buildWebhookCreationTimeQuery returns a SQL query (and arguments) that fetches the DB creation time for a given row"),
			jen.Line(),
			jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildWebhookCreationTimeQuery").Params(jen.ID("webhookID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.Lit("created_on")).
					Dotln("From").Call(jen.ID("webhooksTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("webhookID"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
				jen.Line(),
				jen.Return(jen.List(jen.ID("query"), jen.ID("args"))),
			),
			jen.Line(),
		)
	}

	////////////

	buildCreateWebhookQuery := func() []jen.Code {
		out := []jen.Code{
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook").Valuesln(
				jen.ID("Name").Op(":").ID("input").Dot("Name"),
				jen.ID("ContentType").Op(":").ID("input").Dot("ContentType"),
				jen.ID("URL").Op(":").ID("input").Dot("URL"),
				jen.ID("Method").Op(":").ID("input").Dot("Method"),
				jen.ID("Events").Op(":").ID("input").Dot("Events"),
				jen.ID("DataTypes").Op(":").ID("input").Dot("DataTypes"),
				jen.ID("Topics").Op(":").ID("input").Dot("Topics"),
				jen.ID("BelongsTo").Op(":").ID("input").Dot("BelongsTo"),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildWebhookCreationQuery").Call(jen.ID("x")),
		}

		if isPostgres {
			out = append(out,
				jen.If(jen.ID("err").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("CreatedOn")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing webhook creation query: %w"), jen.ID("err"))),
				),
			)
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return(jen.List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing webhook creation query: %w"), jen.ID("err")))),
				),
				jen.Line(),
				jen.If(jen.List(jen.ID("id"), jen.ID("idErr")).Op(":=").ID("res").Dot("LastInsertId").Call().Op(";").ID("idErr").Op("==").ID("nil")).Block(
					jen.ID("x").Dot("ID").Op("=").ID("uint64").Call(jen.ID("id")),
					jen.Line(),
					jen.List(jen.ID("query"), jen.ID("args")).Op("=").ID(dbfl).Dot("buildWebhookCreationTimeQuery").Call(jen.ID("x").Dot("ID")),
					jen.ID(dbfl).Dot("logCreationTimeRetrievalError").Call(jen.ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("CreatedOn"))),
				),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		)

		return out
	}

	ret.Add(
		jen.Comment("CreateWebhook creates a webhook in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("CreateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "WebhookCreationInput")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook"), jen.ID("error")).Block(
			buildCreateWebhookQuery()...,
		),
		jen.Line(),
	)
	////////////

	buildUpdateWebhookQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("webhooksTableName")).
			Dotln("Set").Call(jen.Lit("name"), jen.ID("input").Dot("Name")).
			Dotln("Set").Call(jen.Lit("content_type"), jen.ID("input").Dot("ContentType")).
			Dotln("Set").Call(jen.Lit("url"), jen.ID("input").Dot("URL")).
			Dotln("Set").Call(jen.Lit("method"), jen.ID("input").Dot("Method")).
			Dotln("Set").Call(jen.Lit("events"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Events"), jen.ID("topicsSeparator"))).
			Dotln("Set").Call(jen.Lit("data_types"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("DataTypes"), jen.ID("typesSeparator"))).
			Dotln("Set").Call(jen.Lit("topics"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Topics"), jen.ID("topicsSeparator"))).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").Op(":").ID("input").Dot("ID"),
			jen.Lit("belongs_to").Op(":").ID("input").Dot("BelongsTo")),
		)

		if isPostgres {
			q.Dot("Suffix").Call(jen.Lit("RETURNING updated_on"))
		}

		q.Dotln("ToSql").Call()

		return q
	}

	ret.Add(
		jen.Comment("buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildUpdateWebhookQuery").Params(jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildUpdateWebhookQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)
	////////////

	buildUpdateWebhookBody := func() []jen.Code {
		if isPostgres {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateWebhookQuery").Call(jen.ID("input")),
				jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("input").Dot("UpdatedOn")),
			}
		} else if isSqlite || isMariaDB {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateWebhookQuery").Call(jen.ID("input")),
				jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
				jen.Return(jen.ID("err")),
			}
		}
		return nil
	}

	ret.Add(
		jen.Comment("UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("UpdateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "Webhook")).Params(jen.ID("error")).Block(
			buildUpdateWebhookBody()...,
		),
		jen.Line(),
	)
	////////////

	buildArchiveWebhookQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("webhooksTableName")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").Op(":").ID("webhookID"),
			jen.Lit("belongs_to").Op(":").ID("userID"),
			jen.Lit("archived_on").Op(":").ID("nil"),
		))

		if isPostgres {
			q.Dot("Suffix").Call(jen.Lit("RETURNING archived_on"))
		}

		q.Dotln("ToSql").Call()

		return q
	}

	ret.Add(
		jen.Comment("buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildArchiveWebhookQuery").Params(jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildArchiveWebhookQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)
	////////////

	ret.Add(
		jen.Comment("ArchiveWebhook archives a webhook from the database by its ID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("ArchiveWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildArchiveWebhookQuery").Call(jen.ID("webhookID"), jen.ID("userID")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
