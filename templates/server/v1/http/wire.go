package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("httpserver")

	utils.AddImports(proj, code)

	buildProviderSet := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("paramFetcherProviders"),
			jen.ID("ProvideServer"),
			jen.ID("ProvideNamespace"),
		}

		// if proj.EnableNewsman {
		lines = append(lines, jen.ID("ProvideNewsmanTypeNameManipulationFunc"))
		// }

		return lines
	}

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				buildProviderSet()...,
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideNamespace provides a namespace."),
		jen.Line(),
		jen.Func().ID("ProvideNamespace").Params().Params(jen.Qual(proj.InternalMetricsV1Package(), "Namespace")).Block(
			jen.Return().ID("serverNamespace"),
		),
		jen.Line(),
	)

	// if proj.EnableNewsman {
	code.Add(
		jen.Comment("ProvideNewsmanTypeNameManipulationFunc provides an WebhookIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideNewsmanTypeNameManipulationFunc").Params().Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "TypeNameManipulationFunc")).Block(
			jen.Return().Func().Params(jen.ID("s").String()).Params(jen.String()).Block(
				jen.Return().ID("s"),
			),
		),
		jen.Line(),
	)
	// }

	return code
}
