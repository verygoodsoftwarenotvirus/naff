package v1

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func databaseMockDotGo() *jen.File {
	ret := jen.NewFile("database")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("Database").Op("=").Parens(jen.Op("*").ID("MockDatabase")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// BuildMockDatabase builds a mock database").ID("BuildMockDatabase").Params().Params(jen.Op("*").ID("MockDatabase")).Block(
		jen.Return().Op("&").ID("MockDatabase").Valuesln(jen.ID("ItemDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataManager").Valuesln(), jen.ID("UserDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager").Valuesln(), jen.ID("OAuth2ClientDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "OAuth2ClientDataManager").Valuesln(), jen.ID("WebhookDataManager").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln()),
	),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("MockDatabase").Struct(jen.ID("mock").Dot(
		"Mock",
	), jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataManager"), jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager"), jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "OAuth2ClientDataManager"), jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// Migrate satisfies the database.Database interface").Params(jen.ID("m").Op("*").ID("MockDatabase")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx")),
		jen.Return().ID("args").Dot(
			"Error",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// IsReady satisfies the database.Database interface").Params(jen.ID("m").Op("*").ID("MockDatabase")).ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("ready").ID("bool")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx")),
		jen.Return().ID("args").Dot(
			"Bool",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	return ret
}
