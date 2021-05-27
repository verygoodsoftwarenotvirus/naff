package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().ID("loginPrompt").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("loginPromptData").Struct(jen.ID("RedirectTo").ID("string")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildLoginView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.ID("contentData").Op(":=").Op("&").ID("loginPromptData").Valuesln(
					jen.ID("RedirectTo").Op(":").ID("pluckRedirectURL").Call(jen.ID("req"))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("loginPrompt"),
						jen.ID("nil"),
					),
					jen.ID("data").Op(":=").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("false"), jen.ID("Title").Op(":").Lit("BeginSession"), jen.ID("ContentData").Op(":").ID("contentData")),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("data"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("loginPrompt"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("contentData"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("usernameFormKey").Op("=").Lit("username"),
			jen.ID("passwordFormKey").Op("=").Lit("password"),
			jen.ID("totpTokenFormKey").Op("=").Lit("totpToken"),
			jen.ID("userIDFormKey").Op("=").Lit("userID"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("parseLoginInputFromForm checks a request for a login form, and returns the parsed login data if relevant."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("parseFormEncodedLoginRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("loginData").Op("*").ID("types").Dot("UserLoginInput"), jen.ID("redirectTo").ID("string")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("form"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(""))),
			jen.ID("loginData").Op("=").Op("&").ID("types").Dot("UserLoginInput").Valuesln(
				jen.ID("Username").Op(":").ID("form").Dot("Get").Call(jen.ID("usernameFormKey")), jen.ID("Password").Op(":").ID("form").Dot("Get").Call(jen.ID("passwordFormKey")), jen.ID("TOTPToken").Op(":").ID("form").Dot("Get").Call(jen.ID("totpTokenFormKey"))),
			jen.If(jen.ID("loginData").Dot("Username").Op("!=").Lit("").Op("&&").ID("loginData").Dot("Password").Op("!=").Lit("").Op("&&").ID("loginData").Dot("TOTPToken").Op("!=").Lit("")).Body(
				jen.Return().List(jen.ID("loginData"), jen.ID("form").Dot("Get").Call(jen.ID("redirectToQueryKey")))),
			jen.Return().List(jen.ID("nil"), jen.Lit("")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleLoginSubmission").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("loginInput"), jen.ID("redirectTo")).Op(":=").ID("s").Dot("parseFormEncodedLoginRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("loginInput").Op("==").ID("nil")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("no input found for login request")),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.If(jen.ID("redirectTo").Op("==").Lit("")).Body(
				jen.ID("redirectTo").Op("=").Lit("/")),
			jen.If(jen.Op("!").ID("s").Dot("useFakeData")).Body(
				jen.List(jen.ID("_"), jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot("authService").Dot("AuthenticateUser").Call(
					jen.ID("ctx"),
					jen.ID("loginInput"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("s").Dot("renderStringToResponse").Call(
						jen.ID("loginPrompt"),
						jen.ID("res"),
					),
					jen.Return(),
				),
				jen.Qual("net/http", "SetCookie").Call(
					jen.ID("res"),
					jen.ID("cookie"),
				),
				jen.ID("htmxRedirectTo").Call(
					jen.ID("res"),
					jen.ID("redirectTo"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleLogoutSubmission").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
					jen.Lit("no session context data attached to request"),
				),
				jen.Qual("net/http", "Redirect").Call(
					jen.ID("res"),
					jen.ID("req"),
					jen.Lit("/login"),
					jen.ID("unauthorizedRedirectResponseCode"),
				),
				jen.Return(),
			),
			jen.If(jen.Op("!").ID("s").Dot("useFakeData")).Body(
				jen.If(jen.ID("err").Op("=").ID("s").Dot("authService").Dot("LogoutUser").Call(
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
					jen.Return(),
				),
				jen.ID("htmxRedirectTo").Call(
					jen.ID("res"),
					jen.Lit("/"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("registrationPrompt").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("registrationComponent").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit(""),
				jen.ID("registrationPrompt"),
				jen.ID("nil"),
			),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID("nil"),
				jen.ID("res"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("registrationView").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("tmpl").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
				jen.ID("registrationPrompt"),
				jen.ID("nil"),
			),
			jen.ID("data").Op(":=").ID("pageData").Valuesln(
				jen.ID("IsLoggedIn").Op(":").ID("false"), jen.ID("Title").Op(":").Lit("Register"), jen.ID("ContentData").Op(":").ID("nil")),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID("data"),
				jen.ID("res"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("successfulRegistrationResponse").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("totpVerificationPrompt").Struct(
			jen.ID("TwoFactorQRCode").Qual("html/template", "URL"),
			jen.ID("UserID").ID("uint64"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("parseFormEncodedRegistrationRequest checks a request for a registration form, and returns the parsed login data if relevant."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("parseFormEncodedRegistrationRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("UserRegistrationInput")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("form"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("nil")),
			jen.ID("input").Op(":=").Op("&").ID("types").Dot("UserRegistrationInput").Valuesln(
				jen.ID("Username").Op(":").ID("form").Dot("Get").Call(jen.ID("usernameFormKey")), jen.ID("Password").Op(":").ID("form").Dot("Get").Call(jen.ID("passwordFormKey"))),
			jen.If(jen.ID("input").Dot("Username").Op("!=").Lit("").Op("&&").ID("input").Dot("Password").Op("!=").Lit("")).Body(
				jen.Return().ID("input")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleRegistrationSubmission").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("registrationInput").Op(":=").ID("s").Dot("parseFormEncodedRegistrationRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("registrationInput").Op("==").ID("nil")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("no input found for registration request")),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.Var().ID("ucr").Op("*").ID("types").Dot("UserCreationResponse"),
			jen.If(jen.Op("!").ID("s").Dot("useFakeData")).Body(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("ucr"), jen.ID("err")).Op("=").ID("s").Dot("usersService").Dot("RegisterUser").Call(
					jen.ID("ctx"),
					jen.ID("registrationInput"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("s").Dot("renderStringToResponse").Call(
						jen.ID("registrationPrompt"),
						jen.ID("res"),
					),
					jen.Return(),
				),
			).Else().Body(
				jen.ID("ucr").Op("=").Op("&").ID("types").Dot("UserCreationResponse").Valuesln(
					jen.ID("TwoFactorQRCode").Op(":").Lit(""))),
			jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
				jen.ID("ctx"),
				jen.Lit(""),
				jen.ID("successfulRegistrationResponse"),
				jen.ID("nil"),
			),
			jen.ID("tmplData").Op(":=").Op("&").ID("totpVerificationPrompt").Valuesln(
				jen.ID("TwoFactorQRCode").Op(":").Qual("html/template", "URL").Call(jen.ID("ucr").Dot("TwoFactorQRCode")), jen.ID("UserID").Op(":").ID("ucr").Dot("CreatedUserID")),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID("tmplData"),
				jen.ID("res"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("parseFormEncodedTOTPSecretVerificationRequest checks a request for a registration form, and returns the parsed input."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("parseFormEncodedTOTPSecretVerificationRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("form"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("nil")),
			jen.List(jen.ID("userID"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("form").Dot("Get").Call(jen.ID("userIDFormKey")),
				jen.Lit(10),
				jen.Lit(64),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("nil")),
			jen.ID("input").Op(":=").Op("&").ID("types").Dot("TOTPSecretVerificationInput").Valuesln(
				jen.ID("UserID").Op(":").ID("userID"), jen.ID("TOTPToken").Op(":").ID("form").Dot("Get").Call(jen.ID("totpTokenFormKey"))),
			jen.If(jen.ID("input").Dot("TOTPToken").Op("!=").Lit("").Op("&&").ID("input").Dot("UserID").Op("!=").Lit(0)).Body(
				jen.Return().ID("input")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("handleTOTPVerificationSubmission").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("verificationInput").Op(":=").ID("s").Dot("parseFormEncodedTOTPSecretVerificationRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("verificationInput").Op("==").ID("nil")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("no input found for registration request")),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.If(jen.ID("err").Op(":=").ID("s").Dot("usersService").Dot("VerifyUserTwoFactorSecret").Call(
				jen.ID("ctx"),
				jen.ID("verificationInput"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("rendering item viewer into dashboard"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("htmxRedirectTo").Call(
				jen.ID("res"),
				jen.Lit("/login"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	return code
}
