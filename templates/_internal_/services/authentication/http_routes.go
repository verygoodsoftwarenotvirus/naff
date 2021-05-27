package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("issueSessionManagedCookie").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("requesterID")).ID("uint64")).Params(jen.ID("cookie").Op("*").Qual("net/http", "Cookie"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.List(jen.ID("ctx"), jen.ID("err")).Op("=").ID("s").Dot("sessionManager").Dot("Load").Call(
				jen.ID("ctx"),
				jen.Lit(""),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("loading token"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("sessionManager").Dot("RenewToken").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("renewing token"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.ID("s").Dot("sessionManager").Dot("Put").Call(
				jen.ID("ctx"),
				jen.ID("accountIDContextKey"),
				jen.ID("accountID"),
			),
			jen.ID("s").Dot("sessionManager").Dot("Put").Call(
				jen.ID("ctx"),
				jen.ID("userIDContextKey"),
				jen.ID("requesterID"),
			),
			jen.List(jen.ID("token"), jen.ID("expiry"), jen.ID("err")).Op(":=").ID("s").Dot("sessionManager").Dot("Commit").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing to session store"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op("=").ID("s").Dot("buildCookie").Call(
				jen.ID("token"),
				jen.ID("expiry"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building cookie"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Return().List(jen.ID("cookie"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("ErrUserNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("user not found")),
			jen.ID("ErrUserBanned").Op("=").Qual("errors", "New").Call(jen.Lit("user is banned")),
			jen.ID("ErrInvalidCredentials").Op("=").Qual("errors", "New").Call(jen.Lit("invalid credentials")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("AuthenticateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("loginData").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.Op("*").ID("types").Dot("User"), jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("loginData").Dot("Username"),
			),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUserByUsername").Call(
				jen.ID("ctx"),
				jen.ID("loginData").Dot("Username"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("user").Op("==").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("nil"), jen.ID("ErrUserNotFound"))),
				jen.Return().List(jen.ID("nil"), jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching user"),
				)),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("user").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachUserToSpan").Call(
				jen.ID("span"),
				jen.ID("user"),
			),
			jen.If(jen.ID("user").Dot("IsBanned").Call()).Body(
				jen.ID("s").Dot("auditLog").Dot("LogBannedUserLoginAttemptEvent").Call(
					jen.ID("ctx"),
					jen.ID("user").Dot("ID"),
				),
				jen.Return().List(jen.ID("user"), jen.ID("nil"), jen.ID("ErrUserBanned")),
			),
			jen.List(jen.ID("loginValid"), jen.ID("err")).Op(":=").ID("s").Dot("validateLogin").Call(
				jen.ID("ctx"),
				jen.ID("user"),
				jen.ID("loginData"),
			),
			jen.ID("logger").Dot("WithValue").Call(
				jen.Lit("login_valid"),
				jen.ID("loginValid"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.ID("authentication").Dot("ErrInvalidTOTPToken"),
				)).Body(
					jen.ID("s").Dot("auditLog").Dot("LogUnsuccessfulLoginBad2FATokenEvent").Call(
						jen.ID("ctx"),
						jen.ID("user").Dot("ID"),
					),
					jen.Return().List(jen.ID("user"), jen.ID("nil"), jen.ID("ErrInvalidCredentials")),
				).Else().If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.ID("authentication").Dot("ErrPasswordDoesNotMatch"),
				)).Body(
					jen.ID("s").Dot("auditLog").Dot("LogUnsuccessfulLoginBadPasswordEvent").Call(
						jen.ID("ctx"),
						jen.ID("user").Dot("ID"),
					),
					jen.Return().List(jen.ID("user"), jen.ID("nil"), jen.ID("ErrInvalidCredentials")),
				),
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("error encountered validating login"),
				),
				jen.Return().List(jen.ID("user"), jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating login"),
				)),
			).Else().If(jen.Op("!").ID("loginValid")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("login was invalid")),
				jen.ID("s").Dot("auditLog").Dot("LogUnsuccessfulLoginBadPasswordEvent").Call(
					jen.ID("ctx"),
					jen.ID("user").Dot("ID"),
				),
				jen.Return().List(jen.ID("user"), jen.ID("nil"), jen.ID("ErrInvalidCredentials")),
			),
			jen.List(jen.ID("defaultAccountID"), jen.ID("err")).Op(":=").ID("s").Dot("accountMembershipManager").Dot("GetDefaultAccountIDForUser").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("ID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("user"), jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching user memberships"),
				))),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot("issueSessionManagedCookie").Call(
				jen.ID("ctx"),
				jen.ID("defaultAccountID"),
				jen.ID("user").Dot("ID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("user"), jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("issuing cookie"),
				))),
			jen.ID("s").Dot("auditLog").Dot("LogSuccessfulLoginEvent").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("ID"),
			),
			jen.Return().List(jen.ID("user"), jen.ID("cookie"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LoginHandler is our login route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("BeginSessionHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("loginData").Op(":=").ID("new").Call(jen.ID("types").Dot("UserLoginInput")),
			jen.If(jen.ID("err").Op(":=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("loginData"),
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
			jen.If(jen.ID("err").Op(":=").ID("loginData").Dot("ValidateWithContext").Call(
				jen.ID("ctx"),
				jen.ID("s").Dot("config").Dot("MinimumUsernameLength"),
				jen.ID("s").Dot("config").Dot("MinimumPasswordLength"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("loginData").Dot("Username"),
			),
			jen.List(jen.ID("user"), jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot("AuthenticateUser").Call(
				jen.ID("ctx"),
				jen.ID("loginData"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Switch().Body(
					jen.Case(jen.Qual("errors", "Is").Call(
						jen.ID("err"),
						jen.ID("ErrUserNotFound"),
					)).Body(
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
						)),
					jen.Case(jen.Qual("errors", "Is").Call(
						jen.ID("err"),
						jen.ID("ErrUserBanned"),
					)).Body(
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
							jen.ID("user").Dot("ReputationExplanation"),
							jen.Qual("net/http", "StatusForbidden"),
						)),
					jen.Case(jen.Qual("errors", "Is").Call(
						jen.ID("err"),
						jen.ID("ErrInvalidCredentials"),
					)).Body(
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
							jen.Lit("login was invalid"),
							jen.Qual("net/http", "StatusUnauthorized"),
						)),
					jen.Default().Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("issuing cookie"),
						), jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
							jen.ID("staticError"),
							jen.Qual("net/http", "StatusInternalServerError"),
						)),
				),
				jen.Return(),
			),
			jen.Qual("net/http", "SetCookie").Call(
				jen.ID("res"),
				jen.ID("cookie"),
			),
			jen.ID("statusResponse").Op(":=").Op("&").ID("types").Dot("UserStatusResponse").Valuesln(jen.ID("UserIsAuthenticated").Op(":").ID("true"), jen.ID("UserReputation").Op(":").ID("user").Dot("ServiceAccountStatus"), jen.ID("UserReputationExplanation").Op(":").ID("user").Dot("ReputationExplanation")),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("statusResponse"),
				jen.Qual("net/http", "StatusAccepted"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("user logged in")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ChangeActiveAccountHandler is our login route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ChangeActiveAccountHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("ChangeActiveAccountInput")),
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
			jen.ID("accountID").Op(":=").ID("input").Dot("AccountID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("new_session_account_id"),
				jen.ID("accountID"),
			),
			jen.ID("requesterID").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("user_id"),
				jen.ID("requesterID"),
			),
			jen.List(jen.ID("authorizedForAccount"), jen.ID("err")).Op(":=").ID("s").Dot("accountMembershipManager").Dot("UserIsMemberOfAccount").Call(
				jen.ID("ctx"),
				jen.ID("requesterID"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("checking permissions"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("staticError"),
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.If(jen.Op("!").ID("authorizedForAccount")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("invalid account ID requested for activation")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot("issueSessionManagedCookie").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("requesterID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("issuing cookie"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("staticError"),
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("successfully changed active session account")),
			jen.Qual("net/http", "SetCookie").Call(
				jen.ID("res"),
				jen.ID("cookie"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("LogoutUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("res").Qual("net/http", "ResponseWriter")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("ctx"), jen.ID("loadErr")).Op(":=").ID("s").Dot("sessionManager").Dot("Load").Call(
				jen.ID("ctx"),
				jen.Lit(""),
			),
			jen.If(jen.ID("loadErr").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("loadErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("loading token"),
				)),
			jen.If(jen.ID("destroyErr").Op(":=").ID("s").Dot("sessionManager").Dot("Destroy").Call(jen.ID("ctx")), jen.ID("destroyErr").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("destroyErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("destroying user session"),
				)),
			jen.List(jen.ID("newCookie"), jen.ID("cookieBuildingErr")).Op(":=").ID("s").Dot("buildCookie").Call(
				jen.Lit("deleted"),
				jen.Qual("time", "Time").Values(),
			),
			jen.If(jen.ID("cookieBuildingErr").Op("!=").ID("nil").Op("||").ID("newCookie").Op("==").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("cookieBuildingErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building cookie"),
				)),
			jen.ID("s").Dot("auditLog").Dot("LogLogoutEvent").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.ID("newCookie").Dot("MaxAge").Op("=").Op("-").Lit(1),
			jen.Qual("net/http", "SetCookie").Call(
				jen.ID("res"),
				jen.ID("newCookie"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("user logged out")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogoutHandler is our logout route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("EndSessionHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
			jen.If(jen.ID("err").Op("=").ID("s").Dot("LogoutUser").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
				jen.ID("res"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("logging out user"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.Qual("net/http", "Redirect").Call(
				jen.ID("res"),
				jen.ID("req"),
				jen.Lit("/"),
				jen.Qual("net/http", "StatusSeeOther"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("StatusHandler returns the user info for the user making the request."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("StatusHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Var().Defs(
				jen.ID("statusResponse").Op("*").ID("types").Dot("UserStatusResponse"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("statusResponse").Op("=").Op("&").ID("types").Dot("UserStatusResponse").Valuesln(jen.ID("ActiveAccount").Op(":").ID("sessionCtxData").Dot("ActiveAccountID"), jen.ID("UserReputation").Op(":").ID("sessionCtxData").Dot("Requester").Dot("Reputation"), jen.ID("UserReputationExplanation").Op(":").ID("sessionCtxData").Dot("Requester").Dot("ReputationExplanation"), jen.ID("UserIsAuthenticated").Op(":").ID("true")),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("statusResponse"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("pasetoRequestTimeThreshold").Op("=").Lit(2).Op("*").Qual("time", "Minute"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PASETOHandler returns the user info for the user making the request."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("PASETOHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("PASETOCreationInput")),
			jen.If(jen.ID("err").Op(":=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
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
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
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
			jen.ID("requestedAccount").Op(":=").ID("input").Dot("AccountID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientClientIDKey"),
				jen.ID("input").Dot("ClientID"),
			),
			jen.If(jen.ID("requestedAccount").Op("!=").Lit(0)).Body(
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
					jen.Lit("requested_account"),
					jen.ID("requestedAccount"),
				)),
			jen.ID("reqTime").Op(":=").Qual("time", "Unix").Call(
				jen.Lit(0),
				jen.ID("input").Dot("RequestTime"),
			),
			jen.If(jen.Qual("time", "Until").Call(jen.ID("reqTime")).Op(">").ID("pasetoRequestTimeThreshold").Op("||").Qual("time", "Since").Call(jen.ID("reqTime")).Op(">").ID("pasetoRequestTimeThreshold")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("provided_request_time"),
					jen.ID("reqTime").Dot("String").Call(),
				).Dot("Debug").Call(jen.Lit("PASETO request denied because its time is out of threshold")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("sum"), jen.ID("err")).Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("DecodeString").Call(jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("signatureHeaderKey"))),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("len").Call(jen.ID("sum")).Op("==").Lit(0)).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("sum_length"),
					jen.ID("len").Call(jen.ID("sum")),
				).Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("invalid signature"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("client"), jen.ID("clientRetrievalErr")).Op(":=").ID("s").Dot("apiClientManager").Dot("GetAPIClientByClientID").Call(
				jen.ID("ctx"),
				jen.ID("input").Dot("ClientID"),
			),
			jen.If(jen.ID("clientRetrievalErr").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching API client"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
				jen.Qual("crypto/sha256", "New"),
				jen.ID("client").Dot("ClientSecret"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("s").Dot("encoderDecoder").Dot("MustEncodeJSON").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			)), jen.ID("macWriteErr").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing HMAC message for comparison"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.If(jen.Op("!").Qual("crypto/hmac", "Equal").Call(
				jen.ID("sum"),
				jen.ID("mac").Dot("Sum").Call(jen.ID("nil")),
			)).Body(
				jen.ID("logger").Dot("Info").Call(jen.Lit("invalid credentials passed to PASETO creation route")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUser").Call(
				jen.ID("ctx"),
				jen.ID("client").Dot("BelongsToUser"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving user"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("user").Dot("ID"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("accountMembershipManager").Dot("BuildSessionContextDataForUser").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("ID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving perms for API client"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.Var().Defs(
				jen.ID("requestedAccountID").ID("uint64"),
			),
			jen.If(jen.ID("requestedAccount").Op("!=").Lit(0)).Body(
				jen.If(jen.List(jen.ID("_"), jen.ID("isMember")).Op(":=").ID("sessionCtxData").Dot("AccountPermissions").Index(jen.ID("requestedAccount")), jen.Op("!").ID("isMember")).Body(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("invalid account ID requested for token")),
					jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
					jen.Return(),
				),
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("requested_account"),
					jen.ID("requestedAccount"),
				).Dot("Debug").Call(jen.Lit("setting token account ID to requested account")),
				jen.ID("requestedAccountID").Op("=").ID("requestedAccount"),
				jen.ID("sessionCtxData").Dot("ActiveAccountID").Op("=").ID("requestedAccount"),
			).Else().Body(
				jen.ID("requestedAccountID").Op("=").ID("sessionCtxData").Dot("ActiveAccountID")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("requestedAccountID"),
			),
			jen.List(jen.ID("tokenRes"), jen.ID("err")).Op(":=").ID("s").Dot("buildPASETOResponse").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("client"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encrypting PASETO"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("PASETO issued")),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("tokenRes"),
				jen.Qual("net/http", "StatusAccepted"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildPASETOToken").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("client").Op("*").ID("types").Dot("APIClient")).Params(jen.ID("paseto").Dot("JSONToken")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("now").Op(":=").Qual("time", "Now").Call().Dot("UTC").Call(),
			jen.ID("lifetime").Op(":=").Qual("time", "Duration").Call(jen.Qual("math", "Min").Call(
				jen.ID("float64").Call(jen.ID("maxPASETOLifetime")),
				jen.ID("float64").Call(jen.ID("s").Dot("config").Dot("PASETO").Dot("Lifetime")),
			)),
			jen.ID("expiry").Op(":=").ID("now").Dot("Add").Call(jen.ID("lifetime")),
			jen.ID("jsonToken").Op(":=").ID("paseto").Dot("JSONToken").Valuesln(jen.ID("Audience").Op(":").Qual("strconv", "FormatUint").Call(
				jen.ID("client").Dot("BelongsToUser"),
				jen.Lit(10),
			), jen.ID("Subject").Op(":").Qual("strconv", "FormatUint").Call(
				jen.ID("client").Dot("BelongsToUser"),
				jen.Lit(10),
			), jen.ID("Jti").Op(":").ID("uuid").Dot("NewString").Call(), jen.ID("Issuer").Op(":").ID("s").Dot("config").Dot("PASETO").Dot("Issuer"), jen.ID("IssuedAt").Op(":").ID("now"), jen.ID("NotBefore").Op(":").ID("now"), jen.ID("Expiration").Op(":").ID("expiry")),
			jen.ID("jsonToken").Dot("Set").Call(
				jen.ID("pasetoDataKey"),
				jen.Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("sessionCtxData").Dot("ToBytes").Call()),
			),
			jen.Return().ID("jsonToken"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildPASETOResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("client").Op("*").ID("types").Dot("APIClient")).Params(jen.Op("*").ID("types").Dot("PASETOResponse"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("jsonToken").Op(":=").ID("s").Dot("buildPASETOToken").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("client"),
			),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("paseto").Dot("NewV2").Call().Dot("Encrypt").Call(
				jen.ID("s").Dot("config").Dot("PASETO").Dot("LocalModeKey"),
				jen.ID("jsonToken"),
				jen.Lit(""),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("s").Dot("logger"),
					jen.ID("span"),
					jen.Lit("encrypting PASETO"),
				))),
			jen.ID("tokenRes").Op(":=").Op("&").ID("types").Dot("PASETOResponse").Valuesln(jen.ID("Token").Op(":").ID("token"), jen.ID("ExpiresAt").Op(":").ID("jsonToken").Dot("Expiration").Dot("String").Call()),
			jen.Return().List(jen.ID("tokenRes"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CycleCookieSecretHandler rotates the cookie building secret with a new random secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("CycleCookieSecretHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("cycling cookie secret!")),
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
			jen.If(jen.Op("!").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("CanCycleCookieSecrets").Call()).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("invalid permissions")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeInvalidPermissionsResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("cookieManager").Op("=").ID("securecookie").Dot("New").Call(
				jen.ID("securecookie").Dot("GenerateRandomKey").Call(jen.ID("cookieSecretSize")),
				jen.Index().ID("byte").Call(jen.ID("s").Dot("config").Dot("Cookies").Dot("SigningKey")),
			),
			jen.ID("s").Dot("auditLog").Dot("LogCycleCookieSecretEvent").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	return code
}
