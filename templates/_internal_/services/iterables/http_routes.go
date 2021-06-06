package iterables

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.Comment("ItemIDURIParamKey is a standard string that we'll use to refer to item IDs with."),
			jen.ID("ItemIDURIParamKey").Op("=").Lit("itemID"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("parseBool differs from strconv.ParseBool in that it returns false by default."),
		jen.Newline(),
		jen.Func().ID("parseBool").Params(jen.ID("str").ID("string")).Params(jen.ID("bool")).Body(
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("str")))).Body(
				jen.Case(jen.Lit("1"), jen.Lit("t"), jen.Lit("true")).Body(
					jen.Return().ID("true")),
				jen.Default().Body(
					jen.Return().ID("false")),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CreateHandler is our item creation route."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.Newline(),
			jen.Comment("check session context data for parsed input struct."),
			jen.ID("input").Op(":=").ID("new").Call(jen.Qual(proj.TypesPackage(), "ItemCreationInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("ctx"), jen.ID("req"), jen.ID("input")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.If(jen.ID("err").Op("=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "ValidationErrorKey"),
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
			jen.Newline(),
			jen.ID("input").Dot("BelongsToAccount").Op("=").ID("sessionCtxData").Dot("ActiveAccountID"),
			jen.Newline(),
			jen.Comment("create item in database."),
			jen.List(jen.ID("item"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("CreateItem").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("item").Dot("ID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "ItemIDKey"),
				jen.ID("item").Dot("ID"),
			),
			jen.Newline(),
			jen.Comment("notify interested parties."),
			func() jen.Code {
				return jen.Null()
			}(),
			jen.If(jen.ID("searchIndexErr").Op(":=").ID("s").Dot("search").Dot("Index").Call(jen.ID("ctx"), jen.ID("item").Dot("ID"), jen.ID("item")), jen.ID("searchIndexErr").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(jen.ID("err"), jen.ID("logger"), jen.ID("span"), jen.Lit("adding item to search index")),
			),
			jen.Newline(),
			jen.ID("s").Dot("itemCounter").Dot("Increment").Call(jen.ID("ctx")),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("item"),
				jen.Qual("net/http", "StatusCreated"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ReadHandler returns a GET handler that returns an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.Newline(),
			jen.Comment("determine item ID."),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.Comment("fetch item from database."),
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
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("x"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ExistenceHandler returns a HEAD handler that returns 200 if an item exists, 404 otherwise."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ExistenceHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.Newline(),
			jen.Comment("determine item ID."),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.Comment("check the database."),
			jen.List(jen.ID("exists"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("ItemExists").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.If(jen.Op("!").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("checking item existence"),
				)),
			jen.Newline(),
			jen.If(jen.Op("!").ID("exists").Op("||").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				)),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ListHandler is our list route."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("filter").Op(":=").Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).
				Dotln("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "FilterLimitKey"),
				jen.ID("filter").Dot("Limit"),
			).
				Dotln("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "FilterPageKey"),
				jen.ID("filter").Dot("Page"),
			).
				Dotln("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "FilterSortByKey"),
				jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
			),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter").Dot("Page"),
				jen.ID("filter").Dot("Limit"),
				jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.Newline(),
			jen.List(jen.ID("items"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("GetItems").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("filter"),
			),
			jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.Comment("in the event no rows exist, return an empty list."),
				jen.ID("items").Op("=").Op("&").Qual(proj.TypesPackage(), "ItemList").Values(jen.ID("Items").Op(":").Index().Op("*").Qual(proj.TypesPackage(), "Item").Values())).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("items"),
			),
		),
		jen.Newline(),
	)

	if typ.SearchEnabled {
		code.Add(
			jen.Comment("SearchHandler is our search route."),
			jen.Newline(),
			jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("SearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Newline(),
				jen.ID("query").Op(":=").ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.Qual(proj.TypesPackage(), "SearchQueryKey")),
				jen.ID("filter").Op(":=").Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).
					Dotln("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "FilterLimitKey"),
					jen.ID("filter").Dot("Limit"),
				).
					Dotln("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "FilterPageKey"),
					jen.ID("filter").Dot("Page"),
				).
					Dotln("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "FilterSortByKey"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				).
					Dotln("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "SearchQueryKey"),
					jen.ID("query"),
				),
				jen.Newline(),
				jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				),
				jen.Newline(),
				jen.Comment("determine user ID."),
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
				jen.Newline(),
				jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
					jen.ID("span"),
					jen.ID("sessionCtxData"),
				),
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				),
				jen.Newline(),
				jen.List(jen.ID("relevantIDs"), jen.ID("err")).Op(":=").ID("s").Dot("search").Dot("Search").Call(
					jen.ID("ctx"),
					jen.ID("query"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
				jen.Newline(),
				jen.Comment("fetch items from database."),
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
					jen.Comment("in the event no rows exist, return an empty list."),
					jen.ID("items").Op("=").Index().Op("*").Qual(proj.TypesPackage(), "Item").Values()).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
				jen.Newline(),
				jen.Comment("encode our response and peace."),
				jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("items"),
				),
			),
			jen.Newline(),
		)
	}

	code.Add(
		jen.Comment("UpdateHandler returns a handler that updates an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("UpdateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.Newline(),
			jen.Comment("check for parsed input attached to session context data."),
			jen.ID("input").Op(":=").ID("new").Call(jen.Qual(proj.TypesPackage(), "ItemUpdateInput")),
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
			jen.Newline(),
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
			jen.Newline(),
			jen.Comment("determine item ID."),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.Comment("fetch item from database."),
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
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Comment("update the item."),
			jen.ID("changeReport").Op(":=").ID("item").Dot("Update").Call(jen.ID("input")),
			jen.Qual(proj.InternalTracingPackage(), "AttachChangeSummarySpan").Call(
				jen.ID("span"),
				jen.Lit("item"),
				jen.ID("changeReport"),
			),
			jen.Newline(),
			jen.Comment("update item in database."),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("itemDataManager").Dot("UpdateItem").Call(
				jen.ID("ctx"),
				jen.ID("item"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("changeReport"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Comment("notify interested parties."),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.If(jen.ID("searchIndexErr").Op(":=").ID("s").Dot("search").Dot("Index").Call(
						jen.ID("ctx"),
						jen.ID("item").Dot("ID"),
						jen.ID("item"),
					), jen.ID("searchIndexErr").Op("!=").ID("nil")).Body(
						jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("updating item in search index"),
						),
					)
				}
				return jen.Null()
			}(),
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("item"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ArchiveHandler returns a handler that archives an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			).Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "AccountIDKey"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.Newline(),
			jen.Comment("determine item ID."),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.Comment("archive the item in the database."),
			jen.ID("err").Op("=").ID("s").Dot("itemDataManager").Dot("ArchiveItem").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Comment("notify interested parties."),
			jen.ID("s").Dot("itemCounter").Dot("Decrement").Call(jen.ID("ctx")),
			jen.Newline(),
			jen.If(jen.ID("indexDeleteErr").Op(":=").ID("s").Dot("search").Dot("Delete").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			), jen.ID("indexDeleteErr").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("removing from search index"),
				)),
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNoContent")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AuditEntryHandler returns a GET handler that returns all audit log entries related to an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.Newline(),
			jen.Comment("determine item ID."),
			jen.ID("itemID").Op(":=").ID("s").Dot("itemIDFetcher").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), "ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("itemDataManager").Dot("GetAuditLogEntriesForItem").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
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
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("x"),
			),
		),
		jen.Newline(),
	)

	return code
}
