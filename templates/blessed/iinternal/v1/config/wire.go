package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("config")

	utils.AddImports(proj, code)

	searchEnabled := false
	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			searchEnabled = true
		}
	}

	code.Add(
		jen.Line(),
		jen.Comment("BEGIN it'd be neat if wire could do this for me one day."),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideConfigServerSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigServerSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("ServerSettings")).Block(
			jen.Return().ID("c").Dot(
				"Server",
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideConfigAuthSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigAuthSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("AuthSettings")).Block(
			jen.Return().ID("c").Dot(
				"Auth",
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideConfigDatabaseSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigDatabaseSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("DatabaseSettings")).Block(
			jen.Return().ID("c").Dot(
				"Database",
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideConfigFrontendSettings is an obligatory function that"),
		jen.Line(),
		jen.Comment("we're required to have because wire doesn't do it for us."),
		jen.Line(),
		jen.Func().ID("ProvideConfigFrontendSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("FrontendSettings")).Block(
			jen.Return().ID("c").Dot(
				"Frontend",
			),
		),
		jen.Line(),
	)

	if searchEnabled {
		code.Add(
			jen.Comment("ProvideSearchSettings is an obligatory function that"),
			jen.Line(),
			jen.Comment("we're required to have because wire doesn't do it for us."),
			jen.Line(),
			jen.Func().ID("ProvideSearchSettings").Params(jen.ID("c").PointerTo().ID("ServerConfig")).Params(jen.ID("SearchSettings")).Block(
				jen.Return().ID("c").Dot(
					"Search",
				),
			),
			jen.Line(),
		)
	}

	code.Add(
		jen.Line(),
		jen.Comment("END it'd be neat if wire could do this for me one day."),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents this package's offering to the dependency manager."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideConfigServerSettings"),
				jen.ID("ProvideConfigAuthSettings"),
				jen.ID("ProvideConfigDatabaseSettings"),
				jen.ID("ProvideConfigFrontendSettings"),
				func() jen.Code {
					if searchEnabled {
						return jen.ID("ProvideSearchSettings")
					}
					return jen.Null()
				}(),
				jen.ID("ProvideSessionManager"),
			),
		),
		jen.Line(),
	)

	return code
}
