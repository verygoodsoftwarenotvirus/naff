package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func sqliteDotGo() *jen.File {
	ret := jen.NewFile("sqlite")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("loggerName").Op("=").Lit("sqlite").Var().ID("sqliteDriverName").Op("=").Lit("wrapped-sqlite-driver").Var().ID("CountQuery").Op("=").Lit("COUNT(id)").Var().ID("CurrentUnixTimeQuery").Op("=").Lit("(strftime('%s','now'))"),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.ID("driver").Op(":=").Qual("contrib.go.opencensus.io/integrations/ocsql",
				"Wrap",
			).Call(jen.Op("&").Qual("github.com/mattn/go-sqlite3", "SQLiteDriver").Valuesln(), jen.Qual("contrib.go.opencensus.io/integrations/ocsql",
				"WithQuery",
			).Call(jen.ID("true")), jen.Qual("contrib.go.opencensus.io/integrations/ocsql",
				"WithAllowRoot",
			).Call(jen.ID("false")), jen.Qual("contrib.go.opencensus.io/integrations/ocsql",
				"WithRowsNext",
			).Call(jen.ID("true")), jen.Qual("contrib.go.opencensus.io/integrations/ocsql",
				"WithRowsClose",
			).Call(jen.ID("true")), jen.Qual("contrib.go.opencensus.io/integrations/ocsql",
				"WithQueryParams",
			).Call(jen.ID("true"))),
			jen.Qual("database/sql", "Register").Call(jen.ID("sqliteDriverName"), jen.ID("driver")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("database").Dot(
			"Database",
		).Op("=").Parens(jen.Op("*").ID("Sqlite")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("Sqlite").Struct(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
		), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("sqlBuilder").ID("squirrel").Dot(
			"StatementBuilderType",
		), jen.ID("migrateOnce").Qual("sync", "Once"), jen.ID("debug").ID("bool")).Type().ID("ConnectionDetails").ID("string").Type().ID("Querier").Interface(jen.ID("ExecContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")), jen.ID("QueryContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")), jen.ID("QueryRowContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row"))),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideSqliteDB provides an instrumented sqlite db"),
		jen.Line(),
		jen.Func().ID("ProvideSqliteDB").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
		), jen.ID("connectionDetails").ID("database").Dot(
			"ConnectionDetails",
		)).Params(jen.Op("*").Qual("database/sql", "DB"), jen.ID("error")).Block(
			jen.ID("logger").Dot(
				"WithValue",
			).Call(jen.Lit("connection_details"), jen.ID("connectionDetails")).Dot(
				"Debug",
			).Call(jen.Lit("Establishing connection to sqlite")),
			jen.Return().Qual("database/sql", "Open").Call(jen.ID("sqliteDriverName"), jen.ID("string").Call(jen.ID("connectionDetails"))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideSqlite provides a sqlite db controller"),
		jen.Line(),
		jen.Func().ID("ProvideSqlite").Params(jen.ID("debug").ID("bool"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
		)).Params(jen.ID("database").Dot(
			"Database",
		)).Block(
			jen.Return().Op("&").ID("Sqlite").Valuesln(jen.ID("db").Op(":").ID("db"), jen.ID("debug").Op(":").ID("debug"), jen.ID("logger").Op(":").ID("logger").Dot(
				"WithName",
			).Call(jen.ID("loggerName")), jen.ID("sqlBuilder").Op(":").ID("squirrel").Dot(
				"StatementBuilder",
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsReady reports whether or not the db is ready"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			jen.Return().ID("true"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("s.logQueryBuildingError logs errors that may occur during query construction."),
		jen.Line(),
		jen.Func().Comment("// Such errors should be few and far between, as the generally only occur with").Comment("// type discrepancies or other misuses of SQL. An alert should be set up for").Comment("// any log entries with the given name, and those alerts should be investigated").Comment("// with the utmost priority.").Params(jen.ID("s").Op("*").ID("Sqlite")).ID("logQueryBuildingError").Params(jen.ID("err").ID("error")).Block(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"WithName",
				).Call(jen.Lit("QUERY_ERROR")).Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("building query")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("logCreationTimeRetrievalError logs errors that may occur during creation time retrieval"),
		jen.Line(),
		jen.Func().Comment("// Such errors should be few and far between, as the generally only occur with").Comment("// type discrepancies or other misuses of SQL. An alert should be set up for").Comment("// any log entries with the given name, and those alerts should be investigated").Comment("// with the utmost priority.").Params(jen.ID("s").Op("*").ID("Sqlite")).ID("logCreationTimeRetrievalError").Params(jen.ID("err").ID("error")).Block(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"WithName",
				).Call(jen.Lit("CREATION_TIME_RETRIEVAL")).Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("building query")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Comment("// buildError takes a given error and wraps it with a message, provided that it").Comment("// IS NOT sql.ErrNoRows, which we want to preserve and surface to the services.").ID("buildError").Params(jen.ID("err").ID("error"), jen.ID("msg").ID("string")).Params(jen.ID("error")).Block(
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
