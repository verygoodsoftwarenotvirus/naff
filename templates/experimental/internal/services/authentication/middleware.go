package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("signatureHeaderKey").Op("=").Lit("Signature"),
			jen.ID("pasetoAuthorizationHeaderKey").Op("=").Lit("Authorization"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errTokenExpired").Op("=").Qual("errors", "New").Call(jen.Lit("token expired")),
			jen.ID("errTokenNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("no token data found")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchSessionContextDataFromPASETO").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.If(jen.ID("rawToken").Op(":=").ID("req").Dot("Header").Dot("Get").Call(jen.ID("pasetoAuthorizationHeaderKey")), jen.ID("rawToken").Op("!=").Lit("")).Body(
				jen.Var().Defs(
					jen.ID("token").ID("paseto").Dot("JSONToken"),
				),
				jen.If(jen.ID("err").Op(":=").ID("paseto").Dot("NewV2").Call().Dot("Decrypt").Call(
					jen.ID("rawToken"),
					jen.ID("s").Dot("config").Dot("PASETO").Dot("LocalModeKey"),
					jen.Op("&").ID("token"),
					jen.ID("nil"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("decrypting PASETO"),
					))),
				jen.If(jen.Qual("time", "Now").Call().Dot("UTC").Call().Dot("After").Call(jen.ID("token").Dot("Expiration"))).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("errTokenExpired"))),
				jen.ID("base64Encoded").Op(":=").ID("token").Dot("Get").Call(jen.ID("pasetoDataKey")),
				jen.List(jen.ID("gobEncoded"), jen.ID("err")).Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("DecodeString").Call(jen.ID("base64Encoded")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("decoding base64 encoded GOB payload"),
					))),
				jen.Var().Defs(
					jen.ID("reqContext").Op("*").ID("types").Dot("SessionContextData"),
				),
				jen.If(jen.ID("err").Op("=").Qual("encoding/gob", "NewDecoder").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("gobEncoded"))).Dot("Decode").Call(jen.Op("&").ID("reqContext")), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("decoding GOB encoded session info payload"),
					))),
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("active_account_id"),
					jen.ID("reqContext").Dot("ActiveAccountID"),
				).Dot("Debug").Call(jen.Lit("returning session context data")),
				jen.Return().List(jen.ID("reqContext"), jen.ID("nil")),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("errTokenNotFound")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CookieRequirementMiddleware requires every request have a valid cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("CookieRequirementMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.If(jen.List(jen.ID("cookie"), jen.ID("cookieErr")).Op(":=").ID("req").Dot("Cookie").Call(jen.ID("s").Dot("config").Dot("Cookies").Dot("Name")), jen.Op("!").Qual("errors", "Is").Call(
					jen.ID("cookieErr"),
					jen.Qual("net/http", "ErrNoCookie"),
				).Op("&&").ID("cookie").Op("!=").ID("nil")).Body(
					jen.Var().Defs(
						jen.ID("token").ID("string"),
					),
					jen.If(jen.ID("err").Op(":=").ID("s").Dot("cookieManager").Dot("Decode").Call(
						jen.ID("s").Dot("config").Dot("Cookies").Dot("Name"),
						jen.ID("cookie").Dot("Value"),
						jen.Op("&").ID("token"),
					), jen.ID("err").Op("==").ID("nil")).Body(
						jen.ID("next").Dot("ServeHTTP").Call(
							jen.ID("res"),
							jen.ID("req"),
						)),
				),
				jen.Qual("net/http", "Redirect").Call(
					jen.ID("res"),
					jen.ID("req"),
					jen.Lit("/users/login"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserAttributionMiddleware is concerned with figuring out who a user is, but not worried about kicking out users who are not known."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("UserAttributionMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.If(jen.List(jen.ID("cookieContext"), jen.ID("userID"), jen.ID("err")).Op(":=").ID("s").Dot("getUserIDFromCookie").Call(
					jen.ID("ctx"),
					jen.ID("req"),
				), jen.ID("err").Op("==").ID("nil").Op("&&").ID("userID").Op("!=").Lit(0)).Body(
					jen.ID("ctx").Op("=").ID("cookieContext"),
					jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
						jen.ID("span"),
						jen.ID("userID"),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("RequesterIDKey"),
						jen.ID("userID"),
					),
					jen.List(jen.ID("sessionCtxData"), jen.ID("sessionCtxDataErr")).Op(":=").ID("s").Dot("accountMembershipManager").Dot("BuildSessionContextDataForUser").Call(
						jen.ID("ctx"),
						jen.ID("userID"),
					),
					jen.If(jen.ID("sessionCtxDataErr").Op("!=").ID("nil")).Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("sessionCtxDataErr"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("fetching user info for cookie"),
						),
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
						),
						jen.Return(),
					),
					jen.ID("s").Dot("overrideSessionContextDataValuesWithSessionData").Call(
						jen.ID("ctx"),
						jen.ID("sessionCtxData"),
					),
					jen.ID("next").Dot("ServeHTTP").Call(
						jen.ID("res"),
						jen.ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
							jen.ID("ctx"),
							jen.ID("types").Dot("SessionContextDataKey"),
							jen.ID("sessionCtxData"),
						)),
					),
					jen.Return(),
				),
				jen.ID("logger").Dot("Debug").Call(jen.Lit("no cookie attached to request")),
				jen.List(jen.ID("tokenSessionContextData"), jen.ID("err")).Op(":=").ID("s").Dot("fetchSessionContextDataFromPASETO").Call(
					jen.ID("ctx"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil").Op("&&").Op("!").Parens(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.ID("errTokenNotFound"),
				).Op("||").Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.ID("errTokenExpired"),
				))).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("extracting token from request"),
					)),
				jen.If(jen.ID("tokenSessionContextData").Op("!=").ID("nil")).Body(
					jen.ID("next").Dot("ServeHTTP").Call(
						jen.ID("res"),
						jen.ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
							jen.ID("ctx"),
							jen.ID("types").Dot("SessionContextDataKey"),
							jen.ID("tokenSessionContextData"),
						)),
					),
					jen.Return(),
				),
				jen.ID("next").Dot("ServeHTTP").Call(
					jen.ID("res"),
					jen.ID("req"),
				),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuthorizationMiddleware checks to see if a user is associated with the request, and then determines whether said request can proceed."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("AuthorizationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.If(jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")), jen.ID("err").Op("==").ID("nil").Op("&&").ID("sessionCtxData").Op("!=").ID("nil")).Body(
					jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("RequesterIDKey"),
						jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
					),
					jen.If(jen.ID("sessionCtxData").Dot("Requester").Dot("Reputation").Op("==").ID("types").Dot("BannedUserAccountStatus").Op("||").ID("sessionCtxData").Dot("Requester").Dot("Reputation").Op("==").ID("types").Dot("TerminatedUserReputation")).Body(
						jen.ID("logger").Dot("Debug").Call(jen.Lit("banned user attempted to make request")),
						jen.Qual("net/http", "Redirect").Call(
							jen.ID("res"),
							jen.ID("req"),
							jen.Lit("/"),
							jen.Qual("net/http", "StatusForbidden"),
						),
						jen.Return(),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
						jen.Lit("requested_account_id"),
						jen.ID("sessionCtxData").Dot("ActiveAccountID"),
					),
					jen.If(jen.List(jen.ID("_"), jen.ID("authorizedForAccount")).Op(":=").ID("sessionCtxData").Dot("AccountPermissions").Index(jen.ID("sessionCtxData").Dot("ActiveAccountID")), jen.Op("!").ID("authorizedForAccount")).Body(
						jen.ID("logger").Dot("Debug").Call(jen.Lit("user trying to access account they are not authorized for")),
						jen.Qual("net/http", "Redirect").Call(
							jen.ID("res"),
							jen.ID("req"),
							jen.Lit("/"),
							jen.Qual("net/http", "StatusUnauthorized"),
						),
						jen.Return(),
					),
					jen.ID("next").Dot("ServeHTTP").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Return(),
				),
				jen.ID("logger").Dot("Debug").Call(jen.Lit("no user attached to request")),
				jen.Qual("net/http", "Redirect").Call(
					jen.ID("res"),
					jen.ID("req"),
					jen.Lit("/users/login"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PermissionFilterMiddleware filters users out of requests based on provided functions."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("PermissionFilterMiddleware").Params(jen.ID("permissions").Op("...").ID("authorization").Dot("Permission")).Params(jen.Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))).Body(
			jen.Return().Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
				jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
					jen.List(jen.ID("sessionContextData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("retrieving session context data"),
						),
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
						),
						jen.Return(),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("sessionContextData").Dot("Requester").Dot("UserID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("sessionContextData").Dot("ActiveAccountID"), jen.Lit("account_perms").Op(":").ID("sessionContextData").Dot("AccountPermissions"))),
					jen.ID("logger").Dot("Debug").Call(jen.Lit("PermissionFilterMiddleware called")),
					jen.ID("isServiceAdmin").Op(":=").ID("sessionContextData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call(),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("UserIsServiceAdminKey"),
						jen.ID("isServiceAdmin"),
					),
					jen.List(jen.ID("_"), jen.ID("allowed")).Op(":=").ID("sessionContextData").Dot("AccountPermissions").Index(jen.ID("sessionContextData").Dot("ActiveAccountID")),
					jen.If(jen.Op("!").ID("allowed").Op("&&").Op("!").ID("isServiceAdmin")).Body(
						jen.ID("logger").Dot("Debug").Call(jen.Lit("not authorized for account!")),
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
						),
						jen.Return(),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("RequesterIDKey"),
						jen.ID("sessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("WithValue").Call(
						jen.ID("keys").Dot("AccountIDKey"),
						jen.ID("sessionContextData").Dot("ActiveAccountID"),
					),
					jen.For(jen.List(jen.ID("_"), jen.ID("perm")).Op(":=").Range().ID("permissions")).Body(
						jen.If(jen.Op("!").ID("sessionContextData").Dot("ServiceRolePermissionChecker").Call().Dot("HasPermission").Call(jen.ID("perm")).Op("&&").Op("!").ID("sessionContextData").Dot("AccountRolePermissionsChecker").Call().Dot("HasPermission").Call(jen.ID("perm"))).Body(
							jen.ID("logger").Dot("WithValue").Call(
								jen.Lit("deficient_permission"),
								jen.ID("perm").Dot("ID").Call(),
							).Dot("Debug").Call(jen.Lit("request filtered out")),
							jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
								jen.ID("ctx"),
								jen.ID("res"),
							),
							jen.Return(),
						)),
					jen.ID("next").Dot("ServeHTTP").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
				)))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ServiceAdminMiddleware restricts requests to admin users only."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ServiceAdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Var().Defs(
				jen.ID("staticError").Op("=").Lit("admin status required"),
			),
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
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
						jen.ID("staticError"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.Return(),
				),
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("RequesterIDKey"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				),
				jen.If(jen.Op("!").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()).Body(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("ServiceAdminMiddleware called by non-admin user")),
					jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("staticError"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.Return(),
				),
				jen.ID("next").Dot("ServeHTTP").Call(
					jen.ID("res"),
					jen.ID("req"),
				),
			)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("ErrNoSessionContextDataAvailable").Op("=").Qual("errors", "New").Call(jen.Lit("no SessionContextData attached to session context data")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("FetchContextFromRequest fetches a SessionContextData from a request."),
		jen.Line(),
		jen.Func().ID("FetchContextFromRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
			jen.If(jen.List(jen.ID("sessionCtxData"), jen.ID("ok")).Op(":=").ID("req").Dot("Context").Call().Dot("Value").Call(jen.ID("types").Dot("SessionContextDataKey")).Assert(jen.Op("*").ID("types").Dot("SessionContextData")), jen.ID("ok").Op("&&").ID("sessionCtxData").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("sessionCtxData"), jen.ID("nil"))),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrNoSessionContextDataAvailable")),
		),
		jen.Line(),
	)

	return code
}
