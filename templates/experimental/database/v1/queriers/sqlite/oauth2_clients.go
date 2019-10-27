package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("sqlite")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("scopesSeparator").Op("=").Lit(`,`),
			jen.ID("oauth2ClientsTableName").Op("=").Lit("oauth2_clients"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("oauth2ClientsTableColumns").Op("=").Index().ID("string").Valuesln(
			jen.Lit("id"),
			jen.Lit("name"),
			jen.Lit("client_id"),
			jen.Lit("scopes"),
			jen.Lit("redirect_uri"),
			jen.Lit("client_secret"),
			jen.Lit("created_on"),
			jen.Lit("updated_on"),
			jen.Lit("archived_on"),
			jen.Lit("belongs_to"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its ressults into an OAuth2Client struct"),
		jen.Line(),
		jen.Func().ID("scanOAuth2Client").Params(jen.ID("scan").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Scanner")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2Client"), jen.ID("error")).Block(
			jen.Var().Defs(
				jen.ID("x").Op("=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2Client").Values(),
				jen.ID("scopes").ID("string"),
			),
			jen.If(jen.ID("err").Op(":=").ID("scan").Dot("Scan").Call(jen.Op("&").ID("x").Dot("ID"),
				jen.Op("&").ID("x").Dot("Name"),
				jen.Op("&").ID("x").Dot("ClientID"),
				jen.Op("&").ID("scopes"), jen.Op("&").ID("x").Dot("RedirectURI"),
				jen.Op("&").ID("x").Dot("ClientSecret"),
				jen.Op("&").ID("x").Dot("CreatedOn"),
				jen.Op("&").ID("x").Dot("UpdatedOn"),
				jen.Op("&").ID("x").Dot("ArchivedOn"),
				jen.Op("&").ID("x").Dot("BelongsTo"),
			), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.If(jen.ID("scopes").Op(":=").Qual("strings", "Split").Call(jen.ID("scopes"), jen.ID("scopesSeparator")), jen.ID("len").Call(jen.ID("scopes")).Op(">=").Lit(1).Op("&&").ID("scopes").Index(jen.Lit(0)).Op("!=").Lit("")).Block(
				jen.ID("x").Dot("Scopes").Op("=").ID("scopes"),
			),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("scanOAuth2Clients").Params(jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2Client"), jen.ID("error")).Block(
			jen.Var().ID("list").Index().Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2Client"),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("scanOAuth2Client").Call(jen.ID("rows")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.ID("list").Op("=").ID("append").Call(jen.ID("list"), jen.ID("client")),
			),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Err").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.ID("s").Dot("logQueryBuildingError").Call(jen.ID("rows").Dot("Close").Call()),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetOAuth2ClientByClientIDQuery").Params(jen.ID("clientID").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot("sqlBuilder").Dot("Select").Call(jen.ID("oauth2ClientsTableColumns").Op("...")).Dot("From").Call(jen.ID("oauth2ClientsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Lit("client_id").Op(":").ID("clientID"), jen.Lit("archived_on").Op(":").ID("nil"))).Dot("ToSql").Call(),
			jen.ID("s").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientByClientID gets an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetOAuth2ClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2Client"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetOAuth2ClientByClientIDQuery",
			).Call(jen.ID("clientID")),
			jen.ID("row").Op(":=").ID("s").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("scanOAuth2Client").Call(jen.ID("row")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientsQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientsQuery").ID("string"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetAllOAuth2ClientsQuery builds a SQL query"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetAllOAuth2ClientsQuery").Params().Params(jen.ID("query").ID("string")).Block(
			jen.ID("getAllOAuth2ClientsQueryBuilder").Dot(
				"Do",
			).Call(jen.Func().Params().Block(

				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("getAllOAuth2ClientsQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID("s").Dot(
					"sqlBuilder",
				).Dot(
					"Select",
				).Call(jen.ID("oauth2ClientsTableColumns").Op("...")).Dot(
					"From",
				).Call(jen.ID("oauth2ClientsTableName")).Dot(
					"Where",
				).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.Lit("archived_on").Op(":").ID("nil"))).Dot(
					"ToSql",
				).Call(),
				jen.ID("s").Dot(
					"logQueryBuildingError",
				).Call(jen.ID("err")),
			)),
			jen.Return().ID("getAllOAuth2ClientsQuery"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetAllOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"QueryContext",
			).Call(jen.ID("ctx"), jen.ID("s").Dot(
				"buildGetAllOAuth2ClientsQuery",
			).Call()),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.ID("err"))),
			),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("s").Dot(
				"scanOAuth2Clients",
			).Call(jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.ID("err"))),
			),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetAllOAuth2ClientsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetOAuth2ClientsQuery",
			).Call(jen.ID("nil"), jen.ID("userID")),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"QueryContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.ID("err"))),
			),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("s").Dot(
				"scanOAuth2Clients",
			).Call(jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.ID("err"))),
			),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("oauth2ClientsTableColumns").Op("...")).Dot(
				"From",
			).Call(jen.ID("oauth2ClientsTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("clientID"), jen.Lit("belongs_to").Op(":").ID("userID"), jen.Lit("archived_on").Op(":").ID("nil"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Client retrieves an OAuth2 client from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetOAuth2ClientQuery",
			).Call(jen.ID("clientID"), jen.ID("userID")),
			jen.ID("row").Op(":=").ID("s").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("scanOAuth2Client").Call(jen.ID("row")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for oauth2 client: %w"), jen.ID("err"))),
			),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetOAuth2ClientCountQuery returns a SQL query (and arguments) that fetches a list of OAuth2 clients that meet certain filter"),
		jen.Line(),
		jen.Comment("restrictions (if relevant) and belong to a given user"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetOAuth2ClientCountQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("CountQuery")).Dot(
				"From",
			).Call(jen.ID("oauth2ClientsTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("belongs_to").Op(":").ID("userID"), jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot(
					"ApplyToQueryBuilder",
				).Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientCount will get the count of OAuth2 clients that match the given filter and belong to the user"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetOAuth2ClientCountQuery",
			).Call(jen.ID("filter"), jen.ID("userID")),
			jen.ID("err").Op("=").ID("s").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
				"Scan",
			).Call(jen.Op("&").ID("count")),
			jen.Return(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientCountQuery").ID("string"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetAllOAuth2ClientCountQuery returns a SQL query for the number of OAuth2 clients"),
		jen.Line(),
		jen.Comment("in the database, regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetAllOAuth2ClientCountQuery").Params().Params(jen.ID("string")).Block(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Dot(
				"Do",
			).Call(jen.Func().Params().Block(

				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("getAllOAuth2ClientCountQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID("s").Dot(
					"sqlBuilder",
				).Dot(
					"Select",
				).Call(jen.ID("CountQuery")).Dot(
					"From",
				).Call(jen.ID("oauth2ClientsTableName")).Dot(
					"Where",
				).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.Lit("archived_on").Op(":").ID("nil"))).Dot(
					"ToSql",
				).Call(),
				jen.ID("s").Dot(
					"logQueryBuildingError",
				).Call(jen.ID("err")),
			)),
			jen.Return().ID("getAllOAuth2ClientCountQuery"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetAllOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(

			jen.Var().ID("count").ID("uint64"),
			jen.ID("err").Op(":=").ID("s").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("s").Dot(
				"buildGetAllOAuth2ClientCountQuery",
			).Call()).Dot(
				"Scan",
			).Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetOAuth2ClientsQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that"),
		jen.Line(),
		jen.Comment("meet the given filter's criteria (if relevant) and belong to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetOAuth2ClientsQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("oauth2ClientsTableColumns").Op("...")).Dot(
				"From",
			).Call(jen.ID("oauth2ClientsTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("belongs_to").Op(":").ID("userID"), jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot(
					"ApplyToQueryBuilder",
				).Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2ClientList",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetOAuth2ClientsQuery",
			).Call(jen.ID("filter"), jen.ID("userID")),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"QueryContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for oauth2 clients: %w"), jen.ID("err"))),
			),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("s").Dot(
				"scanOAuth2Clients",
			).Call(jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
			),
			jen.ID("ll").Op(":=").ID("len").Call(jen.ID("list")),

			jen.Var().ID("clients").Op("=").ID("make").Call(jen.Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"OAuth2Client",
			),
				jen.ID("ll")),
			jen.For(jen.List(jen.ID("i"), jen.ID("t")).Op(":=").Range().ID("list")).Block(
				jen.ID("clients").Index(jen.ID("i")).Op("=").Op("*").ID("t"),
			),
			jen.List(jen.ID("totalCount"), jen.ID("err")).Op(":=").ID("s").Dot(
				"GetOAuth2ClientCount",
			).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching oauth2 client count: %w"), jen.ID("err"))),
			),
			jen.ID("ocl").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"OAuth2ClientList",
			).Valuesln(
				jen.ID("Pagination").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
					"Pagination",
				).Valuesln(
					jen.ID("Page").Op(":").ID("filter").Dot(
						"Page",
					),
					jen.ID("Limit").Op(":").ID("filter").Dot(
						"Limit",
					),
					jen.ID("TotalCount").Op(":").ID("totalCount")), jen.ID("Clients").Op(":").ID("clients")),
			jen.Return().List(jen.ID("ocl"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildCreateOAuth2ClientQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2Client",
		)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Insert",
			).Call(jen.ID("oauth2ClientsTableName")).Dot(
				"Columns",
			).Call(jen.Lit("name"), jen.Lit("client_id"), jen.Lit("client_secret"), jen.Lit("scopes"), jen.Lit("redirect_uri"), jen.Lit("belongs_to")).Dot(
				"Values",
			).Call(jen.ID("input").Dot("Name"),
				jen.ID("input").Dot(
					"ClientID",
				),
				jen.ID("input").Dot(
					"ClientSecret",
				),
				jen.Qual("strings", "Join").Call(jen.ID("input").Dot(
					"Scopes",
				),
					jen.ID("scopesSeparator")), jen.ID("input").Dot(
					"RedirectURI",
				),
				jen.ID("input").Dot("BelongsTo")).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildOAuth2ClientCreationTimeQuery").Params(jen.ID("clientID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.Lit("created_on")).Dot(
				"From",
			).Call(jen.ID("oauth2ClientsTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("clientID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateOAuth2Client creates an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("CreateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2ClientCreationInput",
		)).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"OAuth2Client",
			).Valuesln(
				jen.ID("Name").Op(":").ID("input").Dot("Name"),
				jen.ID("ClientID").Op(":").ID("input").Dot(
					"ClientID",
				),
				jen.ID("ClientSecret").Op(":").ID("input").Dot(
					"ClientSecret",
				),
				jen.ID("RedirectURI").Op(":").ID("input").Dot(
					"RedirectURI",
				),
				jen.ID("Scopes").Op(":").ID("input").Dot(
					"Scopes",
				),
				jen.ID("BelongsTo").Op(":").ID("input").Dot("BelongsTo")),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildCreateOAuth2ClientQuery",
			).Call(jen.ID("x")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing client creation query: %w"), jen.ID("err"))),
			),
			jen.If(jen.List(jen.ID("id"), jen.ID("idErr")).Op(":=").ID("res").Dot(
				"LastInsertId",
			).Call(), jen.ID("idErr").Op("==").ID("nil")).Block(
				jen.ID("x").Dot("ID").Op("=").ID("uint64").Call(jen.ID("id")),
				jen.List(jen.ID("query"), jen.ID("args")).Op("=").ID("s").Dot(
					"buildOAuth2ClientCreationTimeQuery",
				).Call(jen.ID("x").Dot("ID")),
				jen.ID("s").Dot(
					"logCreationTimeRetrievalError",
				).Call(jen.ID("s").Dot("db").Dot(
					"QueryRowContext",
				).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
					"Scan",
				).Call(jen.Op("&").ID("x").Dot("CreatedOn"))),
			),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildUpdateOAuth2ClientQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2Client",
		)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Update",
			).Call(jen.ID("oauth2ClientsTableName")).Dot("Set").Call(jen.Lit("client_id"), jen.ID("input").Dot(
				"ClientID",
			)).Dot("Set").Call(jen.Lit("client_secret"), jen.ID("input").Dot(
				"ClientSecret",
			)).Dot("Set").Call(jen.Lit("scopes"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot(
				"Scopes",
			),
				jen.ID("scopesSeparator"))).Dot("Set").Call(jen.Lit("redirect_uri"), jen.ID("input").Dot(
				"RedirectURI",
			)).Dot("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("input").Dot("ID"),
				jen.Lit("belongs_to").Op(":").ID("input").Dot("BelongsTo"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateOAuth2Client updates a OAuth2 client."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the input's ID field to be valid and non-zero."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("UpdateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"OAuth2Client",
		)).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildUpdateOAuth2ClientQuery",
			).Call(jen.ID("input")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildArchiveOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Update",
			).Call(jen.ID("oauth2ClientsTableName")).Dot("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).Dot("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("clientID"), jen.Lit("belongs_to").Op(":").ID("userID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("ArchiveOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildArchiveOAuth2ClientQuery",
			).Call(jen.ID("clientID"), jen.ID("userID")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
