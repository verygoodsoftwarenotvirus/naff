package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func clientDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("database").Dot(
		"Database",
	).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("Client").Struct(jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("querier").ID("database").Dot(
		"Database",
	), jen.ID("debug").ID("bool"), jen.ID("logger").ID("logging").Dot(
		"Logger",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// Migrate is a simple wrapper around the core querier Migrate call").Params(jen.ID("c").Op("*").ID("Client")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"Migrate",
		).Call(jen.ID("ctx")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// IsReady is a simple wrapper around the core querier IsReady call").Params(jen.ID("c").Op("*").ID("Client")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"IsReady",
		).Call(jen.ID("ctx")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideDatabaseClient provides a new Database client").ID("ProvideDatabaseClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("querier").ID("database").Dot(
		"Database",
	), jen.ID("debug").ID("bool"), jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("database").Dot(
		"Database",
	), jen.ID("error")).Block(
		jen.ID("c").Op(":=").Op("&").ID("Client").Valuesln(jen.ID("db").Op(":").ID("db"), jen.ID("querier").Op(":").ID("querier"), jen.ID("debug").Op(":").ID("debug"), jen.ID("logger").Op(":").ID("logger").Dot(
			"WithName",
		).Call(jen.Lit("db_client"))),
		jen.If(jen.ID("debug")).Block(
			jen.ID("c").Dot(
				"logger",
			).Dot(
				"SetLevel",
			).Call(jen.ID("logging").Dot(
				"DebugLevel",
			)),
		),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"Debug",
		).Call(jen.Lit("migrating querier")),
		jen.If(jen.ID("err").Op(":=").ID("c").Dot(
			"querier",
		).Dot(
			"Migrate",
		).Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"Debug",
		).Call(jen.Lit("querier migrated!")),
		jen.Return().List(jen.ID("c"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachUserIDToSpan provides a consistent way to attach a user's ID to a span").ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachFilterToSpan provides a consistent way to attach a filter's info to a span").ID("attachFilterToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	)).Block(
		jen.If(jen.ID("filter").Op("!=").ID("nil").Op("&&").ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("filter_page"), jen.Qual("strconv", "FormatUint").Call(jen.ID("filter").Dot(
				"QueryPage",
			).Call(), jen.Lit(10))), jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("filter_limit"), jen.Qual("strconv", "FormatUint").Call(jen.ID("filter").Dot(
				"Limit",
			), jen.Lit(10)))),
		),
	),

		jen.Line(),
	)
	return ret
}
