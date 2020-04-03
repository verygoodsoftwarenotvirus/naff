package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestWebhook_Update").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("exampleInput").Assign().VarPointer().ID("WebhookUpdateInput").Valuesln(
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("ContentType").MapAssign().Lit("application/xml"),
					jen.ID("URL").MapAssign().Lit("https://blah.verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Events").MapAssign().Index().ID("string").Values(jen.Lit("more_things")),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(jen.Lit("new_stuff")),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.Lit("blah-blah")),
				),
				jen.Line(),
				jen.ID("actual").Assign().VarPointer().ID("Webhook").Valuesln(
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("ContentType").MapAssign().Lit("application/json"),
					jen.ID("URL").MapAssign().Lit("https://verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPost"),
					jen.ID("Events").MapAssign().Index().ID("string").Values(jen.Lit("things")),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(jen.Lit("stuff")),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.Lit("blah")),
				),
				jen.ID("expected").Assign().VarPointer().ID("Webhook").Valuesln(
					jen.ID("Name").MapAssign().ID("exampleInput").Dot("Name"),
					jen.ID("ContentType").MapAssign().Lit("application/xml"),
					jen.ID("URL").MapAssign().Lit("https://blah.verygoodsoftwarenotvirus.ru"),
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Events").MapAssign().Index().ID("string").Values(jen.Lit("more_things")),
					jen.ID("DataTypes").MapAssign().Index().ID("string").Values(jen.Lit("new_stuff")),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.Lit("blah-blah")),
				),
				jen.Line(),
				jen.ID("actual").Dot("Update").Call(jen.ID("exampleInput")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestWebhook_ToListener").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("w").Assign().VarPointer().ID("Webhook").Values(),
				jen.ID("w").Dot("ToListener").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildErrorLogFunc").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("w").Assign().VarPointer().ID("Webhook").Values(),
				jen.ID("actual").Assign().ID("buildErrorLogFunc").Call(jen.ID("w"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("actual").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			),
		),
		jen.Line(),
	)
	return ret
}
