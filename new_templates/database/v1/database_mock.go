package database

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func databaseMockDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("_").ID("Database").Op("=").Parens(jen.Op("*").ID("MockDatabase")).Call(jen.ID("nil")))
	ret.Add(jen.Func())
	ret.Add(jen.Null().Type().ID("MockDatabase").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
		jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataManager"),
		jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager"),
		jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "OAuth2ClientDataManager"),
		jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager"),
	),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
