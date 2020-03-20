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
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database").Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.Nil()),
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

	ret.Add(buildMigrate()...)
	ret.Add(buildIsReady()...)
	ret.Add(buildProvideDatabaseClient(pkg)...)

	return ret
}

func buildMigrate() []jen.Code {
	funcName := "Migrate"

	lines := []jen.Code{
		jen.Commentf("%s is a simple wrapper around the core querier %s call", funcName, funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID(funcName).Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("error")).Block(
			append(
				utils.StartSpan(true, funcName),
				jen.Return().ID("c").Dot("querier").Dot(funcName).Call(utils.CtxVar()),
			)...,
		),
		jen.Line(),
	}

	return lines
}

func buildIsReady() []jen.Code {
	funcName := "IsReady"

	lines := []jen.Code{
		jen.Commentf("%s is a simple wrapper around the core querier %s call", funcName, funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID(funcName).Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			append(
				utils.StartSpan(true, funcName),
				jen.Return().ID("c").Dot("querier").Dot(funcName).Call(utils.CtxVar()),
			)...,
		),
		jen.Line(),
	}

	return lines
}

func buildProvideDatabaseClient(proj *models.Project) []jen.Code {
	funcName := "ProvideDatabaseClient"

	lines := []jen.Code{
		jen.Commentf("%s provides a new Database client", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Paramsln(
			utils.CtxParam(),
			jen.ID("db").Op("*").Qual("database/sql", "DB"),
			jen.ID("querier").Qual(filepath.Join(proj.OutputPath, "database/v1"), "Database"),
			jen.ID("debug").ID("bool"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
		).Params(jen.Qual(filepath.Join(proj.OutputPath, "database/v1"), "Database"), jen.ID("error")).Block(
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
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("querier migrated!")),
			jen.Line(),
			jen.Return().List(jen.ID("c"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
