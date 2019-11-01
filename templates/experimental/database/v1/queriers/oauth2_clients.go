package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

func oauth2ClientsDotGo(pkgRoot string, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(ret)
	sn := vendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := vendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Const().Defs(
			jen.ID("scopesSeparator").Op("=").Lit(`,`),
			jen.ID("oauth2ClientsTableName").Op("=").Lit("oauth2_clients"),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("oauth2ClientsTableColumns").Op("=").Index().ID("string").Valuesln(
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
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its ressults into an OAuth2Client struct"),
		jen.Line(),
		jen.Func().ID("scanOAuth2Client").Params(jen.ID("scan").Qual(filepath.Join(pkgRoot, "database/v1"), "Scanner")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("error")).Block(
			jen.Var().Defs(
				jen.ID("x").Op("=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client").Values(),
				jen.ID("scopes").ID("string"),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("scan").Dot("Scan").Callln(
				jen.Op("&").ID("x").Dot("ID"),
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
			jen.Line(),
			jen.If(jen.ID("scopes").Op(":=").Qual("strings", "Split").Call(jen.ID("scopes"), jen.ID("scopesSeparator")), jen.ID("len").Call(jen.ID("scopes")).Op(">=").Lit(1).Op("&&").ID("scopes").Index(jen.Lit(0)).Op("!=").Lit("")).Block(
				jen.ID("x").Dot(
					"Scopes",
				).Op("=").ID("scopes"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("scanOAuth2Clients").Params(jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("error")).Block(
			jen.Var().ID("list").Index().Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
			jen.Line(),
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
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID(dbfl).Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetOAuth2ClientByClientIDQuery").Params(jen.ID("clientID").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.Comment("This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't"),
			jen.Comment("care about ownership. It does still care about archived status"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("client_id").Op(":").ID("clientID"),
				jen.Lit("archived_on").Op(":").ID("nil"),
			)).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetOAuth2ClientByClientID gets an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetOAuth2ClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientByClientIDQuery").Call(jen.ID("clientID")),
			jen.ID("row").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("scanOAuth2Client").Call(jen.ID("row")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientsQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientsQuery").ID("string"),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetAllOAuth2ClientsQuery builds a SQL query"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetAllOAuth2ClientsQuery").Params().Params(jen.ID("query").ID("string")).Block(
			jen.ID("getAllOAuth2ClientsQueryBuilder").Dot(
				"Do",
			).Call(jen.Func().Params().Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("getAllOAuth2ClientsQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Op("...")).
					Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			)),
			jen.Line(),
			jen.Return().ID("getAllOAuth2ClientsQuery"),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetAllOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").Qual(filepath.Join(pkgRoot, "models/v1"),
			"OAuth2Client",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID(dbfl).Dot("buildGetAllOAuth2ClientsQuery").Call()),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetAllOAuth2ClientsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientsQuery").Call(jen.ID("nil"), jen.ID("userID")),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("clientID"),
				jen.Lit("belongs_to").Op(":").ID("userID"),
				jen.Lit("archived_on").Op(":").ID("nil"),
			)).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetOAuth2Client retrieves an OAuth2 client from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientQuery").Call(jen.ID("clientID"), jen.ID("userID")),
			jen.ID("row").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Line(),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("scanOAuth2Client").Call(jen.ID("row")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for oauth2 client: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetOAuth2ClientCountQuery returns a SQL query (and arguments) that fetches a list of OAuth2 clients that meet certain filter"),
		jen.Line(),
		jen.Comment("restrictions (if relevant) and belong to a given user"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetOAuth2ClientCountQuery").Params(jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("CountQuery")).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("belongs_to").Op(":").ID("userID"),
				jen.Lit("archived_on").Op(":").ID("nil"),
			)),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot(
					"ApplyToQueryBuilder",
				).Call(jen.ID("builder")),
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
		jen.Comment("GetOAuth2ClientCount will get the count of OAuth2 clients that match the given filter and belong to the user"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientCountQuery").Call(jen.ID("filter"), jen.ID("userID")),
			jen.ID("err").Op("=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return(),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientCountQuery").ID("string"),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetAllOAuth2ClientCountQuery returns a SQL query for the number of OAuth2 clients"),
		jen.Line(),
		jen.Comment("in the database, regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetAllOAuth2ClientCountQuery").Params().Params(jen.ID("string")).Block(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("getAllOAuth2ClientCountQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("CountQuery")).
					Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			)),
			jen.Line(),
			jen.Return().ID("getAllOAuth2ClientCountQuery"),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetAllOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.Var().ID("count").ID("uint64"),
			jen.ID("err").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID(dbfl).Dot("buildGetAllOAuth2ClientCountQuery").Call()).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetOAuth2ClientsQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that"),
		jen.Line(),
		jen.Comment("meet the given filter's criteria (if relevant) and belong to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetOAuth2ClientsQuery").Params(jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
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
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientList"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetOAuth2ClientsQuery").Call(jen.ID("filter"), jen.ID("userID")),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Line(),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for oauth2 clients: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Comment("de-pointer-ize clients"),
			jen.ID("ll").Op(":=").ID("len").Call(jen.ID("list")),
			jen.Var().ID("clients").Op("=").ID("make").Call(jen.Index().Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("ll")),
			jen.For(jen.List(jen.ID("i"), jen.ID("t")).Op(":=").Range().ID("list")).Block(
				jen.ID("clients").Index(jen.ID("i")).Op("=").Op("*").ID("t"),
			),
			jen.Line(),
			jen.List(jen.ID("totalCount"), jen.ID("err")).Op(":=").ID(dbfl).Dot("GetOAuth2ClientCount").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching oauth2 client count: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("ocl").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientList").Valuesln(
				jen.ID("Pagination").Op(":").Qual(filepath.Join(pkgRoot, "models/v1"), "Pagination").Valuesln(
					jen.ID("Page").Op(":").ID("filter").Dot("Page"),
					jen.ID("Limit").Op(":").ID("filter").Dot("Limit"),
					jen.ID("TotalCount").Op(":").ID("totalCount"),
				),
				jen.ID("Clients").Op(":").ID("clients"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("ocl"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	buildCreateQueryCreation := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Insert").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Columns").Callln(
			jen.Lit("name"),
			jen.Lit("client_id"),
			jen.Lit("client_secret"),
			jen.Lit("scopes"),
			jen.Lit("redirect_uri"),
			jen.Lit("belongs_to"),
		).
			Dotln("Values").Callln(
			jen.ID("input").Dot("Name"),
			jen.ID("input").Dot("ClientID"),
			jen.ID("input").Dot("ClientSecret"),
			jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Scopes"), jen.ID("scopesSeparator")),
			jen.ID("input").Dot("RedirectURI"),
			jen.ID("input").Dot("BelongsTo"),
		)
		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	ret.Add(
		jen.Comment("buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildCreateOAuth2ClientQuery").Params(jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildCreateQueryCreation(),
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
			jen.Comment("buildOAuth2ClientCreationTimeQuery takes an item and returns a creation query for that item and the relevant arguments."),
			jen.Line(),
			jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildOAuth2ClientCreationTimeQuery").Params(jen.ID("clientID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
				jen.Var().ID("err").ID("error"),
				jen.Line(),
				jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.Lit("created_on")).
					Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("clientID"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
				jen.Line(),
				jen.Return().List(jen.ID("query"), jen.ID("args")),
			),
			jen.Line(),
		)
	}

	////////////

	buildCreateOauth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client").Valuesln(
				jen.ID("Name").Op(":").ID("input").Dot("Name"),
				jen.ID("ClientID").Op(":").ID("input").Dot("ClientID"),
				jen.ID("ClientSecret").Op(":").ID("input").Dot("ClientSecret"),
				jen.ID("RedirectURI").Op(":").ID("input").Dot("RedirectURI"),
				jen.ID("Scopes").Op(":").ID("input").Dot("Scopes"),
				jen.ID("BelongsTo").Op(":").ID("input").Dot("BelongsTo")),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildCreateOAuth2ClientQuery").Call(jen.ID("x")),
			jen.Line(),
		}

		if isPostgres {
			out = append(out,
				jen.ID("err").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("CreatedOn")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing client creation query: %w"), jen.ID("err"))),
				),
			)
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing client creation query: %w"), jen.ID("err"))),
				),
				jen.Line(),
				jen.Comment("fetch the last inserted ID"),
				jen.If(jen.List(jen.ID("id"), jen.ID("idErr")).Op(":=").ID("res").Dot("LastInsertId").Call()).Op(";").ID("idErr").Op("==").ID("nil").Block(
					jen.ID("x").Dot("ID").Op("=").ID("uint64").Call(jen.ID("id")),
					jen.Line(),
					jen.List(jen.ID("query"), jen.ID("args")).Op("=").ID(dbfl).Dot("buildOAuth2ClientCreationTimeQuery").Call(jen.ID("x").Dot("ID")),
					jen.ID(dbfl).Dot("logCreationTimeRetrievalError").Call(jen.ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("CreatedOn"))),
				),
			)
		}
		out = append(out, jen.Line(), jen.Return().List(jen.ID("x"), jen.ID("nil")))

		return out
	}

	ret.Add(
		jen.Comment("CreateOAuth2Client creates an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("CreateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientCreationInput")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("error")).Block(
			buildCreateOauth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildBuildCreateOauth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{jen.Var().ID("err").ID("error")}
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Set").Call(jen.Lit("client_id"), jen.ID("input").Dot("ClientID")).
			Dotln("Set").Call(jen.Lit("client_secret"), jen.ID("input").Dot("ClientSecret")).
			Dotln("Set").Call(jen.Lit("scopes"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Scopes"), jen.ID("scopesSeparator"))).
			Dotln("Set").Call(jen.Lit("redirect_uri"), jen.ID("input").Dot("RedirectURI")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").Op(":").ID("input").Dot("ID"),
			jen.Lit("belongs_to").Op(":").ID("input").Dot("BelongsTo"),
		))

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING updated_on"))
		}
		q.Dotln("ToSql").Call()

		out = append(out, q,
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		)

		return out
	}

	ret.Add(
		jen.Comment("buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildUpdateOAuth2ClientQuery").Params(jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			buildBuildCreateOauth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildUpdateOAuth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateOAuth2ClientQuery").Call(jen.ID("input")),
		}

		if isPostgres {
			out = append(out, jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("input").Dot("UpdatedOn")))
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
				jen.Return().ID("err"),
			)
		}

		return out
	}

	ret.Add(
		jen.Comment("UpdateOAuth2Client updates a OAuth2 client."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the input's ID field to be valid and non-zero."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("UpdateOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client")).Params(jen.ID("error")).Block(
			buildUpdateOAuth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildBuildArchiveOAuth2ClientQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").Op(":").ID("clientID"),
			jen.Lit("belongs_to").Op(":").ID("userID"),
		))

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING archived_on"))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	ret.Add(
		jen.Comment("buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildArchiveOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildBuildArchiveOAuth2ClientQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("ArchiveOAuth2Client").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildArchiveOAuth2ClientQuery").Call(jen.ID("clientID"), jen.ID("userID")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
