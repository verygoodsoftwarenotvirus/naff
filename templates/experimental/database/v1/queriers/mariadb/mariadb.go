package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mariadbDotGo() *jen.File {
	ret := jen.NewFile("mariadb")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("loggerName").Op("=").Lit("mariadb"),
			jen.ID("mariaDBDriverName").Op("=").Lit("wrapped-mariadb-driver"),
			jen.ID("CountQuery").Op("=").Lit("COUNT(id)"),
			jen.ID("CurrentUnixTimeQuery").Op("=").Lit("UNIX_TIMESTAMP()"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.ID("driver").Op(":=").Qual("contrib.go.opencensus.io/integrations/ocsql", "Wrap").Call(
				jen.Op("&").Qual("github.com/go-sql-driver/mysql", "MySQLDriver").Values(),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithQuery").Call(jen.ID("true")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithAllowRoot").Call(jen.ID("false")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithRowsNext").Call(jen.ID("true")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithRowsClose").Call(jen.ID("true")),
				jen.Qual("contrib.go.opencensus.io/integrations/ocsql", "WithQueryParams").Call(jen.ID("true")),
			),
			jen.Qual("database/sql", "Register").Call(jen.ID("mariaDBDriverName"), jen.ID("driver")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.ID("MariaDB").Struct(
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("db").Op("*").Qual("database/sql", "DB"),
				jen.ID("sqlBuilder").ID("squirrel").Dot("StatementBuilderType"),
				jen.ID("migrateOnce").Qual("sync", "Once"), jen.ID("debug").ID("bool"),
				jen.ID("ConnectionDetails").ID("string"),
				jen.ID("Querier").Interface(
					jen.ID("ExecContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")),
					jen.ID("QueryContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")),
					jen.ID("QueryRowContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideMariaDBConnection provides an instrumented mariadb connection"),
		jen.Line(),
		jen.Func().ID("ProvideMariaDBConnection").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("connectionDetails").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "ConnectionDetails")).Params(jen.Op("*").Qual("database/sql", "DB"), jen.ID("error")).Block(
			jen.ID("logger").Dot("WithValue").Call(jen.Lit("connection_details"), jen.ID("connectionDetails")).Dot("Debug").Call(jen.Lit("Establishing connection to mariadb")),
			jen.Return().Qual("database/sql", "Open").Call(jen.ID("mariaDBDriverName"), jen.ID("string").Call(jen.ID("connectionDetails"))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideMariaDB provides a mariadb controller"),
		jen.Line(),
		jen.Func().ID("ProvideMariaDB").Params(jen.ID("debug").ID("bool"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Database")).Block(
			jen.Return().Op("&").ID("MariaDB").Valuesln(
				jen.ID("db").Op(":").ID("db"),
				jen.ID("debug").Op(":").ID("debug"),
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("loggerName")),
				jen.ID("sqlBuilder").Op(":").ID("squirrel").Dot("StatementBuilder"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsReady reports whether or not the db is ready"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			jen.ID("numberOfUnsuccessfulAttempts").Op(":=").Lit(0),
			jen.ID("waitInterval").Op(":=").Qual("time", "Second"),
			jen.ID("maxAttempts").Op(":=").Lit(100),
			jen.ID("m").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("wait_interval").Op(":").ID("waitInterval"), jen.Lit("max_attempts").Op(":").ID("maxAttempts"))).Dot("Debug").Call(jen.Lit("IsReady called")),
			jen.For(jen.Op("!").ID("ready")).Block(
				jen.ID("err").Op(":=").ID("m").Dot("db").Dot(
					"Ping",
				).Call(),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("m").Dot("logger").Dot("Debug").Call(jen.Lit("ping failed, waiting for db")),
					jen.Qual("time", "Sleep").Call(jen.ID("waitInterval")),
					jen.ID("numberOfUnsuccessfulAttempts").Op("++"),
					jen.If(jen.ID("numberOfUnsuccessfulAttempts").Op(">=").ID("maxAttempts")).Block(
						jen.Return().ID("false"),
					),
				).Else().Block(
					jen.ID("ready").Op("=").ID("true"),
					jen.Return().ID("ready"),
				),
			),
			jen.Return().ID("false"),
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
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("logQueryBuildingError").Params(jen.ID("err").ID("error")).Block(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("m").Dot("logger").Dot("WithName").Call(jen.Lit("QUERY_ERROR")).Dot("Error").Call(jen.ID("err"), jen.Lit("building query")),
			),
		),
		jen.Line(),
	)

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
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("logCreationTimeRetrievalError").Params(jen.ID("err").ID("error")).Block(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("m").Dot("logger").Dot("WithName").Call(jen.Lit("CREATION_TIME_RETRIEVAL")).Dot("Error").Call(jen.ID("err"), jen.Lit("building query")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildError takes a given error and wraps it with a message, provided that it"),
		jen.Line(),
		jen.Comment("IS NOT sql.ErrNoRows, which we want to preserve and surface to the services."),
		jen.Line(),
		jen.Func().ID("buildError").Params(jen.ID("err").ID("error"), jen.ID("msg").ID("string")).Params(jen.ID("error")).Block(
			jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().ID("err"),
			),
			jen.Return().ID("errors").Dot(
				"Wrap",
			).Call(jen.ID("err"), jen.ID("msg")),
		),
		jen.Line(),
	)
	return ret
}
