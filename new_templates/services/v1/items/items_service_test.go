package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func itemsServiceTestDotGo() *jen.File {
	ret := jen.NewFile("items")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("buildTestService").Params().Params(jen.Op("*").ID("Service")).Block(
		jen.Return().Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
			"ProvideNoopLogger",
		).Call(), jen.ID("itemCounter").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(), jen.ID("itemDatabase").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataManager").Valuesln(), jen.ID("userIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().Lit(0),
		), jen.ID("itemIDFetcher").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().Lit(0),
		), jen.ID("encoderDecoder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("reporter").Op(":").ID("nil")),
	),
	)
	ret.Add(jen.Func().ID("TestProvideItemsService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.ID("idm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataManager").Valuesln(),
			jen.ID("idm").Dot(
				"On",
			).Call(jen.Lit("GetAllItemsCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expectation"), jen.ID("nil")),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideItemsService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("idm"), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("s")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error providing unit counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			),
			jen.ID("idm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataManager").Valuesln(),
			jen.ID("idm").Dot(
				"On",
			).Call(jen.Lit("GetAllItemsCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expectation"), jen.ID("nil")),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideItemsService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("idm"), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("nil")),
			jen.ID("require").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("s")),
			jen.ID("require").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching item count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.ID("idm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataManager").Valuesln(),
			jen.ID("idm").Dot(
				"On",
			).Call(jen.Lit("GetAllItemsCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expectation"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expectation")).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("ProvideItemsService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("idm"), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("nil")),
			jen.ID("require").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("s")),
			jen.ID("require").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),
	)
	return ret
}
