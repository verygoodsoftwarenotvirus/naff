package queriers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	uvn := dbvendor.UnexportedVarName()
	cn := dbvendor.SingularCommonName()
	sn := dbvendor.Singular()
	rn := dbvendor.RouteName()
	dbrn := rn
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	squirrelPlaceholder := "Question"
	if dbrn == "postgres" {
		squirrelPlaceholder = "Dollar"
	}
	squirrelInitConfig := jen.ID("sqlBuilder").MapAssign().Qual("github.com/Masterminds/squirrel", "StatementBuilder").Dot("PlaceholderFormat").Call(jen.Qual("github.com/Masterminds/squirrel", squirrelPlaceholder))

	code.Add(
		jen.Const().Defs(
			jen.ID("loggerName").Equals().Lit(rn),
			jen.IDf("%sDriverName", uvn).Equals().Litf("wrapped-%s-driver", dbvendor.KebabName()),
			func() jen.Code {
				if isPostgres(dbvendor) {
					g := &jen.Group{}
					g.Add(
						jen.Line(),
						jen.ID("postgresRowExistsErrorCode").Equals().Lit("23505"),
					)
					return g
				}
				return jen.Null()
			}(),
			jen.Line(),
			jen.List(jen.ID("existencePrefix"), jen.ID("existenceSuffix")).Equals().List(jen.Lit("SELECT EXISTS ("), jen.Lit(")")),
			jen.Line(),
			jen.Comment("countQuery is a generic counter query used in a few query builders."),
			jen.ID("countQuery").Equals().Lit("COUNT(%s.id)"),
			jen.Line(),
			jen.Commentf("currentUnixTimeQuery is the query %s uses to determine the current unix time.", cn),
			jen.ID("currentUnixTimeQuery").Equals().Lit(getTimeQuery(dbvendor)),
		),
		jen.Line(),
	)

	////////////

	var driverInit jen.Code
	if isPostgres(dbvendor) {
		driverInit = jen.AddressOf().Qual("github.com/lib/pq", "Driver").Values()
	} else if isSqlite(dbvendor) {
		driverInit = jen.AddressOf().Qual("github.com/mattn/go-sqlite3", "SQLiteDriver").Values()
	} else if isMariaDB(dbvendor) {
		driverInit = jen.AddressOf().Qual("github.com/go-sql-driver/mysql", "MySQLDriver").Values()
	}

	code.Add(
		jen.Func().ID("init").Params().Block(
			jen.Commentf("Explicitly wrap the %s driver with ocsql.", sn),
			jen.ID("driver").Assign().Qual("contrib.go.opencensus.io/integrations/ocsql", "Wrap").Callln(
				driverInit,
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithQuery").Call(jen.True()),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithAllowRoot").Call(jen.False()),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithRowsNext").Call(jen.True()),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithRowsClose").Call(jen.True()),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithQueryParams").Call(jen.True()),
			),
			jen.Line(),
			jen.Comment("Register our ocsql wrapper as a db driver."),
			jen.Qual("database/sql", "Register").Call(jen.IDf("%sDriverName", uvn), jen.ID("driver")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Underscore().Qual(proj.DatabaseV1Package(), "DataManager").Equals().Params(jen.PointerTo().ID(sn)).Params(jen.Nil()),
		jen.Line(),
		jen.Type().Defs(
			jen.Commentf("%s is our main %s interaction db.", sn, sn),
			jen.ID(sn).Struct(
				constants.LoggerParam(),
				jen.ID("db").PointerTo().Qual("database/sql", "DB"),
				func() jen.Code {
					if isMariaDB(dbvendor) || isSqlite(dbvendor) {
						return jen.ID("timeTeller").ID("timeTeller")
					}
					return jen.Null()
				}(),
				jen.ID("sqlBuilder").Qual("github.com/Masterminds/squirrel", "StatementBuilderType"),
				jen.ID("migrateOnce").Qual("sync", "Once"), jen.ID("debug").Bool(),
			),
			jen.Line(),
			jen.Commentf("ConnectionDetails is a string alias for a %s url.", sn),
			jen.ID("ConnectionDetails").String(),
			jen.Line(),
			jen.Comment("Querier is a subset interface for sql.{DB|Tx|Stmt} objects"),
			jen.ID("Querier").Interface(
				jen.ID("ExecContext").Params(constants.CtxParam(), jen.ID("args").Spread().Interface()).Params(jen.Qual("database/sql", "Result"), jen.Error()),
				jen.ID("QueryContext").Params(constants.CtxParam(), jen.ID("args").Spread().Interface()).Params(jen.PointerTo().Qual("database/sql", "Rows"), jen.Error()),
				jen.ID("QueryRowContext").Params(constants.CtxParam(), jen.ID("args").Spread().Interface()).Params(jen.PointerTo().Qual("database/sql", "Row")),
			),
		),
		jen.Line(),
	)

	////////////
	var (
		dbTrail string
	)
	if !isMariaDB(dbvendor) {
		dbTrail = "DB"
	} else {
		dbTrail = "Connection"
	}

	code.Add(
		jen.Commentf("Provide%s%s provides an instrumented %s db.", sn, dbTrail, cn),
		jen.Line(),
		jen.Func().IDf("Provide%s%s", sn, dbTrail).Params(constants.LoggerParam(), jen.ID("connectionDetails").Qual(proj.DatabaseV1Package(), "ConnectionDetails")).Params(jen.PointerTo().Qual("database/sql", "DB"), jen.Error()).Block(
			jen.ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("connection_details"), jen.ID("connectionDetails")).Dot("Debug").Call(jen.Litf("Establishing connection to %s", cn)),
			jen.Return().Qual("database/sql", "Open").Call(jen.IDf("%sDriverName", uvn), jen.String().Call(jen.ID("connectionDetails"))),
		),
		jen.Line(),
	)

	////////////
	dbTrail = ""
	if !isMariaDB(dbvendor) {
		dbTrail = " db"
	}

	code.Add(
		jen.Commentf("Provide%s provides a %s%s controller.", sn, cn, dbTrail),
		jen.Line(),
		jen.Func().IDf("Provide%s", sn).Params(jen.ID("debug").Bool(), jen.ID("db").PointerTo().Qual("database/sql", "DB"), constants.LoggerParam()).Params(jen.Qual(proj.DatabaseV1Package(), "DataManager")).Block(
			jen.Return().AddressOf().IDf(sn).Valuesln(
				jen.ID("db").MapAssign().ID("db"),
				jen.ID("debug").MapAssign().ID("debug"),
				func() jen.Code {
					if isMariaDB(dbvendor) || isSqlite(dbvendor) {
						return jen.ID("timeTeller").MapAssign().AddressOf().ID("stdLibTimeTeller").Values()
					}
					return jen.Null()
				}(),
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("loggerName")),
				squirrelInitConfig,
			),
		),
		jen.Line(),
	)

	buildIsReadyBody := func() []jen.Code {
		if isSqlite(dbvendor) {
			return []jen.Code{jen.Return(jen.True())}
		} else if isPostgres(dbvendor) {
			return []jen.Code{
				jen.ID("numberOfUnsuccessfulAttempts").Assign().Zero(),
				jen.Line(),
				jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("interval").MapAssign().Qual("time", "Second"),
					jen.Lit("max_attempts").MapAssign().Lit(50)),
				).Dot("Debug").Call(jen.Lit("IsReady called")),
				jen.Line(),
				jen.For(jen.Not().ID("ready")).Block(
					jen.Err().Assign().ID(dbfl).Dot("db").Dot("PingContext").Call(constants.CtxVar()),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("ping failed, waiting for db")),
						jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
						jen.Line(),
						jen.ID("numberOfUnsuccessfulAttempts").Op("++"),
						jen.If(jen.ID("numberOfUnsuccessfulAttempts").Op(">=").Lit(50)).Block(
							jen.Return().False(),
						),
					).Else().Block(
						jen.ID("ready").Equals().True(),
						jen.Return().ID("ready"),
					),
				),
				jen.Return().False(),
			}
		} else if isMariaDB(dbvendor) {
			return []jen.Code{
				jen.ID("numberOfUnsuccessfulAttempts").Assign().Zero(),
				jen.Line(),
				jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("interval").MapAssign().Qual("time", "Second"),
					jen.Lit("max_attempts").MapAssign().Lit(50)),
				).Dot("Debug").Call(jen.Lit("IsReady called")),
				jen.Line(),
				jen.For(jen.Not().ID("ready")).Block(
					jen.Err().Assign().ID(dbfl).Dot("db").Dot("PingContext").Call(constants.CtxVar()),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("ping failed, waiting for db")),
						jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
						jen.Line(),
						jen.ID("numberOfUnsuccessfulAttempts").Op("++"),
						jen.If(jen.ID("numberOfUnsuccessfulAttempts").Op(">=").Lit(50)).Block(
							jen.Return().False(),
						),
					).Else().Block(
						jen.ID("ready").Equals().True(),
						jen.Return().ID("ready"),
					),
				),
				jen.Return().False(),
			}
		}
		return nil
	}

	code.Add(
		jen.Comment("IsReady reports whether or not the db is ready."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("IsReady").Params(
			func() jen.Code {
				if !isSqlite(dbvendor) {
					return constants.CtxParam()
				}
				return jen.Underscore().Qual("context", "Context")
			}(),
		).Params(jen.ID("ready").Bool()).Block(
			buildIsReadyBody()...,
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("logQueryBuildingError logs errors that may occur during query construction."),
		jen.Line(),
		jen.Comment("Such errors should be few and far between, as the generally only occur with"),
		jen.Line(),
		jen.Comment("type discrepancies or other misuses of SQL. An alert should be set up for"),
		jen.Line(),
		jen.Comment("any log entries with the given name, and those alerts should be investigated"),
		jen.Line(),
		jen.Comment("with the utmost priority."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("logQueryBuildingError").Params(jen.Err().Error()).Block(
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("WithName").Call(jen.Lit("QUERY_ERROR")).Dot("Error").Call(jen.Err(), jen.Lit("building query")),
			),
		),
		jen.Line(),
	)

	if isMariaDB(dbvendor) || isSqlite(dbvendor) {
		code.Add(
			jen.Comment("logIDRetrievalError logs errors that may occur during created db row ID retrieval."),
			jen.Line(),
			jen.Comment("Such errors should be few and far between, as the generally only occur with"),
			jen.Line(),
			jen.Comment("type discrepancies or other misuses of SQL. An alert should be set up for"),
			jen.Line(),
			jen.Comment("any log entries with the given name, and those alerts should be investigated"),
			jen.Line(),
			jen.Comment("with the utmost priority."),
			jen.Line(),
			jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("logIDRetrievalError").Params(jen.Err().Error()).Block(
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("WithName").Call(jen.Lit("ROW_ID_ERROR")).Dot("Error").Call(jen.Err(), jen.Lit("fetching row ID")),
				),
			),
			jen.Line(),
		)
	}

	code.Add(
		jen.Comment("buildError takes a given error and wraps it with a message, provided that it"),
		jen.Line(),
		jen.Comment("IS NOT sql.ErrNoRows, which we want to preserve and surface to the services."),
		jen.Line(),
		jen.Func().ID("buildError").Params(jen.Err().Error(), jen.ID("msg").String()).Params(jen.Error()).Block(
			jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().Err(),
			),
			jen.Line(),
			jen.If(jen.Not().Qual("strings", "Contains").Call(jen.ID("msg"), jen.RawString(`%w`))).Block(
				jen.ID("msg").Op("+=").Lit(": %w"),
			),
			jen.Line(),
			jen.Return().Qual("fmt", "Errorf").Call(jen.ID("msg"), jen.Err()),
		),
		jen.Line(),
	)

	return code
}
