package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksServiceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildBuildTestService(proj)...)
	code.Add(buildTestProvideWebhooksService(proj)...)

	return code
}

func buildBuildTestService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params().Params(jen.PointerTo().ID("Service")).Body(
			jen.Return().AddressOf().ID("Service").Valuesln(
				jen.ID(constants.LoggerVarName).MapAssign().Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
				jen.ID("webhookCounter").MapAssign().AddressOf().Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
				jen.ID("webhookDataManager").MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), "WebhookDataManager").Values(),
				jen.ID("userIDFetcher").MapAssign().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
				jen.ID("webhookIDFetcher").MapAssign().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
				jen.ID("encoderDecoder").MapAssign().AddressOf().Qual(proj.EncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("eventManager").MapAssign().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideWebhooksService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideWebhooksService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.MetricsPackage(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter"), jen.Error()).Body(
					jen.Return().List(jen.AddressOf().Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideWebhooksService").Callln(
					jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
					jen.AddressOf().Qual(proj.TypesPackage("mock"), "WebhookDataManager").Values(),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.EncodingPackage("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
				),
				jen.Line(),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error providing counter",
				jen.Var().ID("ucp").Qual(proj.MetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.MetricsPackage(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.MetricsPackage(), "UnitCounter"), jen.Error()).Body(
					jen.Return().List(jen.Nil(), constants.ObligatoryError()),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideWebhooksService").Callln(
					jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
					jen.AddressOf().Qual(proj.TypesPackage("mock"), "WebhookDataManager").Values(),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.EncodingPackage("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
