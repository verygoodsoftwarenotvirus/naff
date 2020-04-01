package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(proj.OutputPath, "database/v1"), "Database").Equals().Parens(jen.Op("*").ID("Client")).Call(jen.Nil()),
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
		jen.Type().ID("Client").Struct(jen.ID("db").ParamPointer().Qual("database/sql", "DB"), jen.ID("querier").Qual(filepath.Join(proj.OutputPath, "database/v1"), "Database"),
			jen.ID("debug").ID("bool"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")),
		jen.Line(),
	)

	ret.Add(buildMigrate(proj)...)
	ret.Add(buildIsReady(proj)...)
	ret.Add(buildProvideDatabaseClient(proj)...)

	return ret
}

func buildMigrate(proj *models.Project) []jen.Code {
	funcName := "Migrate"

	lines := []jen.Code{
		jen.Commentf("%s is a simple wrapper around the core querier %s call", funcName, funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID(funcName).Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("error")).Block(
			utils.StartSpan(proj, true, funcName),
			jen.Return().ID("c").Dot("querier").Dot(funcName).Call(utils.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}

func buildIsReady(proj *models.Project) []jen.Code {
	funcName := "IsReady"

	lines := []jen.Code{
		jen.Commentf("%s is a simple wrapper around the core querier %s call", funcName, funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID(funcName).Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			utils.StartSpan(proj, true, funcName),
			jen.Return().ID("c").Dot("querier").Dot(funcName).Call(utils.CtxVar()),
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
			jen.ID("db").ParamPointer().Qual("database/sql", "DB"),
			jen.ID("querier").Qual(filepath.Join(proj.OutputPath, "database/v1"), "Database"),
			jen.ID("debug").ID("bool"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
		).Params(jen.Qual(filepath.Join(proj.OutputPath, "database/v1"), "Database"), jen.ID("error")).Block(
			jen.ID("c").Assign().VarPointer().ID("Client").Valuesln(
				jen.ID("db").MapAssign().ID("db"),
				jen.ID("querier").MapAssign().ID("querier"),
				jen.ID("debug").MapAssign().ID("debug"),
				jen.ID("logger").MapAssign().ID("logger").Dot("WithName").Call(jen.Lit("db_client")),
			),
			jen.Line(),
			jen.If(jen.ID("debug")).Block(
				jen.ID("c").Dot("logger").Dot("SetLevel").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "DebugLevel")),
			),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("migrating querier")),
			jen.If(jen.Err().Assign().ID("c").Dot("querier").Dot("Migrate").Call(utils.CtxVar()), jen.Err().DoesNotEqual().ID("nil")).Block(
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
