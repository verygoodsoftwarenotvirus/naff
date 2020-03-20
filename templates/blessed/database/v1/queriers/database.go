package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(pkg *models.Project, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.RouteName())

	utils.AddImports(pkg, ret)

	uvn := vendor.UnexportedVarName()
	cn := vendor.SingularCommonName()
	sn := vendor.Singular()
	rn := vendor.RouteName()
	dbrn := rn
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	var squirrelInitConfig jen.Code

	if dbrn == "postgres" {
		squirrelInitConfig = jen.ID("sqlBuilder").Op(":").Qual("github.com/Masterminds/squirrel", "StatementBuilder").Dot("PlaceholderFormat").Call(jen.Qual("github.com/Masterminds/squirrel", "Dollar"))
	} else if dbrn == "sqlite" || dbrn == "mariadb" {
		squirrelInitConfig = jen.ID("sqlBuilder").Op(":").Qual("github.com/Masterminds/squirrel", "StatementBuilder")
	}

	ret.Add(
		jen.Const().Defs(
			jen.ID("loggerName").Op("=").Lit(rn),
			jen.IDf("%sDriverName", uvn).Op("=").Litf("wrapped-%s-driver", vendor.KebabName()),
			jen.Line(),
			jen.Comment("countQuery is a generic counter query used in a few query builders"),
			jen.ID("countQuery").Op("=").Lit("COUNT(id)"),
			jen.Line(),
			jen.Commentf("currentUnixTimeQuery is the query %s uses to determine the current unix time", cn),
			jen.ID("currentUnixTimeQuery").Op("=").Lit(getTimeQuery(vendor)),
		),
		jen.Line(),
	)

	////////////

	var driverInit jen.Code
	if isPostgres {
		driverInit = jen.Op("&").Qual("github.com/lib/pq", "Driver").Values()
	} else if isSqlite {
		driverInit = jen.Op("&").Qual("github.com/mattn/go-sqlite3", "SQLiteDriver").Values()
	} else if isMariaDB {
		driverInit = jen.Op("&").Qual("github.com/go-sql-driver/mysql", "MySQLDriver").Values()
	}

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.Commentf("Explicitly wrap the %s driver with ocsql", sn),
			jen.ID("driver").Op(":=").Qual("contrib.go.opencensus.io/integrations/ocsql", "Wrap").Callln(
				driverInit,
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithQuery").Call(jen.ID("true")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithAllowRoot").Call(jen.ID("false")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithRowsNext").Call(jen.ID("true")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithRowsClose").Call(jen.ID("true")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithQueryParams").Call(jen.ID("true")),
			),
			jen.Line(),
			jen.Comment("Register our ocsql wrapper as a db driver"),
			jen.Qual("database/sql", "Register").Call(jen.IDf("%sDriverName", uvn), jen.ID("driver")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database").Op("=").Params(jen.Op("*").ID(sn)).Params(jen.Nil()),
		jen.Line(),
		jen.Type().Defs(
			jen.Commentf("%s is our main %s interaction db", sn, sn),
			jen.ID(sn).Struct(
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("sqlBuilder").Qual("github.com/Masterminds/squirrel", "StatementBuilderType"),
				jen.ID("migrateOnce").Qual("sync", "Once"), jen.ID("debug").ID("bool"),
			),
			jen.Line(),
			jen.Commentf("ConnectionDetails is a string alias for a %s url", sn),
			jen.ID("ConnectionDetails").ID("string"),
			jen.Line(),
			jen.Comment("Querier is a subset interface for sql.{DB|Tx|Stmt} objects"),
			jen.ID("Querier").Interface(
				jen.ID("ExecContext").Params(utils.CtxParam(), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")),
				jen.ID("QueryContext").Params(utils.CtxParam(), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")),
				jen.ID("QueryRowContext").Params(utils.CtxParam(), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row")),
			),
		),
		jen.Line(),
	)

	////////////
	var (
		dbTrail string
	)
	if !isMariaDB {
		dbTrail = "DB"
	} else {
		dbTrail = "Connection"
	}

	ret.Add(
		jen.Commentf("Provide%s%s provides an instrumented %s db", sn, dbTrail, cn),
		jen.Line(),
		jen.Func().IDf("Provide%s%s", sn, dbTrail).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("connectionDetails").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "ConnectionDetails")).Params(jen.Op("*").Qual("database/sql", "DB"), jen.ID("error")).Block(
			jen.ID("logger").Dot("WithValue").Call(jen.Lit("connection_details"), jen.ID("connectionDetails")).Dot("Debug").Call(jen.Litf("Establishing connection to %s", cn)),
			jen.Return().Qual("database/sql", "Open").Call(jen.IDf("%sDriverName", uvn), jen.ID("string").Call(jen.ID("connectionDetails"))),
		),
		jen.Line(),
	)

	////////////
	dbTrail = ""
	if !isMariaDB {
		dbTrail = " db"
	}

	ret.Add(
		jen.Commentf("Provide%s provides a %s%s controller", sn, cn, dbTrail),
		jen.Line(),
		jen.Func().IDf("Provide%s", sn).Params(jen.ID("debug").ID("bool"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database")).Block(
			jen.Return().Op("&").IDf(sn).Valuesln(
				jen.ID("db").Op(":").ID("db"),
				jen.ID("debug").Op(":").ID("debug"),
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("loggerName")),
				squirrelInitConfig,
			),
		),
		jen.Line(),
	)

	buildIsReadyBody := func() []jen.Code {
		if isSqlite {
			return []jen.Code{jen.Return(jen.ID("true"))}
		} else if isPostgres {
			return []jen.Code{
				jen.ID("numberOfUnsuccessfulAttempts").Op(":=").Lit(0),
				jen.Line(),
				jen.ID(dbfl).Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("interval").Op(":").Qual("time", "Second"),
					jen.Lit("max_attempts").Op(":").Lit(50)),
				).Dot("Debug").Call(jen.Lit("IsReady called")),
				jen.Line(),
				jen.For(jen.Op("!").ID("ready")).Block(
					jen.Err().Op(":=").ID(dbfl).Dot("db").Dot("Ping").Call(),
					jen.If(jen.Err().Op("!=").ID("nil")).Block(
						jen.ID(dbfl).Dot("logger").Dot("Debug").Call(jen.Lit("ping failed, waiting for db")),
						jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
						jen.Line(),
						jen.ID("numberOfUnsuccessfulAttempts").Op("++"),
						jen.If(jen.ID("numberOfUnsuccessfulAttempts").Op(">=").Lit(50)).Block(
							jen.Return().ID("false"),
						),
					).Else().Block(
						jen.ID("ready").Op("=").ID("true"),
						jen.Return().ID("ready"),
					),
				),
				jen.Return().ID("false"),
			}
		} else if isMariaDB {
			return []jen.Code{
				jen.ID("numberOfUnsuccessfulAttempts").Op(":=").Lit(0),
				jen.Line(),
				jen.ID(dbfl).Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("interval").Op(":").Qual("time", "Second"),
					jen.Lit("max_attempts").Op(":").Lit(50)),
				).Dot("Debug").Call(jen.Lit("IsReady called")),
				jen.Line(),
				jen.For(jen.Op("!").ID("ready")).Block(
					jen.Err().Op(":=").ID(dbfl).Dot("db").Dot("Ping").Call(),
					jen.If(jen.Err().Op("!=").ID("nil")).Block(
						jen.ID(dbfl).Dot("logger").Dot("Debug").Call(jen.Lit("ping failed, waiting for db")),
						jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
						jen.Line(),
						jen.ID("numberOfUnsuccessfulAttempts").Op("++"),
						jen.If(jen.ID("numberOfUnsuccessfulAttempts").Op(">=").Lit(50)).Block(
							jen.Return().ID("false"),
						),
					).Else().Block(
						jen.ID("ready").Op("=").ID("true"),
						jen.Return().ID("ready"),
					),
				),
				jen.Return().ID("false"),
			}
		}
		return nil
	}

	ret.Add(
		jen.Comment("IsReady reports whether or not the db is ready"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("IsReady").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			buildIsReadyBody()...,
		),
		jen.Line(),
	)

	ret.Add(
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
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("logQueryBuildingError").Params(jen.Err().ID("error")).Block(
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.ID(dbfl).Dot("logger").Dot("WithName").Call(jen.Lit("QUERY_ERROR")).Dot("Error").Call(jen.Err(), jen.Lit("building query")),
			),
		),
		jen.Line(),
	)

	if isSqlite || isMariaDB {
		ret.Add(
			jen.Comment("logCreationTimeRetrievalError logs errors that may occur during creation time retrieval."),
			jen.Line(),
			jen.Comment("Such errors should be few and far between, as the generally only occur with"),
			jen.Line(),
			jen.Comment("type discrepancies or other misuses of SQL. An alert should be set up for"),
			jen.Line(),
			jen.Comment("any log entries with the given name, and those alerts should be investigated"),
			jen.Line(),
			jen.Comment("with the utmost priority."),
			jen.Line(),
			jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("logCreationTimeRetrievalError").Params(jen.Err().ID("error")).Block(
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID(dbfl).Dot("logger").Dot("WithName").Call(jen.Lit("CREATION_TIME_RETRIEVAL")).Dot("Error").Call(jen.Err(), jen.Lit("retrieving creation time")),
				),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Comment("buildError takes a given error and wraps it with a message, provided that it"),
		jen.Line(),
		jen.Comment("IS NOT sql.ErrNoRows, which we want to preserve and surface to the services."),
		jen.Line(),
		jen.Func().ID("buildError").Params(jen.Err().ID("error"), jen.ID("msg").ID("string")).Params(jen.ID("error")).Block(
			jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().ID("err"),
			),
			jen.Line(),
			jen.If(jen.Op("!").Qual("strings", "Contains").Call(jen.ID("msg"), jen.RawString(`%w`))).Block(
				jen.ID("msg").Op("+=").Lit(": %w"),
			),
			jen.Line(),
			jen.Return().Qual("fmt", "Errorf").Call(jen.ID("msg"), jen.Err()),
		),
		jen.Line(),
	)

	return ret
}
