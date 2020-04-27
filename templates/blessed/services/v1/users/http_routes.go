package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("URIParamKey is used to refer to user IDs in router params."),
			jen.ID("URIParamKey").Equals().Lit("userID"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("this function tests that we have appropriate access to crypto/rand"),
		jen.Line(),
		jen.Func().ID("init").Params().Block(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("randString produces a random string."),
		jen.Line(),
		jen.Comment("https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/"),
		jen.Line(),
		jen.Func().ID("randString").Params().Params(jen.String(), jen.Error()).Block(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.Comment("Note that err == nil only if we read len(b) bytes."),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.EmptyString(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("validateCredentialChangeRequest takes a user's credentials and determines."),
		jen.Line(),
		jen.Comment("if they match what is on record."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("validateCredentialChangeRequest").Paramsln(
			constants.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.Listln(jen.ID("password"), jen.ID("totpToken")).String(),
		).Params(jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"), jen.ID("httpStatus").ID("int")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("validateCredentialChangeRequest")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.Line(),
			jen.Comment("fetch user data."),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("userDataManager").Dot("GetUser").Call(constants.CtxVar(), jen.ID("userID")),
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
			).Else().If(jen.Op("!").ID("valid")).Block(
				jen.ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("valid"), jen.ID("valid")).Dot("Error").Call(jen.Err(), jen.Lit("invalid attempt to cycle TOTP token")),
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusUnauthorized")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("user"), jen.Qual("net/http", "StatusOK")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is a handler for responding with a list of users."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
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
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching users for ListHandler route")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode response."),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("users")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler is our user creation route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CreateHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("in the event that we don't want new users to be able to sign up (a config setting)"),
				jen.Comment("just decline the request from the get-go"),
				jen.If(jen.Op("!").ID("s").Dot("userCreationEnabled")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Info").Call(jen.Lit("disallowing user creation")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusForbidden"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("fetch parsed input from request context."),
				jen.List(jen.ID("userInput"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(
					jen.ID("UserCreationMiddlewareCtxKey"),
				).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Info").Call(jen.Lit("valid input not attached to UsersService CreateHandler request")),
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
				),
				jen.Line(),
				jen.Comment("generate a two factor secret."),
				jen.List(jen.ID("input").Dot("TwoFactorSecret"), jen.Err()).Equals().ID("randString").Call(),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error generating TOTP secret")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
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
				jen.Comment("UserCreationResponse is a struct we can use to notify the user of."),
				jen.Comment("their two factor secret, but ideally just this once and then never again."),
				jen.ID("ucr").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserCreationResponse").Valuesln(
					jen.ID("ID").MapAssign().ID("user").Dot("ID"),
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"),
					jen.ID("PasswordLastChangedOn").MapAssign().ID("user").Dot("PasswordLastChangedOn"),
					jen.ID("CreatedOn").MapAssign().ID("user").Dot("CreatedOn"),
					jen.ID("UpdatedOn").MapAssign().ID("user").Dot("UpdatedOn"),
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
		),
		jen.Line(),
	)

	ret.Add(
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
					jen.Lit("todoservice"),
					jen.ID("username"),
					jen.ID("twoFactorSecret"),
					jen.Lit("todoService"),
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
			jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("data:image/jpeg;base64,%s"), jen.Qual("encoding/base64", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b").Dot("Bytes").Call())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler is our read route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ReadHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out who this is all for."),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("document it for posterity."),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Line(),
				jen.Comment("fetch user data."),
				jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot("userDataManager").Dot("GetUser").Call(constants.CtxVar(), jen.ID("userID")),
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
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("NewTOTPSecretHandler fetches a user, and issues them a new TOTP secret, after validating."),
		jen.Line(),
		jen.Comment("that information received from TOTPSecretRefreshInputContextMiddleware is valid."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("NewTOTPSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("NewTOTPSecretHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("check request context for parsed input."),
				jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.ID("TOTPSecretRefreshMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(),
					"TOTPSecretRefreshInput",
				)),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no input found on TOTP secret refresh request")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("also check for the user's ID."),
				jen.List(jen.ID("userID"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "UserIDKey")).Assert(jen.Uint64()),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no user ID attached to TOTP secret refresh request")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("make sure this is all on the up-and-up"),
				jen.List(jen.ID("user"), jen.ID("sc")).Assign().ID("s").Dot("validateCredentialChangeRequest").Callln(
					constants.CtxVar(),
					jen.ID("userID"),
					jen.ID("input").Dot("CurrentPassword"),
					jen.ID("input").Dot("TOTPToken"),
				),
				jen.Line(),
				jen.Comment("if the above function returns something other than 200, it means some error occurred."),
				jen.If(jen.ID("sc").DoesNotEqual().Qual("net/http", "StatusOK")).Block(
					jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.ID("sc")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("document who this is for."),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("Username")),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user"), jen.ID("user").Dot("ID")),
				jen.Line(),
				jen.Comment("set the two factor secret."),
				jen.List(jen.ID("tfs"), jen.Err()).Assign().ID("randString").Call(),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered generating random TOTP string")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.ID("user").Dot("TwoFactorSecret").Equals().ID("tfs"),
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
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdatePasswordHandler updates a user's password, after validating that information received."),
		jen.Line(),
		jen.Comment("from PasswordUpdateInputContextMiddleware is valid."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UpdatePasswordHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UpdatePasswordHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("check request context for parsed value."),
				jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("PasswordChangeMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "PasswordUpdateInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no input found on UpdatePasswordHandler request")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("check request context for user ID."),
				jen.List(jen.ID("userID"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "UserIDKey")).Assert(jen.Uint64()), jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no user ID attached to UpdatePasswordHandler request")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("make sure everything's on the up-and-up"),
				jen.List(jen.ID("user"), jen.ID("sc")).Assign().ID("s").Dot("validateCredentialChangeRequest").Callln(
					constants.CtxVar(),
					jen.ID("userID"),
					jen.ID("input").Dot("CurrentPassword"),
					jen.ID("input").Dot("TOTPToken"),
				),
				jen.Line(),
				jen.Comment("if the above function returns something other than 200, it means some error occurred."),
				jen.If(jen.ID("sc").DoesNotEqual().Qual("net/http", "StatusOK")).Block(
					jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.ID("sc")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("document who this is all for."),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("Username")),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user"), jen.ID("user").Dot("ID")),
				jen.Line(),
				jen.Comment("hash the new password."),
				jen.Var().Err().Error(),
				jen.List(jen.ID("user").Dot("HashedPassword"), jen.Err()).Equals().ID("s").Dot("authenticator").Dot("HashPassword").Call(constants.CtxVar(), jen.ID("input").Dot("NewPassword")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error hashing password")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("update the user."),
				jen.If(jen.Err().Equals().ID("s").Dot("userDataManager").Dot("UpdateUser").Call(constants.CtxVar(), jen.ID("user")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered updating user")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("we're all good."),
				utils.WriteXHeader(constants.ResponseVarName, "StatusAccepted"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler is a handler for archiving a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ArchiveHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out who this is for."),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Line(),
				jen.Comment("do the deed."),
				jen.If(jen.Err().Assign().ID("s").Dot("userDataManager").Dot("ArchiveUser").Call(constants.CtxVar(), jen.ID("userID")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("deleting user from database")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("inform the relatives."),
				jen.ID("s").Dot("userCounter").Dot("Decrement").Call(constants.CtxVar()),
				jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Archive")),
					jen.ID("Data").MapAssign().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().ID("userID")),
					jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName")),
				)),
				jen.Line(),
				jen.Comment("we're all good."),
				utils.WriteXHeader(constants.ResponseVarName, "StatusNoContent"),
			),
		),
		jen.Line(),
	)

	return ret
}
