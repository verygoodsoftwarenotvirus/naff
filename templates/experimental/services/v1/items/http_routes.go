package items

import (
	"fmt"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	scn := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")       // singular common name
	pcn := strings.Join(strings.Split(typ.Name.PluralRouteName(), "_"), " ") // plural common name
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	pn := typ.Name.Plural()
	// prn := typ.Name.PluralRouteName()

	xID := fmt.Sprintf("%sID", uvn)

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment(fmt.Sprintf("URIParamKey is a standard string that we'll use to refer to %s IDs with", scn)),
			jen.ID("URIParamKey").Op("=").Lit(fmt.Sprintf("%sID", uvn)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID(fmt.Sprintf("attach%sIDToSpan", sn)).Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID(xID).ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit(fmt.Sprintf("%s_id", rn)), jen.Qual("strconv", "FormatUint").Call(jen.ID(xID), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is our list route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ListHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("ensure query filter"),
				jen.ID("qf").Op(":=").ID("models").Dot("ExtractQueryFilter").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("determine user ID"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment(fmt.Sprintf("fetch %s from database", pcn)),
				jen.List(jen.ID(puvn), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot("GetItems").Call(jen.ID("ctx"), jen.ID("qf"), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Comment("in the event no rows exist return an empty list"),
					jen.ID(puvn).Op("=").Op("&").ID("models").Dot("ItemList").Valuesln(
						jen.ID(pn).Op(":").Index().ID("models").Dot(sn).Values(),
					),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered fetching items")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode our response and peace"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID(puvn)), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler is our item creation route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreateHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("determine user ID"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("check request context for parsed input struct"),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot("Value").Call(jen.ID("CreateMiddlewareCtxKey")).Assert(jen.Op("*").ID("models").Dot("ItemCreationInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("logger").Dot("Info").Call(jen.Lit("valid input not attached to request")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(jen.Lit("input"), jen.ID("input")),
				jen.ID("input").Dot("BelongsTo").Op("=").ID("userID"),
				jen.Line(),
				jen.Comment("create item in database"),
				jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot("CreateItem").Call(jen.ID("ctx"), jen.ID("input")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error creating item")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify relevant parties"),
				jen.ID("s").Dot("itemCounter").Dot("Increment").Call(jen.ID("ctx")),
				jen.ID("attachItemIDToSpan").Call(jen.ID("span"), jen.ID("x").Dot("ID")),
				jen.ID("s").Dot("reporter").Dot("Report").Call(jen.ID("newsman").Dot("Event").Valuesln(
					jen.ID("Data").Op(":").ID("x"),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName")),
					jen.ID("EventType").Op(":").ID("string").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Create")),
				)),
				jen.Line(),
				jen.Comment("encode our response and peace"),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusCreated")),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler returns a GET handler that returns an item"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ReadHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("determine relevant information"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID(xID).Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("user_id").Op(":").ID("userID"),
					jen.Lit("item_id").Op(":").ID("itemID"),
				)),
				jen.ID("attachItemIDToSpan").Call(jen.ID("span"), jen.ID(xID)),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("fetch item from database"),
				jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot("GetItem").Call(jen.ID("ctx"), jen.ID(xID), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
					jen.Return(),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching item from database")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode our response and peace"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateHandler returns a handler that updates an item"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UpdateHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("check for parsed input attached to request context"),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot("Value").Call(jen.ID("UpdateMiddlewareCtxKey")).Assert(jen.Op("*").ID("models").Dot("ItemUpdateInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("no input attached to request")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("determine relevant information"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID(xID).Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("user_id").Op(":").ID("userID"),
					jen.Lit("item_id").Op(":").ID("itemID"),
				)),
				jen.ID("attachItemIDToSpan").Call(jen.ID("span"), jen.ID(xID)),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("fetch item from database"),
				jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot("GetItem").Call(jen.ID("ctx"), jen.ID(xID), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
					jen.Return(),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered getting item")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("update the data structure"),
				jen.ID("x").Dot("Update").Call(jen.ID("input")),
				jen.Line(),
				jen.Comment("update item in database"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot("UpdateItem").Call(jen.ID("ctx"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered updating item")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify relevant parties"),
				jen.ID("s").Dot("reporter").Dot("Report").Call(jen.ID("newsman").Dot("Event").Valuesln(
					jen.ID("Data").Op(":").ID("x"),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName")),
					jen.ID("EventType").Op(":").ID("string").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Update")),
				)),
				jen.Line(),
				jen.Comment("encode our response and peace"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler returns a handler that archives an item"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ArchiveHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("determine relevant information"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID(xID).Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("item_id").Op(":").ID("itemID"),
					jen.Lit("user_id").Op(":").ID("userID"),
				)),
				jen.ID("attachItemIDToSpan").Call(jen.ID("span"), jen.ID(xID)),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("archive the item in the database"),
				jen.ID("err").Op(":=").ID("s").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot("ArchiveItem").Call(jen.ID("ctx"), jen.ID(xID), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
					jen.Return(),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered deleting item")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify relevant parties"),
				jen.ID("s").Dot("itemCounter").Dot("Decrement").Call(jen.ID("ctx")),
				jen.ID("s").Dot("reporter").Dot("Report").Call(jen.ID("newsman").Dot("Event").Valuesln(
					jen.ID("EventType").Op(":").ID("string").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Archive")),
					jen.ID("Data").Op(":").Op("&").ID("models").Dot(sn).Values(jen.ID("ID").Op(":").ID("itemID")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName")),
				)),
				jen.Line(),
				jen.Comment("encode our response and peace"),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNoContent")),
			),
		),
		jen.Line(),
	)
	return ret
}
