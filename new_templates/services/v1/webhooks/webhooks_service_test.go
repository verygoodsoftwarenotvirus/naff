package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksServiceTestDotGo() *jen.File {
	ret := jen.NewFile("webhooks")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("buildTestService").Params().Params(jen.Op("*").ID("Service")).Block(
		jen.Return().Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
			"ProvideNoopLogger",
		).Call(), jen.ID("webhookCounter").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(), jen.ID("webhookDatabase").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(), jen.ID("userIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().Lit(0),
		), jen.ID("webhookIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().Lit(0),
		), jen.ID("encoderDecoder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("eventManager").Op(":").ID("newsman").Dot(
			"NewNewsman",
		).Call(jen.ID("nil"), jen.ID("nil"))),
	),
	)
	ret.Add(jen.Func().ID("TestProvideWebhooksService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot(
				"Return",
			).Call(),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.ID("dm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("dm").Dot(
				"On",
			).Call(jen.Lit("GetAllWebhooksCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expectation"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideWebhooksService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("dm"), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("newsman").Dot(
				"NewNewsman",
			).Call(jen.ID("nil"), jen.ID("nil"))),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error providing counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideWebhooksService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("newsman").Dot(
				"NewNewsman",
			).Call(jen.ID("nil"), jen.ID("nil"))),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error setting count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot(
				"Return",
			).Call(),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.ID("dm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataManager").Valuesln(),
			jen.ID("dm").Dot(
				"On",
			).Call(jen.Lit("GetAllWebhooksCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expectation"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideWebhooksService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("dm"), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("newsman").Dot(
				"NewNewsman",
			).Call(jen.ID("nil"), jen.ID("nil"))),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),
	)
	return ret
}
