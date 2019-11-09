package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(pkgRoot string, types []models.DataType, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.RouteName())

	utils.AddImports(pkgRoot, types, ret)

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
			jen.Comment("CountQuery is a generic counter query used in a few query builders"),
			jen.ID("CountQuery").Op("=").Lit("COUNT(id)"),
			jen.Line(),
			jen.Commentf("CurrentUnixTimeQuery is the query %s uses to determine the current unix time", cn),
			jen.ID("CurrentUnixTimeQuery").Op("=").Lit(getTimeQuery(dbrn)),
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
		jen.Var().ID("_").Qual(filepath.Join(pkgRoot, "database/v1"), "Database").Op("=").Params(jen.Op("*").ID(sn)).Params(jen.ID("nil")),
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
				jen.ID("ExecContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")),
				jen.ID("QueryContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")),
				jen.ID("QueryRowContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row")),
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
		jen.Func().IDf("Provide%s%s", sn, dbTrail).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("connectionDetails").Qual(filepath.Join(pkgRoot, "database/v1"), "ConnectionDetails")).Params(jen.Op("*").Qual("database/sql", "DB"), jen.ID("error")).Block(
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
		jen.Func().IDf("Provide%s", sn).Params(jen.ID("debug").ID("bool"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkgRoot, "database/v1"), "Database")).Block(
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
					jen.ID("err").Op(":=").ID(dbfl).Dot("db").Dot("Ping").Call(),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
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
					jen.ID("err").Op(":=").ID(dbfl).Dot("db").Dot("Ping").Call(),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
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
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
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
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("logQueryBuildingError").Params(jen.ID("err").ID("error")).Block(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID(dbfl).Dot("logger").Dot("WithName").Call(jen.Lit("QUERY_ERROR")).Dot("Error").Call(jen.ID("err"), jen.Lit("building query")),
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
			jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("logCreationTimeRetrievalError").Params(jen.ID("err").ID("error")).Block(
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID(dbfl).Dot("logger").Dot("WithName").Call(jen.Lit("CREATION_TIME_RETRIEVAL")).Dot("Error").Call(jen.ID("err"), jen.Lit("retrieving creation time")),
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
		jen.Func().ID("buildError").Params(jen.ID("err").ID("error"), jen.ID("msg").ID("string")).Params(jen.ID("error")).Block(
			jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().ID("err"),
			),
			jen.Line(),
			jen.If(jen.Op("!").Qual("strings", "Contains").Call(jen.ID("msg"), jen.RawString(`%w`))).Block(
				jen.ID("msg").Op("+=").Lit(": %w"),
			),
			jen.Line(),
			jen.Return().Qual("fmt", "Errorf").Call(jen.ID("msg"), jen.ID("err")),
		),
		jen.Line(),
	)

	var f32f, f64f bool
	for _, typ := range types {
		for _, field := range typ.Fields {
			if field.Type == "float32" {
				f32f = true
			} else if field.Type == "float64" {
				f64f = true
			}
			if f32f && f64f {
				break
			}
		}
	}

	if f32f || f64f {
		ret.Add(
			jen.Func().ID("getFormatString").Params(jen.ID("divisor").ID("int64")).Params(jen.ID("string")).Block(
				jen.Switch(jen.ID("divisor")).Block(
					jen.Case(jen.Lit(100)).Block(
						jen.Return(jen.RawString(`%.2f`)),
					),
					jen.Case(jen.Lit(1000)).Block(
						jen.Return(jen.RawString(`%.3f`)),
					),
					jen.Case(jen.Lit(10000)).Block(
						jen.Return(jen.RawString(`%.4f`)),
					),
					jen.Case(jen.Lit(100000)).Block(
						jen.Return(jen.RawString(`%.5f`)),
					),
					jen.Case(jen.Lit(1000000)).Block(
						jen.Return(jen.RawString(`%.6f`)),
					),
					jen.Case(jen.Lit(10000000)).Block(
						jen.Return(jen.RawString(`%.7f`)),
					),
					jen.Case(jen.Lit(100000000)).Block(
						jen.Return(jen.RawString(`%.8f`)),
					),
					jen.Case(jen.Lit(1000000000)).Block(
						jen.Return(jen.RawString(`%.9f`)),
					),
					jen.Default().Block(
						jen.Return(jen.Lit("")),
					),
				),
			),
			jen.Line(),
		)
	}

	if f32f {
		ret.Add(
			jen.Func().ID("truncateFloat32").Params(jen.ID("x").ID("float32"), jen.ID("divisor").ID("int64")).Params(jen.ID("float32")).Block(
				jen.List(jen.ID("s"), jen.ID("err")).Op(":=").Qual("strconv", "ParseFloat").Call(jen.Qual("fmt", "Sprintf").Call(jen.ID("getFormatString").Call(jen.ID("divisor")), jen.ID("x")), jen.Lit(32)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Panic(jen.ID("err")),
				),
				jen.ID("dbr").Op(":=").ID("int64").Call(jen.ID("s").Op("*").ID("float64").Call(jen.ID("divisor"))),
				jen.Return().ID("float32").Call(jen.ID("dbr")).Op("/").ID("float32").Call(jen.ID("divisor")),
			),
			jen.Line(),
		)
	}

	if f64f {
		ret.Add(
			jen.Func().ID("truncateFloat64").Params(jen.ID("x").ID("float64"), jen.ID("divisor").ID("int64")).Params(jen.ID("float64")).Block(
				jen.List(jen.ID("s"), jen.ID("err")).Op(":=").Qual("strconv", "ParseFloat").Call(jen.Qual("fmt", "Sprintf").Call(jen.ID("getFormatString").Call(jen.ID("divisor")), jen.ID("x")), jen.Lit(64)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Panic(jen.ID("err")),
				),
				jen.ID("dbr").Op(":=").ID("int64").Call(jen.ID("s").Op("*").ID("float64").Call(jen.ID("divisor"))),
				jen.Return().ID("float64").Call(jen.ID("dbr")).Op("/").ID("float64").Call(jen.ID("divisor")),
			),
			jen.Line(),
		)
	}

	return ret
}
