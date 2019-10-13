package v1

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day").ID("ProvideReporter").Params(jen.ID("n").Op("*").ID("newsman").Dot(
		"Newsman",
	)).Params(jen.ID("newsman").Dot(
		"Reporter",
	)).Block(
		jen.Return().ID("n"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// BuildServer builds a server").ID("BuildServer").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("config").Dot(
		"ServerConfig",
	), jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("database").ID("database").Dot(
		"Database",
	)).Params(jen.Op("*").ID("server").Dot(
		"Server",
	), jen.ID("error")).Block(
		jen.ID("wire").Dot(
			"Build",
		).Call(jen.ID("config").Dot(
			"Providers",
		), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth", "Providers"), jen.ID("server").Dot(
			"Providers",
		), jen.ID("encoding").Dot(
			"Providers",
		), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "Providers"), jen.ID("metrics").Dot(
			"Providers",
		), jen.ID("newsman").Dot(
			"NewNewsman",
		), jen.ID("ProvideReporter"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth", "Providers"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users", "Providers"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items", "Providers"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/frontend", "Providers"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks", "Providers"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients", "Providers")),
		jen.Return().List(jen.ID("nil"), jen.ID("nil")),
	),

		jen.Line(),
	)
	return ret
}
