package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookTestDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("TestWebhook_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("actual").Op(":=").Op("&").ID("Webhook").Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("ContentType").Op(":").Lit("application/json"),
					jen.ID("URL").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").Op(":").Qual("net/http", "MethodPost"),
					jen.ID("Events").Op(":").Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(jen.Lit("stuff")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.Lit("blah")),
				),
				jen.ID("expected").Op(":=").Op("&").ID("Webhook").Valuesln(
					jen.ID("Name").Op(":").Lit("new name"),
					jen.ID("ContentType").Op(":").Lit("application/xml"),
					jen.ID("URL").Op(":").Lit("https://blah.verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"),
					jen.ID("Events").Op(":").Index().ID("string").Values(jen.Lit("more_things")),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(jen.Lit("new_stuff")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.Lit("blah-blah")),
				),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").ID("WebhookUpdateInput").Valuesln(
					jen.ID("Name").Op(":").Lit("new name"),
					jen.ID("ContentType").Op(":").Lit("application/xml"),
					jen.ID("URL").Op(":").Lit("https://blah.verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"),
					jen.ID("Events").Op(":").Index().ID("string").Values(jen.Lit("more_things")),
					jen.ID("DataTypes").Op(":").Index().ID("string").Values(jen.Lit("new_stuff")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.Lit("blah-blah")),
				),
				jen.Line(),
				jen.ID("actual").Dot("Update").Call(jen.ID("exampleInput")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestWebhook_ToListener").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("w").Op(":=").Op("&").ID("Webhook").Values(),
				jen.ID("w").Dot("ToListener").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildErrorLogFunc").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("w").Op(":=").Op("&").ID("Webhook").Values(),
				jen.ID("actual").Op(":=").ID("buildErrorLogFunc").Call(jen.ID("w"), jen.ID("noop").Dot("ProvideNoopLogger").Call()),
				jen.ID("actual").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			)),
		),
		jen.Line(),
	)
	return ret
}
