package dbclient

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func clientDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	ret.Add(jen.Null().Var().ID("_").ID("database").Dot("Database").Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("Client").Struct(
		jen.ID("db").Op("*").Qual("database/sql", "DB"),
		jen.ID("querier").ID("database").Dot(
			"Database",
		),
		jen.ID("debug").ID("bool"),
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
	),
	)
	return ret
}
