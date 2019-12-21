package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Line(),
		jen.Comment("BEGIN it'd be neat if wire could do this for me one day."),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideConfigServerSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigServerSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("ServerSettings")).Block(
			jen.Return().ID("c").Dot(
				"Server",
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideConfigAuthSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigAuthSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("AuthSettings")).Block(
			jen.Return().ID("c").Dot(
				"Auth",
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideConfigDatabaseSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigDatabaseSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("DatabaseSettings")).Block(
			jen.Return().ID("c").Dot(
				"Database",
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideConfigFrontendSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigFrontendSettings").Params(jen.ID("c").Op("*").ID("ServerConfig")).Params(jen.ID("FrontendSettings")).Block(
			jen.Return().ID("c").Dot(
				"Frontend",
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Line(),
		jen.Comment("END it'd be neat if wire could do this for me one day."),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents this package's offering to the dependency manager"),
			jen.ID("Providers").Op("=").Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID("ProvideConfigServerSettings"),
				jen.ID("ProvideConfigAuthSettings"),
				jen.ID("ProvideConfigDatabaseSettings"),
				jen.ID("ProvideConfigFrontendSettings"),
			),
		),
		jen.Line(),
	)
	return ret
}
