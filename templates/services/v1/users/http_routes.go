package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildUsersHTTPRoutesVarDeclarations(proj)...)
	code.Add(buildUsersHTTPRoutesValidateCredentialChangeRequest(proj)...)
	code.Add(buildUsersHTTPRoutesListHandler(proj)...)
	code.Add(buildUsersHTTPRoutesCreateHandler(proj)...)
	code.Add(buildUsersHTTPRoutesBuildQRCode(proj)...)
	code.Add(buildUsersHTTPRoutesReadHandler(proj)...)
	code.Add(buildUsersHTTPRoutesTOTPSecretVerificationHandler(proj)...)
	code.Add(buildUsersHTTPRoutesNewTOTPSecretHandler(proj)...)
	code.Add(buildUsersHTTPRoutesUpdatePasswordHandler(proj)...)
	code.Add(buildUsersHTTPRoutesArchiveHandler(proj)...)

	return code
}

func buildUsersHTTPRoutesVarDeclarations(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("URIParamKey is used to refer to user IDs in router params."),
			jen.ID("URIParamKey").Equals().Lit("userID"),
			jen.Line(),
			jen.ID("totpIssuer").Equals().Litf("%sService", proj.Name.UnexportedVarName()),
			jen.ID("base64ImagePrefix").Equals().Lit("data:image/jpeg;base64,"),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesValidateCredentialChangeRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("validateCredentialChangeRequest takes a user's credentials and determines"),
		jen.Line(),
		jen.Comment("if they match what is on record."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("validateCredentialChangeRequest").Paramsln(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.Listln(jen.ID("password"), jen.ID("totpToken")).String(),
		).Params(jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"), jen.ID("httpStatus").ID("int")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("validateCredentialChangeRequest")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)),
			jen.Line(),
			jen.Comment("fetch user data."),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("userDataManager").Dot("GetUser").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusNotFound")),
			).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered fetching user")),
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusInternalServerError")),
			),
			jen.Line(),
			jen.Comment("validate login."),
			jen.List(jen.ID("valid"), jen.Err()).Assign().ID("s").Dot("authenticator").Dot("ValidateLogin").Callln(
				constants.CtxVar(),
				jen.ID("user").Dot("HashedPassword"),
				jen.ID("password"),
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.ID("totpToken"),
				jen.ID("user").Dot("Salt"),
			),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered generating random TOTP string")),
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusInternalServerError")),
			).Else().If(jen.Not().ID("valid")).Block(
				jen.ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("valid"), jen.ID("valid")).Dot("Error").Call(jen.Err(), jen.Lit("invalid attempt to cycle TOTP token")),
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusUnauthorized")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("user"), jen.Qual("net/http", "StatusOK")),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesListHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ListHandler is a handler for responding with a list of users."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ListHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ListHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("determine desired filter."),
			jen.ID("qf").Assign().Qual(proj.ModelsV1Package(), "ExtractQueryFilter").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("fetch user data."),
			jen.List(jen.ID("users"), jen.Err()).Assign().ID("s").Dot("userDataManager").Dot("GetUsers").Call(constants.CtxVar(), jen.ID("qf")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching users for ListHandler route")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("encode response."),
			jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("users")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesCreateHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateHandler is our user creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreateHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CreateHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("in the event that we don't want new users to be able to sign up (a config setting)"),
			jen.Comment("just decline the request from the get-go"),
			jen.If(jen.Not().ID("s").Dot("userCreationEnabled")).Block(
				jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("disallowing user creation")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusForbidden"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("fetch parsed input from request context."),
			jen.List(jen.ID("userInput"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(
				jen.ID("userCreationMiddlewareCtxKey"),
			).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput")),
			jen.If(jen.Not().ID("ok")).Block(
				jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("valid input not attached to UsersService CreateHandler request")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
				jen.Return(),
			),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userInput").Dot("Username")),
			jen.Line(),
			jen.Comment("NOTE: I feel comfortable letting username be in the logger, since"),
			jen.Comment("the logging statements below are only in the event of errors. If"),
			jen.Comment("and when that changes, this can/should be removed."),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("username"), jen.ID("userInput").Dot("Username")),
			jen.Line(),
			jen.Comment("hash the password."),
			jen.List(jen.ID("hp"), jen.Err()).Assign().ID("s").Dot("authenticator").Dot("HashPassword").Call(constants.CtxVar(), jen.ID("userInput").Dot("Password")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("valid input not attached to request")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.ID("input").Assign().Qual(proj.ModelsV1Package(), "UserDatabaseCreationInput").Valuesln(
				jen.ID("Username").MapAssign().ID("userInput").Dot("Username"),
				jen.ID("HashedPassword").MapAssign().ID("hp"),
				jen.ID("TwoFactorSecret").MapAssign().EmptyString(),
				jen.ID("Salt").MapAssign().Index().Byte().Values(),
			),
			jen.Line(),
			jen.Comment("generate a two factor secret."),
			jen.List(jen.ID("input").Dot("TwoFactorSecret"), jen.Err()).Equals().ID("s").Dot("secretGenerator").Dot("GenerateTwoFactorSecret").Call(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error generating TOTP secret")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("generate a salt."),
			jen.List(jen.ID("input").Dot("Salt"), jen.Err()).Equals().ID("s").Dot("secretGenerator").Dot("GenerateSalt").Call(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error generating salt")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("create the user."),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("userDataManager").Dot("CreateUser").Call(constants.CtxVar(), jen.ID("input")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual(proj.DatabaseV1Package("client"), "ErrUserExists")).Block(
					jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("duplicate username attempted")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error creating user")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("UserCreationResponse is a struct we can use to notify the user of"),
			jen.Comment("their two factor secret, but ideally just this once and then never again."),
			jen.ID("ucr").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserCreationResponse").Valuesln(
				jen.ID("ID").MapAssign().ID("user").Dot("ID"),
				jen.ID("Username").MapAssign().ID("user").Dot("Username"),
				jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"),
				jen.ID("PasswordLastChangedOn").MapAssign().ID("user").Dot("PasswordLastChangedOn"),
				jen.ID("CreatedOn").MapAssign().ID("user").Dot("CreatedOn"),
				jen.ID("LastUpdatedOn").MapAssign().ID("user").Dot("LastUpdatedOn"),
				jen.ID("ArchivedOn").MapAssign().ID("user").Dot("ArchivedOn"),
				jen.ID("TwoFactorQRCode").MapAssign().ID("s").Dot("buildQRCode").Call(constants.CtxVar(), jen.ID("user").Dot("Username"), jen.ID("user").Dot("TwoFactorSecret")),
			),
			jen.Line(),
			jen.Comment("notify the relevant parties."),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("ID")),
			jen.ID("s").Dot("userCounter").Dot("Increment").Call(constants.CtxVar()),
			jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
				jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Create")),
				jen.ID("Data").MapAssign().ID("ucr"),
				jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName")),
			)),
			jen.Line(),
			jen.Comment("encode and peace."),
			utils.WriteXHeader(constants.ResponseVarName, "StatusCreated"),
			jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("ucr")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesBuildQRCode(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildQRCode builds a QR code for a given username and secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("buildQRCode").Params(constants.CtxParam(), jen.List(jen.ID("username"), jen.ID("twoFactorSecret")).String()).Params(jen.String()).Block(
			jen.List(jen.Underscore(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("buildQRCode")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Comment("encode two factor secret as authenticator-friendly QR code"),
			jen.List(jen.ID("qrcode"), jen.Err()).Assign().Qual("github.com/boombuler/barcode/qr", "Encode").Callln(
				jen.Comment(`"otpauth://totp/{{ .Issuer }}:{{ .Username }}?secret={{ .Secret }}&issuer={{ .Issuer }}"`),
				jen.Qual("fmt", "Sprintf").Callln(
					jen.Lit("otpauth://totp/%s:%s?secret=%s&issuer=%s"),
					jen.ID("totpIssuer"),
					jen.ID("username"),
					jen.ID("twoFactorSecret"),
					jen.ID("totpIssuer"),
				),
				jen.Qual("github.com/boombuler/barcode/qr", "L"),
				jen.Qual("github.com/boombuler/barcode/qr", "Auto"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("trying to encode secret to qr code")),
				jen.Return().EmptyString(),
			),
			jen.Line(),
			jen.Comment("scale the QR code so that it's not a PNG for ants."),
			jen.List(jen.ID("qrcode"), jen.Err()).Equals().Qual("github.com/boombuler/barcode", "Scale").Call(jen.ID("qrcode"), jen.Lit(256), jen.Lit(256)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("trying to enlarge qr code")),
				jen.Return().EmptyString(),
			),
			jen.Line(),
			jen.Comment("encode the QR code to PNG."),
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.If(jen.Err().Equals().Qual("image/png", "Encode").Call(jen.AddressOf().ID("b"), jen.ID("qrcode")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("trying to encode qr code to png")),
				jen.Return().EmptyString(),
			),
			jen.Line(),
			jen.Comment("base64 encode the image for easy HTML use."),
			jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("%s%s"), jen.ID("base64ImagePrefix"), jen.Qual("encoding/base64", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b").Dot("Bytes").Call())),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesReadHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ReadHandler is our read route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ReadHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ReadHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("figure out who this is all for."),
			jen.ID(constants.UserIDVarName).Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)),
			jen.Line(),
			jen.Comment("document it for posterity."),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Line(),
			jen.Comment("fetch user data."),
			jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot("userDataManager").Dot("GetUser").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no such user")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
				jen.Return(),
			).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching user from database")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("encode response and peace."),
			jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesTOTPSecretVerificationHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("TOTPSecretVerificationHandler accepts a TOTP token as input and returns 200 if the TOTP token"),
		jen.Line(),
		jen.Comment("is validated by the user's TOTP secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("TOTPSecretVerificationHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("TOTPSecretVerificationHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("check request context for parsed input."),
			jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.ID("totpSecretVerificationMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(),
				"TOTPSecretVerificationInput",
			)),
			jen.If(jen.Not().ID("ok").Or().ID("input").IsEqualTo().Nil()).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no input found on TOTP secret refresh request")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
				jen.Return(),
			),
			jen.Line(),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("userDataManager").Dot("GetUserWithUnverifiedTwoFactorSecret").Call(
				constants.CtxVar(),
				jen.ID("input").Dot(constants.UserIDFieldName),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching user")),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("ID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("Username")),
			jen.Line(),
			jen.If(jen.ID("user").Dot("TwoFactorSecretVerifiedOn").DoesNotEqual().Nil()).Block(
				jen.Comment("I suppose if this happens too many times, we'll want to keep track of that"),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAlreadyReported")),
				jen.Return(),
			),
			jen.Line(),
			jen.If(
				jen.Qual("github.com/pquerna/otp/totp", "Validate").Call(
					jen.ID("input").Dot("TOTPToken"),
					jen.ID("user").Dot("TwoFactorSecret"),
				),
			).Block(
				jen.If(
					jen.ID("updateUserErr").Assign().ID("s").Dot("userDataManager").Dot("VerifyUserTwoFactorSecret").Call(
						constants.CtxVar(),
						jen.ID("user").Dot("ID"),
					),
					jen.ID("updateUserErr").DoesNotEqual().Nil(),
				).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.ID("updateUserErr"), jen.Lit("updating user to indicate their 2FA secret is validated")),
					jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
			).Else().Block(
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesNewTOTPSecretHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("NewTOTPSecretHandler fetches a user, and issues them a new TOTP secret, after validating"),
		jen.Line(),
		jen.Comment("that information received from TOTPSecretRefreshInputContextMiddleware is valid."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("NewTOTPSecretHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("NewTOTPSecretHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("check request context for parsed input."),
			jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.ID("totpSecretRefreshMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(),
				"TOTPSecretRefreshInput",
			)),
			jen.If(jen.Not().ID("ok")).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no input found on TOTP secret refresh request")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("also check for the user's ID."),
			jen.List(jen.ID("si"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "SessionInfoKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "SessionInfo")),
			jen.If(jen.Not().ID("ok").Or().ID("si").IsEqualTo().Nil()).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no user ID attached to TOTP secret refresh request")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("make sure this is all on the up-and-up"),
			jen.List(jen.ID("user"), jen.ID("httpStatus")).Assign().ID("s").Dot("validateCredentialChangeRequest").Callln(
				constants.CtxVar(),
				jen.ID("si").Dot(constants.UserIDFieldName),
				jen.ID("input").Dot("CurrentPassword"),
				jen.ID("input").Dot("TOTPToken"),
			),
			jen.Line(),
			jen.Comment("if the above function returns something other than 200, it means some error occurred."),
			jen.If(jen.ID("httpStatus").DoesNotEqual().Qual("net/http", "StatusOK")).Block(
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.ID("httpStatus")),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("document who this is for."),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("si").Dot(constants.UserIDFieldName)),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("Username")),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user"), jen.ID("user").Dot("ID")),
			jen.Line(),
			jen.Comment("set the two factor secret."),
			jen.List(jen.ID("tfs"), jen.Err()).Assign().ID("s").Dot("secretGenerator").Dot("GenerateTwoFactorSecret").Call(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered generating random TOTP string")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.ID("user").Dot("TwoFactorSecret").Equals().ID("tfs"),
			jen.ID("user").Dot("TwoFactorSecretVerifiedOn").Equals().ID("nil"),
			jen.Line(),
			jen.Comment("update the user in the database."),
			jen.If(jen.Err().Assign().ID("s").Dot("userDataManager").Dot("UpdateUser").Call(constants.CtxVar(), jen.ID("user")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered updating TOTP token")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("let the requester know we're all good."),
			utils.WriteXHeader(constants.ResponseVarName, "StatusAccepted"),
			jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.AddressOf().Qual(proj.ModelsV1Package(), "TOTPSecretRefreshResponse").Values(jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"))), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesUpdatePasswordHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdatePasswordHandler updates a user's password, after validating that information received"),
		jen.Line(),
		jen.Comment("from PasswordUpdateInputContextMiddleware is valid."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UpdatePasswordHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UpdatePasswordHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("check request context for parsed value."),
			jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("passwordChangeMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "PasswordUpdateInput")),
			jen.If(jen.Not().ID("ok")).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no input found on UpdatePasswordHandler request")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("check request context for user ID."),
			jen.List(jen.ID("si"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "SessionInfoKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "SessionInfo")),
			jen.If(jen.Not().ID("ok").Or().ID("si").IsEqualTo().Nil()).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no user ID attached to UpdatePasswordHandler request")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("determine relevant user ID."),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("si").Dot(constants.UserIDFieldName)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("si").Dot(constants.UserIDFieldName)),
			jen.Line(),
			jen.Comment("make sure everything's on the up-and-up"),
			jen.List(jen.ID("user"), jen.ID("httpStatus")).Assign().ID("s").Dot("validateCredentialChangeRequest").Callln(
				constants.CtxVar(),
				jen.ID("si").Dot(constants.UserIDFieldName),
				jen.ID("input").Dot("CurrentPassword"),
				jen.ID("input").Dot("TOTPToken"),
			),
			jen.Line(),
			jen.Comment("if the above function returns something other than 200, it means some error occurred."),
			jen.If(jen.ID("httpStatus").DoesNotEqual().Qual("net/http", "StatusOK")).Block(
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.ID("httpStatus")),
				jen.Return(),
			),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("Username")),
			jen.Line(),
			jen.Comment("hash the new password."),
			jen.List(jen.ID("newPasswordHash"), jen.Err()).Assign().ID("s").Dot("authenticator").Dot("HashPassword").Call(constants.CtxVar(), jen.ID("input").Dot("NewPassword")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error hashing password")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("update the user."),
			jen.If(jen.Err().Equals().ID("s").Dot("userDataManager").Dot("UpdateUserPassword").Call(
				constants.CtxVar(),
				jen.ID("user").Dot("ID"),
				jen.ID("newPasswordHash"),
			), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered updating user")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("we're all good."),
			utils.WriteXHeader(constants.ResponseVarName, "StatusAccepted"),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersHTTPRoutesArchiveHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveHandler is a handler for archiving a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ArchiveHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ArchiveHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("figure out who this is for."),
			jen.ID(constants.UserIDVarName).Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Line(),
			jen.Comment("do the deed."),
			jen.If(jen.Err().Assign().ID("s").Dot("userDataManager").Dot("ArchiveUser").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("deleting user from database")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
				jen.Return(),
			),
			jen.Line(),
			jen.Comment("inform the relatives."),
			jen.ID("s").Dot("userCounter").Dot("Decrement").Call(constants.CtxVar()),
			jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
				jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Archive")),
				jen.ID("Data").MapAssign().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().ID(constants.UserIDVarName)),
				jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName")),
			)),
			jen.Line(),
			jen.Comment("we're all good."),
			utils.WriteXHeader(constants.ResponseVarName, "StatusNoContent"),
		),
		jen.Line(),
	}

	return lines
}
