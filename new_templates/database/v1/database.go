package database

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func databaseDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("_").ID("Scanner").Op("=").Parens(jen.Op("*").Qual("database/sql", "Row")).Call(jen.ID("nil")).Var().ID("_").ID("Querier").Op("=").Parens(jen.Op("*").Qual("database/sql", "DB")).Call(jen.ID("nil")).Var().ID("_").ID("Querier").Op("=").Parens(jen.Op("*").Qual("database/sql", "Tx")).Call(jen.ID("nil")))
	ret.Add(jen.Null().Type().ID("Scanner").Interface(jen.ID("Scan").Params(jen.ID("dest").Op("...").Interface()).Params(jen.ID("error"))).Type().ID("Querier").Interface(jen.ID("ExecContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")), jen.ID("QueryContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")), jen.ID("QueryRowContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row"))).Type().ID("ConnectionDetails").ID("string").Type().ID("Database").Interface(jen.ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")), jen.ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")), jen.ID("models").Dot(
		"ItemDataManager",
	), jen.ID("models").Dot(
		"UserDataManager",
	), jen.ID("models").Dot(
		"OAuth2ClientDataManager",
	), jen.ID("models").Dot(
		"WebhookDataManager",
	)),
	)
	return ret
}
