package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mariadbDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("loggerName").Op("=").Lit("mariadb"),
			jen.ID("driverName").Op("=").Lit("wrapped-mariadb-driverName"),
			jen.ID("columnCountQueryTemplate").Op("=").Lit(`COUNT(%s.id)`),
			jen.ID("allCountQuery").Op("=").Lit(`COUNT(*)`),
			jen.ID("jsonPluckQuery").Op("=").Lit(`JSON_CONTAINS(%s.%s, '%d', '$.%s')`),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("currentUnixTimeQuery").Op("=").ID("squirrel").Dot("Expr").Call(jen.Lit(`UNIX_TIMESTAMP()`)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("SQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("MariaDB")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("MariaDB").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("sqlBuilder").ID("squirrel").Dot("StatementBuilderType"),
				jen.ID("externalIDGenerator").ID("querybuilding").Dot("ExternalIDGenerator"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("instrumentedDriverRegistration").Qual("sync", "Once"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideMariaDBConnection provides an instrumented maria DB db."),
		jen.Line(),
		jen.Func().ID("ProvideMariaDBConnection").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("connectionDetails").ID("database").Dot("ConnectionDetails")).Params(jen.Op("*").Qual("database/sql", "DB"), jen.ID("error")).Body(
			jen.ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ConnectionDetailsKey"),
				jen.ID("connectionDetails"),
			).Dot("Debug").Call(jen.Lit("Establishing connection to maria DB")),
			jen.ID("instrumentedDriverRegistration").Dot("Do").Call(jen.Func().Params().Body(
				jen.Qual("database/sql", "Register").Call(
					jen.ID("driverName"),
					jen.ID("instrumentedsql").Dot("WrapDriver").Call(
						jen.Op("&").ID("mysql").Dot("MySQLDriver").Valuesln(),
						jen.ID("instrumentedsql").Dot("WithOmitArgs").Call(),
						jen.ID("instrumentedsql").Dot("WithTracer").Call(jen.ID("tracing").Dot("NewInstrumentedSQLTracer").Call(jen.Lit("mariadb_connection"))),
						jen.ID("instrumentedsql").Dot("WithLogger").Call(jen.ID("tracing").Dot("NewInstrumentedSQLLogger").Call(jen.ID("logger"))),
					),
				))),
			jen.Return().Qual("database/sql", "Open").Call(
				jen.Lit("mysql"),
				jen.ID("string").Call(jen.ID("connectionDetails")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideMariaDB provides a maria DB controller."),
		jen.Line(),
		jen.Func().ID("ProvideMariaDB").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.Op("*").ID("MariaDB")).Body(
			jen.Return().Op("&").ID("MariaDB").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("loggerName")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("mariadb_query_builder")), jen.ID("sqlBuilder").Op(":").ID("squirrel").Dot("StatementBuilder").Dot("PlaceholderFormat").Call(jen.ID("squirrel").Dot("Question")), jen.ID("externalIDGenerator").Op(":").ID("querybuilding").Dot("UUIDExternalIDGenerator").Valuesln())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("logQueryBuildingError logs errs that may occur during query construction. Such errors should be few and far between,"),
		jen.Line(),
		jen.Func().Comment("as the generally only occur with type discrepancies or other misuses of SQL. An alert should be set up for any log").Comment("entries with the given name, and those alerts should be investigated quickly.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("logQueryBuildingError").Params(jen.ID("span").ID("tracing").Dot("Span"), jen.ID("err").ID("error")).Body(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("QueryErrorKey"),
					jen.ID("true"),
				),
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building query"),
				),
			)),
		jen.Line(),
	)

	return code
}
