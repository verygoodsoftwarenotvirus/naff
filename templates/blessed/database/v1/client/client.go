package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database").Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
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
		jen.Type().ID("Client").Struct(jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("querier").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"),
			jen.ID("debug").ID("bool"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Migrate is a simple wrapper around the core querier Migrate call"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("Migrate").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("error")).Block(
			jen.Return().ID("c").Dot("querier").Dot("Migrate").Call(utils.CtxVar()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsReady is a simple wrapper around the core querier IsReady call"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("IsReady").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			jen.Return().ID("c").Dot("querier").Dot("IsReady").Call(utils.CtxVar()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideDatabaseClient provides a new Database client"),
		jen.Line(),
		jen.Func().ID("ProvideDatabaseClient").Paramsln(
			utils.CtxParam(),
			jen.ID("db").Op("*").Qual("database/sql", "DB"),
			jen.ID("querier").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"),
			jen.ID("debug").ID("bool"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
		).Params(jen.Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"), jen.ID("error")).Block(
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
			jen.If(jen.Err().Op(":=").ID("c").Dot("querier").Dot("Migrate").Call(utils.CtxVar()), jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Err()),
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
		jen.Func().ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual(utils.TracingLibrary, "Span"), jen.ID("userID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Callln(
					jen.Qual(utils.TracingLibrary, "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10))),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachFilterToSpan provides a consistent way to attach a filter's info to a span"),
		jen.Line(),
		jen.Func().ID("attachFilterToSpan").Params(jen.ID("span").Op("*").Qual(utils.TracingLibrary, "Span"), jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Block(
			jen.If(jen.ID("filter").Op("!=").ID("nil").Op("&&").ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Callln(
					jen.Qual(utils.TracingLibrary, "StringAttribute").Call(jen.Lit("filter_page"), jen.Qual("strconv", "FormatUint").Call(jen.ID("filter").Dot("QueryPage").Call(), jen.Lit(10))),
					jen.Qual(utils.TracingLibrary, "StringAttribute").Call(jen.Lit("filter_limit"), jen.Qual("strconv", "FormatUint").Call(jen.ID("filter").Dot("Limit"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)
	return ret
}
