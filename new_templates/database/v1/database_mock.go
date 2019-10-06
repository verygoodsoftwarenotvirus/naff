package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func databaseMockDotGo() *jen.File {
	ret := jen.NewFile("database")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("Database").Op("=").Parens(jen.Op("*").ID("MockDatabase")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Func(),
	)
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
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
