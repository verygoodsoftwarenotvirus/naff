package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksServiceTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params().Params(jen.PointerTo().ID("Service")).Block(
			jen.Return().AddressOf().ID("Service").Valuesln(
				jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("webhookCounter").MapAssign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("webhookDatabase").MapAssign().AddressOf().Qual(proj.ModelsV1Package("mock"), "WebhookDataManager").Values(),
				jen.ID("userIDFetcher").MapAssign().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
				jen.ID("webhookIDFetcher").MapAssign().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
				jen.ID("encoderDecoder").MapAssign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("eventManager").MapAssign().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhooksService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"), jen.Error()).Block(
					jen.Return().List(jen.AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideWebhooksService").Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.AddressOf().Qual(proj.ModelsV1Package("mock"), "WebhookDataManager").Values(),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
				),
				jen.Line(),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error providing counter",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"), jen.Error()).Block(
					jen.Return().List(jen.Nil(), utils.ObligatoryError()),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideWebhooksService").Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.AddressOf().Qual(proj.ModelsV1Package("mock"), "WebhookDataManager").Values(),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	)

	return ret
}
