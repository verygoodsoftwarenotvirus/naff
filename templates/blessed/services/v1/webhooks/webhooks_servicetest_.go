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
			utils.BuildSubTest(
				"happy path",
				jen.ID("expectation").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.ID("uc").Assign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.Qual(utils.MockPkg, "Anything"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Paramsln(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"), jen.Error()).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.ID("dm").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "WebhookDataManager").Values(),
				jen.ID("dm").Dot("On").Call(jen.Lit("GetAllWebhooksCount"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideWebhooksService").Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("dm"),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
				),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error providing counter",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Paramsln(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String()).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error()).Block(
					jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideWebhooksService").Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.AddressOf().Qual(proj.ModelsV1Package("mock"), "WebhookDataManager").Values(),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil())),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error getting count",
				jen.ID("expectation").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.ID("uc").Assign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("IncrementBy"), jen.Qual(utils.MockPkg, "Anything"), jen.ID("expectation")).Dot("Return").Call(),
				jen.Line(),
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Paramsln(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"), jen.Error()).Block(
					jen.Return().List(jen.ID("uc"), jen.Nil()),
				),
				jen.Line(),
				jen.ID("dm").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "WebhookDataManager").Values(),
				jen.ID("dm").Dot("On").Call(jen.Lit("GetAllWebhooksCount"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.ID("expectation"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideWebhooksService").Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("dm"),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(), jen.ID("ucp"),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
				),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	)

	return ret
}
