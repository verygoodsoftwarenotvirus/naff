package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

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

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				buildProviderSet()...,
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideNamespace provides a namespace"),
		jen.Line(),
		jen.Func().ID("ProvideNamespace").Params().Params(jen.Qual(proj.InternalMetricsV1Package(), "Namespace")).Block(
			jen.Return().ID("serverNamespace"),
		),
		jen.Line(),
	)

	// if proj.EnableNewsman {
	ret.Add(
		jen.Comment("ProvideNewsmanTypeNameManipulationFunc provides an WebhookIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideNewsmanTypeNameManipulationFunc").Params().Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "TypeNameManipulationFunc")).Block(
			jen.Return().Func().Params(jen.ID("s").String()).Params(jen.String()).Block(
				jen.Return().ID("s"),
			),
		),
		jen.Line(),
	)
	// }

	return ret
}
