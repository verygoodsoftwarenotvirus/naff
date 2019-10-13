package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideConfigServerSettings is an obligatory function that").Comment("// we're required to have because wire doesn't do it for us.").ID("ProvideConfigServerSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("ServerSettings")).Block(
		jen.Return().ID("c").Dot(
			"Server",
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideConfigAuthSettings is an obligatory function that").Comment("// we're required to have because wire doesn't do it for us.").ID("ProvideConfigAuthSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("AuthSettings")).Block(
		jen.Return().ID("c").Dot(
			"Auth",
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideConfigDatabaseSettings is an obligatory function that").Comment("//  we're required to have because wire doesn't do it for us.").ID("ProvideConfigDatabaseSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("DatabaseSettings")).Block(
		jen.Return().ID("c").Dot(
			"Database",
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideConfigFrontendSettings is an obligatory function that").Comment("//  we're required to have because wire doesn't do it for us.").ID("ProvideConfigFrontendSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("FrontendSettings")).Block(
		jen.Return().ID("c").Dot(
			"Frontend",
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideConfigServerSettings"), jen.ID("ProvideConfigAuthSettings"), jen.ID("ProvideConfigDatabaseSettings"), jen.ID("ProvideConfigFrontendSettings")),

		jen.Line(),
	)
	return ret
}
