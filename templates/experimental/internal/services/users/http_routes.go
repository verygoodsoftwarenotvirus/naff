package users

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
			jen.ID("UserIDURIParamKey").Op("=").Lit("userID"),
			jen.ID("totpIssuer").Op("=").Lit("todoService"),
			jen.ID("base64ImagePrefix").Op("=").Lit("data:image/jpeg;base64,"),
			jen.ID("minimumPasswordEntropy").Op("=").Lit(75),
			jen.ID("totpSecretSize").Op("=").Lit(64),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("validateCredentialChangeRequest takes a user's credentials and determines"),
		jen.Line(),
		jen.Func().Comment("if they match what is on record.").Params(jen.ID("s").Op("*").ID("service")).ID("validateCredentialChangeRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.List(jen.ID("password"), jen.ID("totpToken")).ID("string")).Params(jen.ID("user").Op("*").ID("types").Dot("User"), jen.ID("httpStatus").ID("int")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusNotFound"))).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("error encountered fetching user"),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusInternalServerError")),
			),
			jen.List(jen.ID("valid"), jen.ID("validationErr")).Op(":=").ID("s").Dot("authenticator").Dot("ValidateLogin").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("HashedPassword"),
				jen.ID("password"),
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.ID("totpToken"),
			),
			jen.If(jen.ID("validationErr").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("validation_error"),
					jen.ID("validationErr"),
				).Dot("Debug").Call(jen.Lit("error validating credentials")),
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusBadRequest")),
			).Else().If(jen.Op("!").ID("valid")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("valid"),
					jen.ID("valid"),
				).Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("invalid credentials"),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusUnauthorized")),
			),
			jen.Return().List(jen.ID("user"), jen.Qual("net/http", "StatusOK")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsernameSearchHandler is a handler for responding to username queries."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("UsernameSearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("query").Op(":=").ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.ID("types").Dot("SearchQueryKey")),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("WithValue").Call(
				jen.Lit("query"),
				jen.ID("query"),
			),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("users"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("SearchForUsersByUsername").Call(
				jen.ID("ctx"),
				jen.ID("query"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("searching for users"),
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
				jen.ID("users"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ListHandler is a handler for responding with a list of users."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("qf").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
			jen.List(jen.ID("users"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUsers").Call(
				jen.ID("ctx"),
				jen.ID("qf"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching users"),
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
				jen.ID("users"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("RegisterUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("registrationInput").Op("*").ID("types").Dot("UserRegistrationInput")).Params(jen.Op("*").ID("types").Dot("UserCreationResponse"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("registrationInput").Dot("Username"),
			),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("registrationInput").Dot("Username"),
			),
			jen.List(jen.ID("hp"), jen.ID("err")).Op(":=").ID("s").Dot("authenticator").Dot("HashPassword").Call(
				jen.ID("ctx"),
				jen.ID("registrationInput").Dot("Password"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("hashing password"),
				))),
			jen.ID("input").Op(":=").Op("&").ID("types").Dot("UserDataStoreCreationInput").Valuesln(jen.ID("Username").Op(":").ID("registrationInput").Dot("Username"), jen.ID("HashedPassword").Op(":").ID("hp"), jen.ID("TwoFactorSecret").Op(":").Lit("")),
			jen.If(jen.List(jen.ID("input").Dot("TwoFactorSecret"), jen.ID("err")).Op("=").ID("s").Dot("secretGenerator").Dot("GenerateBase32EncodedString").Call(
				jen.ID("ctx"),
				jen.ID("totpSecretSize"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("generating TOTP secret"),
				))),
			jen.List(jen.ID("user"), jen.ID("userCreationErr")).Op(":=").ID("s").Dot("userDataManager").Dot("CreateUser").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("userCreationErr").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("userCreationErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating user"),
				))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("user").Dot("ID"),
			),
			jen.ID("s").Dot("userCounter").Dot("Increment").Call(jen.ID("ctx")),
			jen.ID("ucr").Op(":=").Op("&").ID("types").Dot("UserCreationResponse").Valuesln(jen.ID("CreatedUserID").Op(":").ID("user").Dot("ID"), jen.ID("Username").Op(":").ID("user").Dot("Username"), jen.ID("CreatedOn").Op(":").ID("user").Dot("CreatedOn"), jen.ID("TwoFactorSecret").Op(":").ID("user").Dot("TwoFactorSecret"), jen.ID("TwoFactorQRCode").Op(":").ID("s").Dot("buildQRCode").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("Username"),
				jen.ID("user").Dot("TwoFactorSecret"),
			)),
			jen.Return().List(jen.ID("ucr"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler is our user creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.Op("!").ID("s").Dot("authSettings").Dot("EnableUserSignup")).Body(
				jen.ID("logger").Dot("Info").Call(jen.Lit("disallowing user creation")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("user creation is disabled"),
					jen.Qual("net/http", "StatusForbidden"),
				),
				jen.Return(),
			),
			jen.ID("userInput").Op(":=").ID("new").Call(jen.ID("types").Dot("UserRegistrationInput")),
			jen.If(jen.ID("err").Op(":=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("userInput"),
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
			jen.If(jen.ID("err").Op(":=").ID("userInput").Dot("ValidateWithContext").Call(
				jen.ID("ctx"),
				jen.ID("s").Dot("authSettings").Dot("MinimumUsernameLength"),
				jen.ID("s").Dot("authSettings").Dot("MinimumPasswordLength"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("ValidationErrorKey"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("provided input was invalid")),
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
				jen.ID("userInput").Dot("Username"),
			),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("userInput").Dot("Username"),
			),
			jen.If(jen.ID("err").Op(":=").Qual("github.com/wagslane/go-password-validator", "Validate").Call(
				jen.ID("userInput").Dot("Password"),
				jen.ID("minimumPasswordEntropy"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("password_validation_error"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("weak password provided to user creation route")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("password too weak"),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("ucr"), jen.ID("err")).Op(":=").ID("s").Dot("RegisterUser").Call(
				jen.ID("ctx"),
				jen.ID("userInput"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating user"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("ucr"),
				jen.Qual("net/http", "StatusCreated"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildQRCode builds a QR code for a given username and secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildQRCode").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("username"), jen.ID("twoFactorSecret")).ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("username"),
			),
			jen.ID("otpString").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("otpauth://totp/%s:%s?secret=%s&issuer=%s"),
				jen.ID("totpIssuer"),
				jen.ID("username"),
				jen.ID("twoFactorSecret"),
				jen.ID("totpIssuer"),
			),
			jen.List(jen.ID("qrCode"), jen.ID("err")).Op(":=").ID("qr").Dot("Encode").Call(
				jen.ID("otpString"),
				jen.ID("qr").Dot("L"),
				jen.ID("qr").Dot("Auto"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encoding OTP string"),
				),
				jen.Return().Lit(""),
			),
			jen.List(jen.ID("qrCode"), jen.ID("err")).Op("=").ID("barcode").Dot("Scale").Call(
				jen.ID("qrCode"),
				jen.Lit(256),
				jen.Lit(256),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scaling QR code"),
				),
				jen.Return().Lit(""),
			),
			jen.Var().Defs(
				jen.ID("b").Qual("bytes", "Buffer"),
			),
			jen.If(jen.ID("err").Op("=").Qual("image/png", "Encode").Call(
				jen.Op("&").ID("b"),
				jen.ID("qrCode"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encoding QR code to PNG"),
				),
				jen.Return().Lit(""),
			),
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s%s"),
				jen.ID("base64ImagePrefix"),
				jen.Qual("encoding/base64", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b").Dot("Bytes").Call()),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SelfHandler returns information about the user making the request."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("SelfHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("requester"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("requester"),
			),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUser").Call(
				jen.ID("ctx"),
				jen.ID("requester"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("no such user")),
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
					jen.Lit("fetching user"),
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
				jen.ID("user"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler is our read route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
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
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("no such user")),
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
					jen.Lit("fetching user from database"),
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
		jen.Var().Defs(
			jen.ID("errSecretAlreadyVerified").Op("=").Qual("errors", "New").Call(jen.Lit("secret already verified")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("VerifyUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("input").Dot("UserID"),
			),
			jen.List(jen.ID("user"), jen.ID("fetchUserErr")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUserWithUnverifiedTwoFactorSecret").Call(
				jen.ID("ctx"),
				jen.ID("input").Dot("UserID"),
			),
			jen.If(jen.ID("fetchUserErr").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("fetchUserErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching user with unverified two factor secret"),
				)),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("user").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("user").Dot("Username"),
			),
			jen.If(jen.ID("user").Dot("TwoFactorSecretVerifiedOn").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("two factor secret already verified")),
				jen.Return().ID("errSecretAlreadyVerified"),
			),
			jen.ID("totpValid").Op(":=").ID("totp").Dot("Validate").Call(
				jen.ID("input").Dot("TOTPToken"),
				jen.ID("user").Dot("TwoFactorSecret"),
			),
			jen.If(jen.Op("!").ID("totpValid")).Body(
				jen.Return().ID("authentication").Dot("ErrInvalidTOTPToken")),
			jen.If(jen.ID("updateUserErr").Op(":=").ID("s").Dot("userDataManager").Dot("MarkUserTwoFactorSecretAsVerified").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("ID"),
			), jen.ID("updateUserErr").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("updateUserErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("marking 2FA secret as validated"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TOTPSecretVerificationHandler accepts a TOTP token as input and returns 200 if the TOTP token"),
		jen.Line(),
		jen.Func().Comment("is validated by the user's TOTP secret.").Params(jen.ID("s").Op("*").ID("service")).ID("TOTPSecretVerificationHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("TOTPSecretVerificationInput")),
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
				).Dot("Debug").Call(jen.Lit("provided input was invalid")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("input").Dot("UserID"),
			),
			jen.If(jen.ID("twoFactorSecretVerificationError").Op(":=").ID("s").Dot("VerifyUserTwoFactorSecret").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			), jen.ID("twoFactorSecretVerificationError").Op("!=").ID("nil")).Body(
				jen.Switch().Body(
					jen.Case(jen.Qual("errors", "Is").Call(
						jen.ID("twoFactorSecretVerificationError"),
						jen.ID("authentication").Dot("ErrInvalidTOTPToken"),
					)).Body(
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeInvalidInputResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
						), jen.Return()),
					jen.Case(jen.Qual("errors", "Is").Call(
						jen.ID("twoFactorSecretVerificationError"),
						jen.ID("errSecretAlreadyVerified"),
					)).Body(
						jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
							jen.Lit("TOTP secret already verified"),
							jen.Qual("net/http", "StatusAlreadyReported"),
						), jen.Return()),
					jen.Default().Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("twoFactorSecretVerificationError"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("verifying user two factor secret"),
						), jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
							jen.ID("ctx"),
							jen.ID("res"),
						), jen.Return()),
				)),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewTOTPSecretHandler fetches a user, and issues them a new TOTP secret, after validating"),
		jen.Line(),
		jen.Func().Comment("that information received from TOTPSecretRefreshInputContextMiddleware is valid.").Params(jen.ID("s").Op("*").ID("service")).ID("NewTOTPSecretHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("TOTPSecretRefreshInput")),
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
				).Dot("Debug").Call(jen.Lit("provided input was invalid")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
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
			jen.List(jen.ID("user"), jen.ID("httpStatus")).Op(":=").ID("s").Dot("validateCredentialChangeRequest").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("input").Dot("CurrentPassword"),
				jen.ID("input").Dot("TOTPToken"),
			),
			jen.If(jen.ID("httpStatus").Op("!=").Qual("net/http", "StatusOK")).Body(
				jen.ID("res").Dot("WriteHeader").Call(jen.ID("httpStatus")),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("user").Dot("Username"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("user").Dot("ID"),
			),
			jen.List(jen.ID("tfs"), jen.ID("err")).Op(":=").ID("s").Dot("secretGenerator").Dot("GenerateBase32EncodedString").Call(
				jen.ID("ctx"),
				jen.ID("totpSecretSize"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("generating 2FA secret"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("user").Dot("TwoFactorSecret").Op("=").ID("tfs"),
			jen.ID("user").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("userDataManager").Dot("UpdateUser").Call(
				jen.ID("ctx"),
				jen.ID("user"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating 2FA secret"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("result").Op(":=").Op("&").ID("types").Dot("TOTPSecretRefreshResponse").Valuesln(jen.ID("TwoFactorSecret").Op(":").ID("user").Dot("TwoFactorSecret"), jen.ID("TwoFactorQRCode").Op(":").ID("s").Dot("buildQRCode").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("Username"),
				jen.ID("user").Dot("TwoFactorSecret"),
			)),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("result"),
				jen.Qual("net/http", "StatusAccepted"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdatePasswordHandler updates a user's password."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("UpdatePasswordHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("PasswordUpdateInput")),
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
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(
				jen.ID("ctx"),
				jen.ID("s").Dot("authSettings").Dot("MinimumPasswordLength"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("ValidationErrorKey"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("provided input was invalid")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
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
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.List(jen.ID("user"), jen.ID("httpStatus")).Op(":=").ID("s").Dot("validateCredentialChangeRequest").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("input").Dot("CurrentPassword"),
				jen.ID("input").Dot("TOTPToken"),
			),
			jen.If(jen.ID("httpStatus").Op("!=").Qual("net/http", "StatusOK")).Body(
				jen.ID("res").Dot("WriteHeader").Call(jen.ID("httpStatus")),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("user").Dot("Username"),
			),
			jen.If(jen.ID("err").Op("=").Qual("github.com/wagslane/go-password-validator", "Validate").Call(
				jen.ID("input").Dot("NewPassword"),
				jen.ID("minimumPasswordEntropy"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("password_validation_error"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("invalid password provided")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("new passwords is too weak!"),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.List(jen.ID("newPasswordHash"), jen.ID("err")).Op(":=").ID("s").Dot("authenticator").Dot("HashPassword").Call(
				jen.ID("ctx"),
				jen.ID("input").Dot("NewPassword"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("hashing password"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("userDataManager").Dot("UpdateUserPassword").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("ID"),
				jen.ID("newPasswordHash"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encountered updating user"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.Qual("net/http", "SetCookie").Call(
				jen.ID("res"),
				jen.Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("MaxAge").Op(":").Op("-").Lit(1)),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("stringPointer").Params(jen.ID("storageProviderPath").ID("string")).Params(jen.Op("*").ID("string")).Body(
			jen.Return().Op("&").ID("storageProviderPath")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AvatarUploadHandler updates a user's avatar."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("AvatarUploadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("session context data data extracted")),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUser").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching associated user"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("user").Dot("ID"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("retrieved user from database")),
			jen.List(jen.ID("img"), jen.ID("err")).Op(":=").ID("s").Dot("imageUploadProcessor").Dot("Process").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Lit("avatar"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("img").Op("==").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("processing provided avatar upload file"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeInvalidInputResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("internalPath").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("avatar_%d"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("file_size"),
				jen.ID("len").Call(jen.ID("img").Dot("Data")),
			).Dot("WithValue").Call(
				jen.Lit("internal_path"),
				jen.ID("internalPath"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("uploadManager").Dot("SaveFile").Call(
				jen.ID("ctx"),
				jen.ID("internalPath"),
				jen.ID("img").Dot("Data"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("saving provided avatar"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("user").Dot("AvatarSrc").Op("=").ID("stringPointer").Call(jen.ID("internalPath")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("userDataManager").Dot("UpdateUser").Call(
				jen.ID("ctx"),
				jen.ID("user"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating user info"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler is a handler for archiving a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
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
			jen.ID("err").Op(":=").ID("s").Dot("userDataManager").Dot("ArchiveUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
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
					jen.Lit("archiving user"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("userCounter").Dot("Decrement").Call(jen.ID("ctx")),
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
			jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetAuditLogEntriesForUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
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
					jen.Lit("fetching audit log entries for user"),
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
