package items

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(ret)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	srn := typ.Name.RouteName()

	ret.Add(
		jen.Const().Defs(
			jen.Comment("CreateMiddlewareCtxKey is a string alias we can use for referring to item input data in contexts"),
			jen.ID("CreateMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("item_create_input"),
			jen.Comment("UpdateMiddlewareCtxKey is a string alias we can use for referring to item update data in contexts"),
			jen.ID("UpdateMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("item_update_input"),
			jen.Line(),
			jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName").Op("=").Lit("items"),
			jen.ID("counterDescription").Op("=").Lit("the number of items managed by the items service"),
			jen.ID("topicName").ID("string").Op("=").Lit("items"),
			jen.ID("serviceName").ID("string").Op("=").Lit("items_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").ID("models").Dot("ItemDataServer").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Service handles to-do list items"),
			jen.ID("Service").Struct(
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID(fmt.Sprintf("%sCounter", srn)).Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
				jen.ID(fmt.Sprintf("%sDatabase", srn)).ID("models").Dot("ItemDataManager"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID(fmt.Sprintf("%sIDFetcher", srn)).ID(fmt.Sprintf("%sIDFetcher", sn)),
				jen.ID("encoderDecoder").ID("encoding").Dot("EncoderDecoder"),
				jen.ID("reporter").ID("newsman").Dot("Reporter"),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher is a function that fetches user IDs"),
			jen.ID("UserIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.Line(),
			jen.Comment("ItemIDFetcher is a function that fetches item IDs"),
			jen.ID(fmt.Sprintf("%sIDFetcher", sn)).Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Provide%sService builds a new ItemsService", pn)),
		jen.Line(),
		jen.Func().ID("ProvideItemsService").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("db").ID("models").Dot("ItemDataManager"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
			jen.ID("itemIDFetcher").ID("ItemIDFetcher"),
			jen.ID("encoder").ID("encoding").Dot("EncoderDecoder"),
			jen.ID("itemCounterProvider").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider"),
			jen.ID("reporter").ID("newsman").Dot("Reporter"),
		).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID("itemCounter"), jen.ID("err")).Op(":=").ID("itemCounterProvider").Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID(fmt.Sprintf("%sDatabase", uvn)).Op(":").ID("db"),
				jen.ID("encoderDecoder").Op(":").ID("encoder"),
				jen.ID(fmt.Sprintf("%sCounter", uvn)).Op(":").ID(fmt.Sprintf("%sCounter", uvn)),
				jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"),
				jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).Op(":").ID(fmt.Sprintf("%sIDFetcher", uvn)),
				jen.ID("reporter").Op(":").ID("reporter"),
			),
			jen.Line(),
			jen.List(jen.ID("itemCount"), jen.ID("err")).Op(":=").ID("svc").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("GetAll%sCount", pn)).Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting current item count: %w"), jen.ID("err"))),
			),
			jen.ID("svc").Dot("itemCounter").Dot("IncrementBy").Call(jen.ID("ctx"), jen.ID("itemCount")),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Line(),
	)
	return ret
}
