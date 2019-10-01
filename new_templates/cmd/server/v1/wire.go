package cmdv1server

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	newsman := "gitlab.com/verygoodsoftwarenotvirus/newsman"
	nill := jen.ID("nil")

	ret := jen.NewFile("main")
	ret.HeaderComment("+build wireinject")

	utils.AddImports(ret)
	ret.ImportNames(map[string]string{
		"gitlab.com/verygoodsoftwarenotvirus/todo/database/v1":               "database",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/config/v1":        "config",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1":      "encoding",
		"gitlab.com/verygoodsoftwarenotvirus/todo/internal/logging/v1":       "logging",
		"gitlab.com/verygoodsoftwarenotvirus/todo/server/v1":                 "server",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth":          "auth",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/frontend":      "frontend",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items":         "items",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients": "oauth2clients",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users":         "users",
		"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks":      "webhooks",
		newsman: "newsman",
	})
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "httpserver")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1", "libauth")
	ret.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1", "metricsProvider")

	ret.Add(
		jen.Comment("ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day"),
		jen.Line(),
		jen.Func().ID("ProvideReporter").Params(jen.ID("n").Op("*").Qual(newsman, "Newsman")).Params(jen.Qual(newsman, "Reporter")).Block(
			jen.Return().ID("n"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildServer builds a server"),
		jen.Line(),
		jen.Func().ID("BuildServer").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("cfg").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config/v1", "ServerConfig"),
			jen.ID("logger").Qual(utils.LoggingPkg, "Logger"),
			jen.ID("database").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Database"),
		).Params(
			jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1", "Server"),
			jen.ID("error"),
		).Block(
			jen.Qual("github.com/google/wire", "Build").Callln(
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config/v1", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1", "Providers"),
				jen.Line().Add(jen.Comment("Server things")),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "Providers"),
				jen.Line().Add(jen.Comment("Server things")),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1", "Providers"),
				jen.Line().Add(jen.Comment("external libs")),
				jen.ID("ProvideReporter"),
				jen.Qual(newsman, "NewNewsman"),
				jen.Line().Add(jen.Comment("services")),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/frontend", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks", "Providers"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients", "Providers"),
			),
			jen.Line(),
			jen.Return().List(nill, nill),
		),
	)
	return ret
}
