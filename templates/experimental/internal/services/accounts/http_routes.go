package accounts

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
			jen.ID("AccountIDURIParamKey").Op("=").Lit("accountID"),
			jen.ID("UserIDURIParamKey").Op("=").Lit("userID"),
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
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.List(jen.ID("accounts"), jen.ID("err")).Op(":=").ID("s").Dot("accountDataManager").Dot("GetAccounts").Call(
				jen.ID("ctx"),
				jen.ID("requester"),
				jen.ID("filter"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("accounts").Op("=").Op("&").ID("types").Dot("AccountList").Valuesln(jen.ID("Accounts").Op(":").Index().Op("*").ID("types").Dot("Account").Valuesln())).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching accounts"),
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
				jen.ID("accounts"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler is our account creation route."),
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
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("AccountCreationInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request body"),
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
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("NameKey"),
				jen.ID("input").Dot("Name"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("input").Dot("BelongsToUser").Op("=").ID("requester"),
			jen.List(jen.ID("account"), jen.ID("err")).Op(":=").ID("s").Dot("accountDataManager").Dot("CreateAccount").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("requester"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating account"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("account").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("account").Dot("ID"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("created account")),
			jen.ID("s").Dot("accountCounter").Dot("Increment").Call(jen.ID("ctx")),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("account"),
				jen.Qual("net/http", "StatusCreated"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler returns a GET handler that returns an account."),
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
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("account"), jen.ID("err")).Op(":=").ID("s").Dot("accountDataManager").Dot("GetAccount").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("requester"),
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
					jen.Lit("fetching account from database"),
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
				jen.ID("account"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateHandler returns a handler that updates an account."),
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
					jen.Lit("fetching session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("AccountUpdateInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request body"),
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
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("input").Dot("BelongsToUser").Op("=").ID("requester"),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("account"), jen.ID("err")).Op(":=").ID("s").Dot("accountDataManager").Dot("GetAccount").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("requester"),
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
					jen.Lit("fetching account from database"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("changeReport").Op(":=").ID("account").Dot("Update").Call(jen.ID("input")),
			jen.ID("tracing").Dot("AttachChangeSummarySpan").Call(
				jen.ID("span"),
				jen.Lit("account"),
				jen.ID("changeReport"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("accountDataManager").Dot("UpdateAccount").Call(
				jen.ID("ctx"),
				jen.ID("account"),
				jen.ID("requester"),
				jen.ID("changeReport"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating account"),
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
				jen.ID("account"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler returns a handler that archives an account."),
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
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("err").Op("=").ID("s").Dot("accountDataManager").Dot("ArchiveAccount").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("requester"),
				jen.ID("requester"),
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
					jen.Lit("archiving account"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("accountCounter").Dot("Decrement").Call(jen.ID("ctx")),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNoContent")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AddMemberHandler is our account creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("AddMemberHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("AddUserToAccountInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request body"),
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
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("accountMembershipDataManager").Dot("AddUserToAccount").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("requester"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("adding user to account"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ModifyMemberPermissionsHandler is our account creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ModifyMemberPermissionsHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("ModifyUserPermissionsInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request body"),
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
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("accountMembershipDataManager").Dot("ModifyUserPermissions").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("requester"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("modifying user permissions"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TransferAccountOwnershipHandler is our account creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("TransferAccountOwnershipHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
					jen.Lit("transferring account ownership"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("AccountOwnershipTransferInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request body"),
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
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("accountMembershipDataManager").Dot("TransferAccountOwnership").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("requester"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("transferring account ownership"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RemoveMemberHandler is our account creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("RemoveMemberHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("reason").Op(":=").ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.Lit("reason")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ReasonKey"),
				jen.ID("reason"),
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
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("accountMembershipDataManager").Dot("RemoveUserFromAccount").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("requester"),
				jen.ID("reason"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("removing user from account"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MarkAsDefaultAccountHandler is our account creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("MarkAsDefaultAccountHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
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
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("accountMembershipDataManager").Dot("MarkAccountAsUserDefault").Call(
				jen.ID("ctx"),
				jen.ID("requester"),
				jen.ID("accountID"),
				jen.ID("requester"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("marking account as default"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuditEntryHandler returns a GET handler that returns all audit log entries related to an account."),
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
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("accountID").Op(":=").ID("s").Dot("accountIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("accountDataManager").Dot("GetAuditLogEntriesForAccount").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
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
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching audit log entries"),
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
