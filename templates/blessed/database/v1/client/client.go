package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func clientDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").ID("database").Dot("Database").Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Line()
	ret.Comment(`		NOTE: the primary purpose of this client is to allow convenient
		wrapping of actual query execution.
`)
	ret.Line()

	ret.Add(
		jen.Comment("Client is a wrapper around a database querier. Client is where all"),
		jen.Line(),
		jen.Comment("logging and trace propagation should happen, the querier is where"),
		jen.Line(),
		jen.Comment("the actual database querying is performed."),
		jen.Line(),
		jen.Type().ID("Client").Struct(jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("querier").ID("database").Dot("Database"),
			jen.ID("debug").ID("bool"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Migrate is a simple wrapper around the core querier Migrate call"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
			jen.Return().ID("c").Dot("querier").Dot("Migrate").Call(jen.ID("ctx")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsReady is a simple wrapper around the core querier IsReady call"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			jen.Return().ID("c").Dot("querier").Dot("IsReady").Call(jen.ID("ctx")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideDatabaseClient provides a new Database client"),
		jen.Line(),
		jen.Func().ID("ProvideDatabaseClient").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("db").Op("*").Qual("database/sql", "DB"),
			jen.ID("querier").ID("database").Dot("Database"),
			jen.ID("debug").ID("bool"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
		).Params(jen.ID("database").Dot("Database"), jen.ID("error")).Block(
			jen.ID("c").Op(":=").Op("&").ID("Client").Valuesln(
				jen.ID("db").Op(":").ID("db"),
				jen.ID("querier").Op(":").ID("querier"),
				jen.ID("debug").Op(":").ID("debug"),
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.Lit("db_client")),
			),
			jen.Line(),
			jen.If(jen.ID("debug")).Block(
				jen.ID("c").Dot("logger").Dot("SetLevel").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "DebugLevel")),
			),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("migrating querier")),
			jen.If(jen.ID("err").Op(":=").ID("c").Dot("querier").Dot("Migrate").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("querier migrated!")),
			jen.Line(),
			jen.Return().List(jen.ID("c"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachUserIDToSpan provides a consistent way to attach a user's ID to a span"),
		jen.Line(),
		jen.Func().ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Callln(
					jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10))),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachFilterToSpan provides a consistent way to attach a filter's info to a span"),
		jen.Line(),
		jen.Func().ID("attachFilterToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Block(
			jen.If(jen.ID("filter").Op("!=").ID("nil").Op("&&").ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Callln(
					jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("filter_page"), jen.Qual("strconv", "FormatUint").Call(jen.ID("filter").Dot("QueryPage").Call(), jen.Lit(10))),
					jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("filter_limit"), jen.Qual("strconv", "FormatUint").Call(jen.ID("filter").Dot("Limit"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)
	return ret
}