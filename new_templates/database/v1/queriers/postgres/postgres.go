package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func postgresDotGo() *jen.File {
	ret := jen.NewFile("postgres")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("loggerName").Op("=").Lit("postgres").Var().ID("postgresDriverName").Op("=").Lit("wrapped-postgres-driver").Var().ID("CountQuery").Op("=").Lit("COUNT(id)").Var().ID("CurrentUnixTimeQuery").Op("=").Lit("extract(epoch FROM NOW())"),
	)
	ret.Add(jen.Func().ID("init").Params().Block(
		jen.ID("driver").Op(":=").ID("ocsql").Dot(
			"Wrap",
		).Call(jen.Op("&").Qual("github.com/lib/pq", "Driver").Valuesln(), jen.ID("ocsql").Dot(
			"WithQuery",
		).Call(jen.ID("true")), jen.ID("ocsql").Dot(
			"WithAllowRoot",
		).Call(jen.ID("false")), jen.ID("ocsql").Dot(
			"WithRowsNext",
		).Call(jen.ID("true")), jen.ID("ocsql").Dot(
			"WithRowsClose",
		).Call(jen.ID("true")), jen.ID("ocsql").Dot(
			"WithQueryParams",
		).Call(jen.ID("true"))),
		jen.Qual("database/sql", "Register").Call(jen.ID("postgresDriverName"), jen.ID("driver")),
	),
	)
	ret.Add(jen.Null().Type().ID("Postgres").Struct(
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("db").Op("*").Qual("database/sql", "DB"),
		jen.ID("sqlBuilder").ID("squirrel").Dot(
			"StatementBuilderType",
		),
		jen.ID("migrateOnce").Qual("sync", "Once"),
		jen.ID("debug").ID("bool"),
	).Type().ID("ConnectionDetails").ID("string").Type().ID("Querier").Interface(jen.ID("ExecContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")), jen.ID("QueryContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")), jen.ID("QueryRowContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row"))),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
