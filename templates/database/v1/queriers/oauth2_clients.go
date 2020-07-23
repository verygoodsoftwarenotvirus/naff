package queriers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	sn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(sn)[0]))
	dbrn := dbvendor.RouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	code.Add(
		jen.Const().Defs(
			jen.ID("scopesSeparator").Equals().Lit(","),
			jen.ID("oauth2ClientsTableName").Equals().Lit("oauth2_clients"),
			jen.ID("oauth2ClientsTableNameColumn").Equals().Lit("name"),
			jen.ID("oauth2ClientsTableClientIDColumn").Equals().Lit("client_id"),
			jen.ID("oauth2ClientsTableScopesColumn").Equals().Lit("scopes"),
			jen.ID("oauth2ClientsTableRedirectURIColumn").Equals().Lit("redirect_uri"),
			jen.ID("oauth2ClientsTableClientSecretColumn").Equals().Lit("client_secret"),
			jen.ID("oauth2ClientsTableOwnershipColumn").Equals().Lit("belongs_to_user"),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Var().Defs(
			jen.ID("oauth2ClientsTableColumns").Equals().Index().String().Valuesln(
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("idColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableNameColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableClientIDColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableScopesColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableRedirectURIColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableClientSecretColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("createdOnColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("lastUpdatedOnColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("archivedOnColumn")),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableOwnershipColumn")),
			),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("scanOAuth2Client").Params(
			jen.ID("scan").Qual(proj.DatabaseV1Package(), "Scanner"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("x").Equals().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(),
				jen.ID("scopes").String(),
			),
			jen.Line(),
			jen.ID("targetVars").Assign().Index().Interface().Valuesln(
				jen.AddressOf().ID("x").Dot("ID"),
				jen.AddressOf().ID("x").Dot("Name"),
				jen.AddressOf().ID("x").Dot("ClientID"),
				jen.AddressOf().ID("scopes"), jen.AddressOf().ID("x").Dot("RedirectURI"),
				jen.AddressOf().ID("x").Dot("ClientSecret"),
				jen.AddressOf().ID("x").Dot("CreatedOn"),
				jen.AddressOf().ID("x").Dot("LastUpdatedOn"),
				jen.AddressOf().ID("x").Dot("ArchivedOn"),
				jen.AddressOf().ID("x").Dot(constants.UserOwnershipFieldName),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Spread()), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.If(
				jen.ID("scopes").Assign().Qual("strings", "Split").Call(jen.ID("scopes"), jen.ID("scopesSeparator")),
				jen.Len(jen.ID("scopes")).Op(">=").One().And().ID("scopes").Index(jen.Zero()).DoesNotEqual().EmptyString(),
			).Block(
				jen.ID("x").Dot("Scopes").Equals().ID("scopes"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("scanOAuth2Clients").Params(
			jen.ID("rows").Qual(proj.DatabaseV1Package(), "ResultIterator"),
		).Params(
			jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("list").Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("client"), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Client").Call(
					jen.ID("rows"),
				),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Line(),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.ID("client")),
			),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID."),
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
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableClientIDColumn")).MapAssign().ID("clientID"),
				utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("archivedOnColumn")).MapAssign().ID("nil"),
			)).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("GetOAuth2ClientByClientID gets an OAuth2 client."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetOAuth2ClientByClientID").Params(
			constants.CtxParam(),
			jen.ID("clientID").String(),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error(),
		).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientByClientIDQuery").Call(jen.ID("clientID")),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.Return().ID(dbfl).Dot("scanOAuth2Client").Call(jen.ID("row")),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientsQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientsQuery").String(),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("buildGetAllOAuth2ClientsQuery builds a SQL query."),
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
						utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("archivedOnColumn")).MapAssign().Nil(),
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

	code.Add(
		jen.Comment("GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllOAuth2Clients").Params(constants.CtxParam()).Params(jen.Index().PointerTo().Qual(proj.ModelsV1Package(),
			"OAuth2Client",
		),
			jen.Error()).Block(
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID(dbfl).Dot("buildGetAllOAuth2ClientsQuery").Call()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllOAuth2ClientsForUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientsForUserQuery").Call(jen.ID(constants.UserIDVarName), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying database for oauth2 clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching list of OAuth2Clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Spread()).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("idColumn")).MapAssign().ID("clientID"),
					utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableOwnershipColumn")).MapAssign().ID(constants.UserIDVarName),
					utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("archivedOnColumn")).MapAssign().Nil(),
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

	code.Add(
		jen.Comment("GetOAuth2Client retrieves an OAuth2 client from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientQuery").Call(jen.ID("clientID"), jen.ID(constants.UserIDVarName)),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("client"), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Client").Call(jen.ID("row")),
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

	code.Add(
		jen.Var().Defs(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Qual("sync", "Once"),
			jen.ID("getAllOAuth2ClientCountQuery").String(),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("buildGetAllOAuth2ClientsCountQuery returns a SQL query for the number of OAuth2 clients"),
		jen.Line(),
		jen.Comment("in the database, regardless of ownership."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllOAuth2ClientsCountQuery").Params().Params(jen.String()).Block(
			jen.ID("getAllOAuth2ClientCountQueryBuilder").Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.ID("getAllOAuth2ClientCountQuery"), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(utils.FormatStringWithArg(jen.ID("countQuery"), jen.ID("oauth2ClientsTableName"))).
					Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
					Dotln("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
						utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("archivedOnColumn")).MapAssign().ID("nil"),
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

	code.Add(
		jen.Comment("GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllOAuth2ClientCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.Var().ID("count").Uint64(),
			jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID(dbfl).Dot("buildGetAllOAuth2ClientsCountQuery").Call()).Dot("Scan").Call(jen.AddressOf().ID("count")),
			jen.Return().List(jen.ID("count"), jen.Err()),
		),
		jen.Line(),
	)

	////////////

	code.Add(
		jen.Comment("buildGetOAuth2ClientsForUserQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that"),
		jen.Line(),
		jen.Comment("meet the given filter's criteria (if relevant) and belong to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetOAuth2ClientsForUserQuery").Params(
			constants.UserIDParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("oauth2ClientsTableColumns").Spread()).
				Dotln("From").Call(jen.ID("oauth2ClientsTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("oauth2ClientsTableOwnershipColumn")).MapAssign().ID(constants.UserIDVarName),
					utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("archivedOnColumn")).MapAssign().Nil(),
				),
			).Dotln("OrderBy").Call(utils.FormatString("%s.%s", jen.ID("oauth2ClientsTableName"), jen.ID("idColumn"))),
			jen.Line(),
			jen.If(jen.ID(constants.FilterVarName).DoesNotEqual().ID("nil")).Block(
				jen.ID("builder").Equals().ID(constants.FilterVarName).Dot("ApplyToQueryBuilder").Call(jen.ID("builder"), jen.ID("oauth2ClientsTableName")),
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

	code.Add(
		jen.Comment("GetOAuth2ClientsForUser gets a list of OAuth2 clients."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetOAuth2ClientsForUser").Params(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList"),
			jen.Error(),
		).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetOAuth2ClientsForUserQuery").Call(jen.ID(constants.UserIDVarName), jen.ID(constants.FilterVarName)),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for oauth2 clients: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().ID(dbfl).Dot("scanOAuth2Clients").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("ocl").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID(constants.FilterVarName).Dot("Page"),
					jen.ID("Limit").MapAssign().ID(constants.FilterVarName).Dot("Limit"),
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
			jen.ID("oauth2ClientsTableNameColumn"),
			jen.ID("oauth2ClientsTableClientIDColumn"),
			jen.ID("oauth2ClientsTableClientSecretColumn"),
			jen.ID("oauth2ClientsTableScopesColumn"),
			jen.ID("oauth2ClientsTableRedirectURIColumn"),
			jen.ID("oauth2ClientsTableOwnershipColumn"),
		}

		vals := []jen.Code{
			jen.ID("input").Dot("Name"),
			jen.ID("input").Dot("ClientID"),
			jen.ID("input").Dot("ClientSecret"),
			jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Scopes"), jen.ID("scopesSeparator")),
			jen.ID("input").Dot("RedirectURI"),
			jen.ID("input").Dot(constants.UserOwnershipFieldName),
		}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Insert").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Columns").Callln(cols...).
			Dotln("Values").Callln(vals...)
		if isPostgres {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s, %s"), jen.ID("idColumn"), jen.ID("createdOnColumn")))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	code.Add(
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

	buildCreateOauth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
				jen.ID("Name").MapAssign().ID("input").Dot("Name"),
				jen.ID("ClientID").MapAssign().ID("input").Dot("ClientID"),
				jen.ID("ClientSecret").MapAssign().ID("input").Dot("ClientSecret"),
				jen.ID("RedirectURI").MapAssign().ID("input").Dot("RedirectURI"),
				jen.ID("Scopes").MapAssign().ID("input").Dot("Scopes"),
				jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("input").Dot(constants.UserOwnershipFieldName)),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildCreateOAuth2ClientQuery").Call(jen.ID("x")),
			jen.Line(),
		}

		if isPostgres {
			out = append(out,
				jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("ID"), jen.AddressOf().ID("x").Dot("CreatedOn")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing client creation query: %w"), jen.Err())),
				),
			)
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing client creation query: %w"), jen.Err())),
				),
				jen.Line(),
				jen.Comment("fetch the last inserted ID."),
				jen.List(jen.ID("id"), jen.ID("err")).Assign().ID(constants.ResponseVarName).Dot("LastInsertId").Call(),
				jen.ID(dbfl).Dot("logIDRetrievalError").Call(jen.Err()),
				jen.ID("x").Dot("ID").Equals().Uint64().Call(jen.ID("id")),
				jen.Line(),
				jen.Comment("this won't be completely accurate, but it will suffice."),
				jen.ID("x").Dot("CreatedOn").Equals().ID(dbfl).Dot("timeTeller").Dot("Now").Call(),
			)
		}
		out = append(out, jen.Line(), jen.Return().List(jen.ID("x"), jen.Nil()))

		return out
	}

	code.Add(
		jen.Comment("CreateOAuth2Client creates an OAuth2 client."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("CreateOAuth2Client").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			buildCreateOauth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildBuildCreateOauth2ClientTestBody := func() []jen.Code {
		out := []jen.Code{jen.Var().Err().Error(), jen.Line()}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Set").Call(jen.ID("oauth2ClientsTableClientIDColumn"), jen.ID("input").Dot("ClientID")).
			Dotln("Set").Call(jen.ID("oauth2ClientsTableClientSecretColumn"), jen.ID("input").Dot("ClientSecret")).
			Dotln("Set").Call(jen.ID("oauth2ClientsTableScopesColumn"), jen.Qual("strings", "Join").Call(jen.ID("input").Dot("Scopes"), jen.ID("scopesSeparator"))).
			Dotln("Set").Call(jen.ID("oauth2ClientsTableRedirectURIColumn"), jen.ID("input").Dot("RedirectURI")).
			Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.ID("idColumn").MapAssign().ID("input").Dot("ID"),
			jen.ID("oauth2ClientsTableOwnershipColumn").MapAssign().ID("input").Dot(constants.UserOwnershipFieldName),
		))

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.ID("lastUpdatedOnColumn")))
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

	code.Add(
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
			out = append(out, jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("input").Dot("LastUpdatedOn")))
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.Return().Err(),
			)
		}

		return out
	}

	code.Add(
		jen.Comment("UpdateOAuth2Client updates a OAuth2 client."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the input's ID field to be valid and non-zero."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("UpdateOAuth2Client").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.Error()).Block(
			buildUpdateOAuth2ClientTestBody()...,
		),
		jen.Line(),
	)

	////////////

	buildBuildArchiveOAuth2ClientQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("oauth2ClientsTableName")).
			Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Set").Call(jen.ID("archivedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.ID("idColumn").MapAssign().ID("clientID"),
			jen.ID("oauth2ClientsTableOwnershipColumn").MapAssign().ID(constants.UserIDVarName),
		))

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.ID("archivedOnColumn")))
		}
		q.Dotln("ToSql").Call()

		return q
	}

	code.Add(
		jen.Comment("buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildArchiveOAuth2ClientQuery").Params(jen.List(jen.ID("clientID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
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

	code.Add(
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("ArchiveOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveOAuth2ClientQuery").Call(jen.ID("clientID"), jen.ID(constants.UserIDVarName)),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return().Err(),
		),
		jen.Line(),
	)

	return code
}
