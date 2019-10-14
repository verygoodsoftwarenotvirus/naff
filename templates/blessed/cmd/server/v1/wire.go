package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.HeaderComment("+build wireinject")

	utils.AddImports(ret)

	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "httpserver")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth", "authservice")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/frontend", "frontendservice")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items", "itemsservice")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients", "oauth2clientsservice")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users", "usersservice")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks", "webhooksservice")

	ret.Add(
		jen.Comment("ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day"),
		jen.Line(),
		jen.Func().ID("ProvideReporter").Params(jen.ID("n").Op("*").Qual(newsmanImp, "Newsman")).Params(jen.Qual(newsmanImp, "Reporter")).Block(
			jen.Return().ID("n"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildServer builds a server"),
		jen.Line(),
		jen.Func().ID("BuildServer").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("cfg").Op("*").Qual(internalConfigImp, "ServerConfig"),
			jen.ID("logger").Qual(loggingImp, "Logger"),
			jen.ID("database").Qual(databaseClientImp, "Database")).
			Params(jen.Op("*").Qual(serverImp, "Server"), jen.ID("error")).Block(
			jen.ID("wire").Dot(
				"Build",
			).Callln(
				jen.Qual(internalConfigImp, "Providers"),
				jen.Qual(internalAuthImp, "Providers"),
				jen.Comment("server things"),
				jen.Qual(serverImp, "Providers"),
				jen.Qual(internalEncodingImp, "Providers"),
				jen.Qual(httpServerImp, "Providers"),
				jen.Comment("metrics"),
				jen.Qual(internalMetricsImp, "Providers"),
				jen.Comment("external libs"),
				jen.Qual(newsmanImp, "NewNewsman"),
				jen.ID("ProvideReporter"),
				jen.Comment("services"),
				jen.Qual(authServiceImp, "Providers"),
				jen.Qual(usersServiceImp, "Providers"),
				jen.Qual(itemsServiceImp, "Providers"),
				jen.Qual(frontendServiceImp, "Providers"),
				jen.Qual(webhooksServiceImp, "Providers"),
				jen.Qual(oauth2ClientsServiceImp, "Providers"),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return ret
}
