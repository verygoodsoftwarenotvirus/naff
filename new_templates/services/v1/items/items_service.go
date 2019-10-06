package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func itemsServiceDotGo() *jen.File {
	ret := jen.NewFile("items")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("CreateMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("item_create_input").Var().ID("UpdateMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("item_update_input").Var().ID("counterName").ID("metrics").Dot(
		"CounterName",
	).Op("=").Lit("items").Var().ID("counterDescription").Op("=").Lit("the number of items managed by the items service").Var().ID("topicName").ID("string").Op("=").Lit("items").Var().ID("serviceName").ID("string").Op("=").Lit("items_service"),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"ItemDataServer",
	).Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("Service").Struct(
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("itemCounter").ID("metrics").Dot(
			"UnitCounter",
		),
		jen.ID("itemDatabase").ID("models").Dot(
			"ItemDataManager",
		),
		jen.ID("userIDFetcher").ID("UserIDFetcher"),
		jen.ID("itemIDFetcher").ID("ItemIDFetcher"),
		jen.ID("encoderDecoder").ID("encoding").Dot(
			"EncoderDecoder",
		),
		jen.ID("reporter").ID("newsman").Dot(
			"Reporter",
		),
	).Type().ID("UserIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Type().ID("ItemIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
