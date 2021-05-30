package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ItemIDURIParamKey").Op("=").Lit("itemID"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("parseBool differs from strconv.ParseBool in that it returns false by default."),
		jen.Line(),
		jen.Func().ID("parseBool").Params(jen.ID("str").ID("string")).Params(jen.ID("bool")).Body(
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("str")))).Body(
				jen.Case(jen.Lit("1"), jen.Lit("t"), jen.Lit("true")).Body(
					jen.Return().ID("true")),
				jen.Default().Body(
					jen.Return().ID("false")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler is our item creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("ItemCreationInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("invalid request content"),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.If(jen.ID("err").Op("=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("ValidationErrorKey"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("invalid input attached to request")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.ID("input").Dot("BelongsToAccount").Op("=").ID("sessionCtxData").Dot("ActiveAccountID"),
			jen.List(jen.ID("item"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("CreateItem").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating item"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("item").Dot("ID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("item").Dot("ID"),
			),
			jen.If(jen.ID("searchIndexErr").Op(":=").ID("s").Dot("search").Dot("Index").Call(
				jen.ID("ctx"),
				jen.ID("item").Dot("ID"),
				jen.ID("item"),
			), jen.ID("searchIndexErr").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("adding item to search index"),
				)),
			jen.ID("s").Dot("itemCounter").Dot("Increment").Call(jen.ID("ctx")),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("item"),
				jen.Qual("net/http", "StatusCreated"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler returns a GET handler that returns an item."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("GetItem").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving item"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("x"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ExistenceHandler returns a HEAD handler that returns 200 if an item exists, 404 otherwise."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ExistenceHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.List(jen.ID("exists"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("ItemExists").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.If(jen.Op("!").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("checking item existence"),
				)),
			jen.If(jen.Op("!").ID("exists").Op("||").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ListHandler is our list route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("filter").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterLimitKey"),
				jen.ID("filter").Dot("Limit"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterPageKey"),
				jen.ID("filter").Dot("Page"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterSortByKey"),
				jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
			),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("tracing").Dot("AttachFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter").Dot("Page"),
				jen.ID("filter").Dot("Limit"),
				jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.List(jen.ID("items"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("GetItems").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("filter"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("items").Op("=").Op("&").ID("types").Dot("ItemList").Valuesln(jen.ID("Items").Op(":").Index().Op("*").ID("types").Dot("Item").Valuesln())).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving items"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("items"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SearchHandler is our search route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("SearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("query").Op(":=").ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.ID("types").Dot("SearchQueryKey")),
			jen.ID("filter").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterLimitKey"),
				jen.ID("filter").Dot("Limit"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterPageKey"),
				jen.ID("filter").Dot("Page"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("FilterSortByKey"),
				jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("SearchQueryKey"),
				jen.ID("query"),
			),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("tracing").Dot("AttachFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter").Dot("Page"),
				jen.ID("filter").Dot("Limit"),
				jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.List(jen.ID("relevantIDs"), jen.ID("err")).Op(":=").ID("s").Dot("search").Dot("Search").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing item search query"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("items"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("GetItemsWithIDs").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("filter").Dot("Limit"),
				jen.ID("relevantIDs"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("items").Op("=").Index().Op("*").ID("types").Dot("Item").Valuesln()).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("searching items"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("items"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateHandler returns a handler that updates an item."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("UpdateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("ItemUpdateInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("error encountered decoding request body"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("invalid request content"),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.If(jen.ID("err").Op("=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("provided input was invalid"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.ID("input").Dot("BelongsToAccount").Op("=").ID("sessionCtxData").Dot("ActiveAccountID"),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.List(jen.ID("item"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("GetItem").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving item for update"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("changeReport").Op(":=").ID("item").Dot("Update").Call(jen.ID("input")),
			jen.ID("tracing").Dot("AttachChangeSummarySpan").Call(
				jen.ID("span"),
				jen.Lit("item"),
				jen.ID("changeReport"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("itemDataManager").Dot("UpdateItem").Call(
				jen.ID("ctx"),
				jen.ID("item"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("changeReport"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating item"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.If(jen.ID("searchIndexErr").Op(":=").ID("s").Dot("search").Dot("Index").Call(
				jen.ID("ctx"),
				jen.ID("item").Dot("ID"),
				jen.ID("item"),
			), jen.ID("searchIndexErr").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating item in search index"),
				)),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("item"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler returns a handler that archives an item."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("err").Op("=").ID("s").Dot("itemDataManager").Dot("ArchiveItem").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving item"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("itemCounter").Dot("Decrement").Call(jen.ID("ctx")),
			jen.If(jen.ID("indexDeleteErr").Op(":=").ID("s").Dot("search").Dot("Delete").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			), jen.ID("indexDeleteErr").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("removing from search index"),
				)),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNoContent")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuditEntryHandler returns a GET handler that returns all audit log entries related to an item."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("GetAuditLogEntriesForItem").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving audit log entries for item"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("x"),
			),
		),
		jen.Line(),
	)

	return code
}
