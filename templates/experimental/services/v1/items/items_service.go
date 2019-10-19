package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsServiceDotGo() *jen.File {
	ret := jen.NewFile("items")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("CreateMiddlewareCtxKey").ID("models").Dot(
			"ContextKey",
		).Op("=").Lit("item_create_input").Var().ID("UpdateMiddlewareCtxKey").ID("models").Dot(
			"ContextKey",
		).Op("=").Lit("item_update_input").Var().ID("counterName").ID("metrics").Dot(
			"CounterName",
		).Op("=").Lit("items").Var().ID("counterDescription").Op("=").Lit("the number of items managed by the items service").Var().ID("topicName").ID("string").Op("=").Lit("items").Var().ID("serviceName").ID("string").Op("=").Lit("items_service"),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("models").Dot(
			"ItemDataServer",
		).Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("Service").Struct(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
	),
	jen.ID("itemCounter").ID("metrics").Dot(
			"UnitCounter",
	),
	jen.ID("itemDatabase").ID("models").Dot(
			"ItemDataManager",
	),
	jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("itemIDFetcher").ID("ItemIDFetcher"), jen.ID("encoderDecoder").ID("encoding").Dot(
			"EncoderDecoder",
	),
	jen.ID("reporter").ID("newsman").Dot(
			"Reporter",
		)).Type().ID("UserIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Type().ID("ItemIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideItemsService builds a new ItemsService"),
		jen.Line(),
		jen.Func().ID("ProvideItemsService").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
	),
	jen.ID("db").ID("models").Dot(
			"ItemDataManager",
	),
	jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("itemIDFetcher").ID("ItemIDFetcher"), jen.ID("encoder").ID("encoding").Dot(
			"EncoderDecoder",
	),
	jen.ID("itemCounterProvider").ID("metrics").Dot(
			"UnitCounterProvider",
	),
	jen.ID("reporter").ID("newsman").Dot(
			"Reporter",
		)).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID("itemCounter"), jen.ID("err")).Op(":=").ID("itemCounterProvider").Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.ID("err"))),
			),
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("logger").Dot(
				"WithName",
			).Call(jen.ID("serviceName")), jen.ID("itemDatabase").Op(":").ID("db"), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("itemCounter").Op(":").ID("itemCounter"), jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"), jen.ID("itemIDFetcher").Op(":").ID("itemIDFetcher"), jen.ID("reporter").Op(":").ID("reporter")),
			jen.List(jen.ID("itemCount"), jen.ID("err")).Op(":=").ID("svc").Dot(
				"itemDatabase",
			).Dot(
				"GetAllItemsCount",
			).Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting current item count: %w"), jen.ID("err"))),
			),
			jen.ID("svc").Dot(
				"itemCounter",
			).Dot(
				"IncrementBy",
			).Call(jen.ID("ctx"), jen.ID("itemCount")),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Line(),
	)
	return ret
}
