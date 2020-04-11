package queriers

import (
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(proj, ret)
	sn := vendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := vendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Const().Defs(
			jen.ID("scopesSeparator").Equals().Lit(`,`),
			jen.ID("oauth2ClientsTableName").Equals().Lit("oauth2_clients"),
			jen.ID("oauth2ClientsTableOwnershipColumn").Equals().Lit("belongs_to_user"),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("oauth2ClientsTableColumns").Equals().Index().String().Valuesln(
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.name"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.client_id"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.scopes"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.redirect_uri"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.client_secret"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.created_on"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.updated_on"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("oauth2ClientsTableName")),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableOwnershipColumn")),
			),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct"),
		jen.Line(),
		jen.Func().ID("scanOAuth2Client").Params(
			jen.ID("scan").Qual(proj.DatabaseV1Package(), "Scanner"),
			jen.ID("includeCount").Bool(),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Uint64(),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("x").Equals().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.ID("scopes").String(),
				jen.ID("count").Uint64(),
			),
			jen.Line(),
			jen.ID("targetVars").Assign().Index().Interface().Valuesln(
				jen.AddressOf().ID("x").Dot("ID"),
				jen.AddressOf().ID("x").Dot("Name"),
				jen.AddressOf().ID("x").Dot("ClientID"),
				jen.AddressOf().ID("scopes"), jen.AddressOf().ID("x").Dot("RedirectURI"),
				jen.AddressOf().ID("x").Dot("ClientSecret"),
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
			jen.If(
				jen.ID("scopes").Assign().Qual("strings", "Split").Call(jen.ID("scopes"), jen.ID("scopesSeparator")),
				jen.ID("len").Call(jen.ID("scopes")).Op(">=").One().And().ID("scopes").Index(jen.Zero()).DoesNotEqual().EmptyString(),
			).Block(
				jen.ID("x").Dot("Scopes").Equals().ID("scopes"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("count"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("scanOAuth2Clients").Params(
			jen.ID("rows").PointerTo().Qual("database/sql", "Rows"),
		).Params(
			jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Uint64(),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("list").Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
				jen.ID("count").Uint64(),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("client"), jen.ID("c"), jen.Err()).Assign().ID("scanOAuth2Client").Call(
					jen.ID("rows"),
					jen.True(),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Zero(), jen.Err()),
				),
				jen.Line(),
				jen.If(jen.ID("count").IsEqualTo().Zero()).Block(
					jen.ID("count").Equals().ID("c"),
				),
				jen.Line(),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.ID("client")),
			),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(dbfl).Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("count"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetOAuth2ClientByClientIDQuery").Params(jen.ID("clientID").String()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.Comment("This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't"),
			jen.Comment("care about ownership. It does still care about archived status"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Spread()).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.client_id"), jen.ID("oauth2ClientsTableName")).MapAssign().ID("clientID"),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("oauth2ClientsTableName")).MapAssign().ID("nil"),
			)).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetOAuth2ClientByClientID gets an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetOAuth2ClientByClientID").Params(
			utils.CtxParam(),
			jen.ID("clientID").String(),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error(),
		).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientByClientIDQuery").Call(jen.ID("clientID")),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Underscore(), jen.Err()).Assign().ID("scanOAuth2Client").Call(jen.ID("row"), jen.False()),
			jen.Return().List(jen.ID("client"), jen.Err()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientsQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientsQuery").String(),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetAllOAuth2ClientsQuery builds a SQL query"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllOAuth2ClientsQuery").Params().Params(jen.ID("query").String()).Block(
			jen.ID("getAllOAuth2ClientsQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("getAllOAuth2ClientsQuery"), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Spread()).
					Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
					Dotln("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
						jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("oauth2ClientsTableName")).MapAssign().Nil(),
					),
				).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
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
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllOAuth2Clients").Params(utils.CtxParam()).Params(jen.Index().PointerTo().Qual(proj.ModelsV1Package(),
			"OAuth2Client",
		),
			jen.Error()).Block(
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID(dbfl).Dot("buildGetAllOAuth2ClientsQuery").Call()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllOAuth2ClientsForUser").Params(utils.CtxParam(), jen.ID("userID").Uint64()).Params(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientsQuery").Call(jen.ID("userID"), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Spread()).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.ID("oauth2ClientsTableName")).MapAssign().ID("clientID"),
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableOwnershipColumn")).MapAssign().ID("userID"),
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("oauth2ClientsTableName")).MapAssign().Nil(),
				),
			).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetOAuth2Client retrieves an OAuth2 client from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetOAuth2Client").Params(utils.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientQuery").Call(jen.ID("clientID"), jen.ID("userID")),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Underscore(), jen.Err()).Assign().ID("scanOAuth2Client").Call(jen.ID("row"), jen.False()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for oauth2 client: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientCountQuery").String(),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetAllOAuth2ClientCountQuery returns a SQL query for the number of OAuth2 clients"),
		jen.Line(),
		jen.Comment("in the database, regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllOAuth2ClientCountQuery").Params().Params(jen.String()).Block(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("getAllOAuth2ClientCountQuery"), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(jen.ID("countQuery"), jen.ID("oauth2ClientsTableName"))).
					Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
					Dotln("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
						jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("oauth2ClientsTableName")).MapAssign().ID("nil"),
					),
				).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
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
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllOAuth2ClientCount").Params(utils.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.Var().ID("count").Uint64(),
			jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID(dbfl).Dot("buildGetAllOAuth2ClientCountQuery").Call()).Dot("Scan").Call(jen.AddressOf().ID("count")),
			jen.Return().List(jen.ID("count"), jen.Err()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetOAuth2ClientsQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that"),
		jen.Line(),
		jen.Comment("meet the given filter's criteria (if relevant) and belong to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetOAuth2ClientsQuery").Params(
			jen.ID("userID").Uint64(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.Append(jen.ID("oauth2ClientsTableColumns"), jen.Qual("fmt", "Sprintf").Call(jen.ID("countQuery"), jen.ID("oauth2ClientsTableName"))).Spread()).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableOwnershipColumn")).MapAssign().ID("userID"),
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.ID("oauth2ClientsTableName")).MapAssign().Nil(),
				),
			).Dotln("GroupBy").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.ID("oauth2ClientsTableName"))),
			jen.Line(),
			jen.If(jen.ID(utils.FilterVarName).DoesNotEqual().ID("nil")).Block(
				jen.ID("builder").Equals().ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder"), jen.ID("oauth2ClientsTableName")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetOAuth2Clients").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList"),
			jen.Error(),
		).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientsQuery").Call(jen.ID("userID"), jen.ID(utils.FilterVarName)),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for oauth2 clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("count"), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("ocl").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID("filter").Dot("Page"),
					jen.ID("Limit").MapAssign().ID("filter").Dot("Limit"),
					jen.ID("TotalCount").MapAssign().ID("count"),
				),
			),
			jen.Line(),
			jen.Comment("de-pointer-ize clients"),
			jen.ID("ocl").Dot("Clients").Equals().Make(jen.Index().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Len(jen.ID("list"))),
			jen.For(jen.List(jen.ID("i"), jen.ID("t")).Assign().Range().ID("list")).Block(
				jen.ID("ocl").Dot("Clients").Index(jen.ID("i")).Equals().PointerTo().ID("t"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("ocl"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	buildCreateQueryCreation := func() jen.Code {
		cols := []jen.Code{
			jen.Lit("name"),
			jen.Lit("client_id"),
			jen.Lit("client_secret"),
			jen.Lit("scopes"),
			jen.Lit("redirect_uri"),
			jen.ID("oauth2ClientsTableOwnershipColumn"),
		}

		vals := []jen.Code{
			jen.ID("input").Dot("Name"),
			jen.ID("input").Dot("ClientID"),
			jen.ID("input").Dot("ClientSecret"),
			jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Scopes"), jen.ID("scopesSeparator")),
			jen.ID("input").Dot("RedirectURI"),
			jen.ID("input").Dot("BelongsToUser"),
		}

		if isMariaDB {
			cols = append(cols, jen.Lit("created_on"))
			vals = append(vals, jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery")))
		}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Insert").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Columns").Callln(cols...).
			Dotln("Values").Callln(vals...)
		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	ret.Add(
		jen.Comment("buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildCreateOAuth2ClientQuery").Params(jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			buildCreateQueryCreation(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	if isSqlite || isMariaDB {
		ret.Add(
			jen.Comment("buildOAuth2ClientCreationTimeQuery takes an oauth2 client ID and returns a creation query"),
			jen.Line(),
			jen.Comment("for that oauth2 client and the relevant arguments."),
			jen.Line(),
			jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildOAuth2ClientCreationTimeQuery").Params(jen.ID("clientID").Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.Lit("created_on")).
					Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.Lit("id").MapAssign().ID("clientID"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
				jen.Line(),
				jen.Return().List(jen.ID("query"), jen.ID("args")),
			),
			jen.Line(),
		)
	}

	////////////

	buildCreateOauth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
				jen.ID("Name").MapAssign().ID("input").Dot("Name"),
				jen.ID("ClientID").MapAssign().ID("input").Dot("ClientID"),
				jen.ID("ClientSecret").MapAssign().ID("input").Dot("ClientSecret"),
				jen.ID("RedirectURI").MapAssign().ID("input").Dot("RedirectURI"),
				jen.ID("Scopes").MapAssign().ID("input").Dot("Scopes"),
				jen.ID("BelongsToUser").MapAssign().ID("input").Dot("BelongsToUser")),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildCreateOAuth2ClientQuery").Call(jen.ID("x")),
			jen.Line(),
		}

		if isPostgres {
			out = append(out,
				jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("ID"), jen.AddressOf().ID("x").Dot("CreatedOn")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing client creation query: %w"), jen.Err())),
				),
			)
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.List(jen.ID("res"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing client creation query: %w"), jen.Err())),
				),
				jen.Line(),
				jen.Comment("fetch the last inserted ID"),
				jen.If(jen.List(jen.ID("id"), jen.ID("idErr")).Assign().ID("res").Dot("LastInsertId").Call()).Op(";").ID("idErr").IsEqualTo().ID("nil").Block(
					jen.ID("x").Dot("ID").Equals().Uint64().Call(jen.ID("id")),
					jen.Line(),
					jen.List(jen.ID("query"), jen.ID("args")).Equals().ID(dbfl).Dot("buildOAuth2ClientCreationTimeQuery").Call(jen.ID("x").Dot("ID")),
					jen.ID(dbfl).Dot("logCreationTimeRetrievalError").Call(jen.ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("CreatedOn"))),
				),
			)
		}
		out = append(out, jen.Line(), jen.Return().List(jen.ID("x"), jen.Nil()))

		return out
	}

	ret.Add(
		jen.Comment("CreateOAuth2Client creates an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("CreateOAuth2Client").Params(utils.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			buildCreateOauth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildBuildCreateOauth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{jen.Var().Err().Error(), jen.Line()}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Set").Call(jen.Lit("client_id"), jen.ID("input").Dot("ClientID")).
			Dotln("Set").Call(jen.Lit("client_secret"), jen.ID("input").Dot("ClientSecret")).
			Dotln("Set").Call(jen.Lit("scopes"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Scopes"), jen.ID("scopesSeparator"))).
			Dotln("Set").Call(jen.Lit("redirect_uri"), jen.ID("input").Dot("RedirectURI")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").MapAssign().ID("input").Dot("ID"),
			jen.ID("oauth2ClientsTableOwnershipColumn").MapAssign().ID("input").Dot("BelongsToUser"),
		))

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING updated_on"))
		}
		q.Dotln("ToSql").Call()

		out = append(out, q,
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		)

		return out
	}

	ret.Add(
		jen.Comment("buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildUpdateOAuth2ClientQuery").Params(jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			buildBuildCreateOauth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildUpdateOAuth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateOAuth2ClientQuery").Call(jen.ID("input")),
		}

		if isPostgres {
			out = append(out, jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("input").Dot("UpdatedOn")))
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.Return().Err(),
			)
		}

		return out
	}

	ret.Add(
		jen.Comment("UpdateOAuth2Client updates a OAuth2 client."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the input's ID field to be valid and non-zero."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("UpdateOAuth2Client").Params(utils.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.Error()).Block(
			buildUpdateOAuth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildBuildArchiveOAuth2ClientQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").MapAssign().ID("clientID"),
			jen.ID("oauth2ClientsTableOwnershipColumn").MapAssign().ID("userID"),
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
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildArchiveOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			buildBuildArchiveOAuth2ClientQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("ArchiveOAuth2Client").Params(utils.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveOAuth2ClientQuery").Call(jen.ID("clientID"), jen.ID("userID")),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return().Err(),
		),
		jen.Line(),
	)
	return ret
}
