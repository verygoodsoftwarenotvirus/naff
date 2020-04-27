package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.DatabaseV1Package(), "Database").Equals().Parens(jen.PointerTo().ID("Client")).Call(jen.Nil()),
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
		jen.Type().ID("Client").Struct(jen.ID("db").PointerTo().Qual("database/sql", "DB"), jen.ID("querier").Qual(proj.DatabaseV1Package(), "Database"),
			jen.ID("debug").Bool(), jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")),
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
		jen.Commentf("%s is a simple wrapper around the core querier %s call.", funcName, funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID(funcName).Params(constants.CtxParam()).Params(jen.Error()).Block(
			utils.StartSpan(proj, true, funcName),
			jen.Return().ID("c").Dot("querier").Dot(funcName).Call(constants.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}

func buildIsReady(proj *models.Project) []jen.Code {
	funcName := "IsReady"

	lines := []jen.Code{
		jen.Commentf("%s is a simple wrapper around the core querier %s call.", funcName, funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID(funcName).Params(constants.CtxParam()).Params(jen.ID("ready").Bool()).Block(
			utils.StartSpan(proj, true, funcName),
			jen.Return().ID("c").Dot("querier").Dot(funcName).Call(constants.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideDatabaseClient(proj *models.Project) []jen.Code {
	funcName := "ProvideDatabaseClient"

	lines := []jen.Code{
		jen.Commentf("%s provides a new Database client.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Paramsln(
			constants.CtxParam(),
			jen.ID("db").PointerTo().Qual("database/sql", "DB"),
			jen.ID("querier").Qual(proj.DatabaseV1Package(), "Database"),
			jen.ID("debug").Bool(),
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
		).Params(jen.Qual(proj.DatabaseV1Package(), "Database"), jen.Error()).Block(
			jen.ID("c").Assign().AddressOf().ID("Client").Valuesln(
				jen.ID("db").MapAssign().ID("db"),
				jen.ID("querier").MapAssign().ID("querier"),
				jen.ID("debug").MapAssign().ID("debug"),
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.Lit("db_client")),
			),
			jen.Line(),
			jen.If(jen.ID("debug")).Block(
				jen.ID("c").Dot(constants.LoggerVarName).Dot("SetLevel").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "DebugLevel")),
			),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("migrating querier")),
			jen.If(jen.Err().Assign().ID("c").Dot("querier").Dot("Migrate").Call(constants.CtxVar()), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("querier migrated!")),
			jen.Line(),
			jen.Return().List(jen.ID("c"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
