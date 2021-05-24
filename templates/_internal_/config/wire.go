package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Line(),
		jen.Comment("BEGIN it'd be neat if wire could do this for me one day."),
		jen.Line(),
	)

	code.Add(buildProvideConfigServerSettings()...)

	code.Add(buildProvideConfigAuthSettings()...)

	code.Add(buildProvideConfigDatabaseSettings()...)

	code.Add(buildProvideConfigFrontendSettings()...)

	if proj.SearchEnabled() {
		code.Add(buildProvideSearchSettings()...)
	}

	code.Add(
		jen.Line(),
		jen.Comment("END it'd be neat if wire could do this for me one day."),
		jen.Line(),
	)

	code.Add(buildWireProviders(proj)...)

	return code
}

func buildProvideConfigServerSettings() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideConfigServerSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigServerSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("ServerSettings")).Body(
			jen.Return().ID("c").Dot(
				"Server",
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideConfigAuthSettings() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideConfigAuthSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigAuthSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("AuthSettings")).Body(
			jen.Return().ID("c").Dot(
				"Auth",
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideConfigDatabaseSettings() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideConfigDatabaseSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigDatabaseSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("DatabaseSettings")).Body(
			jen.Return().ID("c").Dot(
				"Database",
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideConfigFrontendSettings() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideConfigFrontendSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigFrontendSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("FrontendSettings")).Body(
			jen.Return().ID("c").Dot(
				"Frontend",
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideSearchSettings() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideSearchSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideSearchSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("SearchSettings")).Body(
			jen.Return().ID("c").Dot(
				"Search",
			),
		),
		jen.Line(),
	}

	return lines
}

func buildWireProviders(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers represents this package's offering to the dependency manager."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideConfigServerSettings"),
				jen.ID("ProvideConfigAuthSettings"),
				jen.ID("ProvideConfigDatabaseSettings"),
				jen.ID("ProvideConfigFrontendSettings"),
				func() jen.Code {
					if proj.SearchEnabled() {
						return jen.ID("ProvideSearchSettings")
					}
					return jen.Null()
				}(),
				jen.ID("ProvideSessionManager"),
			),
		),
		jen.Line(),
	}

	return lines
}
