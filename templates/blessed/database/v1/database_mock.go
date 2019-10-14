package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

const (
	mockImp       = "github.com/stretchr/testify/mock"
	mockModelsImp = "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock"
)

func databaseMockDotGo() *jen.File {
	ret := jen.NewFile("database")

	ret.ImportAlias(mockModelsImp, "mockmodels")
	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").ID("Database").Op("=").Parens(jen.Op("*").ID("MockDatabase")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildMockDatabase builds a mock database"),
		jen.Line(),
		jen.Func().ID("BuildMockDatabase").Params().Params(jen.Op("*").ID("MockDatabase")).Block(
			jen.Return().Op("&").ID("MockDatabase").Valuesln(
				jen.ID("ItemDataManager").Op(":").Op("&").Qual(mockModelsImp, "ItemDataManager").Values(),
				jen.ID("UserDataManager").Op(":").Op("&").Qual(mockModelsImp, "UserDataManager").Values(),
				jen.ID("OAuth2ClientDataManager").Op(":").Op("&").Qual(mockModelsImp, "OAuth2ClientDataManager").Values(),
				jen.ID("WebhookDataManager").Op(":").Op("&").Qual(mockModelsImp, "WebhookDataManager").Values(),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("MockDatabase is our mock database structure"),
		jen.Line(),
		jen.Type().ID("MockDatabase").Struct(
			jen.Qual(mockImp, "Mock"),
			jen.Line(),
			jen.Op("*").Qual(mockModelsImp, "ItemDataManager"),
			jen.Op("*").Qual(mockModelsImp, "UserDataManager"),
			jen.Op("*").Qual(mockModelsImp, "OAuth2ClientDataManager"),
			jen.Op("*").Qual(mockModelsImp, "WebhookDataManager"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Migrate satisfies the database.Database interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockDatabase")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsReady satisfies the database.Database interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockDatabase")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().ID("args").Dot("Bool").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
