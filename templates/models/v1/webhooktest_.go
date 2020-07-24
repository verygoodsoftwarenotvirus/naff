package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)

	code.Add(
		jen.Func().ID("TestWebhook_Update").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID(utils.BuildFakeVarName("Input")).Assign().AddressOf().ID("WebhookUpdateInput").Valuesln(
					jen.ID("Name").MapAssign().Lit("whatever"),
					jen.ID("ContentType").MapAssign().Lit("application/xml"),
					jen.ID("URL").MapAssign().Lit("https://blah.verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Events").MapAssign().Index().String().Values(jen.Lit("more_things")),
					jen.ID("DataTypes").MapAssign().Index().String().Values(jen.Lit("new_stuff")),
					jen.ID("Topics").MapAssign().Index().String().Values(jen.Lit("blah-blah")),
				),
				jen.Line(),
				jen.ID("actual").Assign().AddressOf().ID("Webhook").Valuesln(
					jen.ID("Name").MapAssign().Lit("something_else"),
					jen.ID("ContentType").MapAssign().Lit("application/json"),
					jen.ID("URL").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPost"),
					jen.ID("Events").MapAssign().Index().String().Values(jen.Lit("things")),
					jen.ID("DataTypes").MapAssign().Index().String().Values(jen.Lit("stuff")),
					jen.ID("Topics").MapAssign().Index().String().Values(jen.Lit("blah")),
				),
				jen.ID("expected").Assign().AddressOf().ID("Webhook").Valuesln(
					jen.ID("Name").MapAssign().ID("exampleInput").Dot("Name"),
					jen.ID("ContentType").MapAssign().Lit("application/xml"),
					jen.ID("URL").MapAssign().Lit("https://blah.verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Events").MapAssign().Index().String().Values(jen.Lit("more_things")),
					jen.ID("DataTypes").MapAssign().Index().String().Values(jen.Lit("new_stuff")),
					jen.ID("Topics").MapAssign().Index().String().Values(jen.Lit("blah-blah")),
				),
				jen.Line(),
				jen.ID("actual").Dot("Update").Call(jen.ID(utils.BuildFakeVarName("Input"))),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhook_ToListener").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("w").Assign().AddressOf().ID("Webhook").Values(),
				jen.ID("w").Dot("ToListener").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildErrorLogFunc").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("w").Assign().AddressOf().ID("Webhook").Values(),
				jen.ID("actual").Assign().ID("buildErrorLogFunc").Call(jen.ID("w"), jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("actual").Call(constants.ObligatoryError()),
			),
		),
		jen.Line(),
	)

	return code
}
