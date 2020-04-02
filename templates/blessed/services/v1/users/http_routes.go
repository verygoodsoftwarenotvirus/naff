package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("URIParamKey is used to refer to user IDs in router params"),
			jen.ID("URIParamKey").Equals().Lit("userID"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("this function tests that we have appropriate access to crypto/rand"),
		jen.Line(),
		jen.Func().ID("init").Params().Block(
			jen.ID("b").Assign().ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("randString produces a random string"),
		jen.Line(),
		jen.Comment("https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/"),
		jen.Line(),
		jen.Func().ID("randString").Params().Params(jen.ID("string"), jen.ID("error")).Block(
			jen.ID("b").Assign().ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.Comment("Note that err == nil only if we read len(b) bytes."),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Lit(""), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("validateCredentialChangeRequest takes a user's credentials and determines"),
		jen.Line(),
		jen.Comment("if they match what is on record"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("validateCredentialChangeRequest").Paramsln(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
			jen.Listln(jen.ID("password"), jen.ID("totpToken")).ID("string"),
		).Params(jen.ID("user").Op("*").Qual(pkg.ModelsV1Package(), "User"), jen.ID("httpStatus").ID("int")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(utils.CtxVar(), jen.Lit("validateCredentialChangeRequest")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.Line(),
			jen.Comment("fetch user data"),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetUser").Call(utils.CtxVar(), jen.ID("userID")),
			jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusNotFound")),
			).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered fetching user")),
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusInternalServerError")),
			),
			jen.Line(),
			jen.Comment("validate login"),
			jen.List(jen.ID("valid"), jen.Err()).Assign().ID("s").Dot("authenticator").Dot("ValidateLogin").Callln(
				utils.CtxVar(),
				jen.ID("user").Dot("HashedPassword"),
				jen.ID("password"),
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.ID("totpToken"),
				jen.ID("user").Dot("Salt"),
			),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered generating random TOTP string")),
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusInternalServerError")),
			).Else().If(jen.Op("!").ID("valid")).Block(
				jen.ID("logger").Dot("WithValue").Call(jen.Lit("valid"), jen.ID("valid")).Dot("Error").Call(jen.Err(), jen.Lit("invalid attempt to cycle TOTP token")),
				jen.Return().List(jen.Nil(), jen.Qual("net/http", "StatusUnauthorized")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("user"), jen.Qual("net/http", "StatusOK")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is a handler for responding with a list of users"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ListHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("determine desired filter"),
				jen.ID("qf").Assign().Qual(pkg.ModelsV1Package(), "ExtractQueryFilter").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("fetch user data"),
				jen.List(jen.ID("users"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetUsers").Call(utils.CtxVar(), jen.ID("qf")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching users for ListHandler route")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode response"),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("users")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler is our user creation route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreateHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("in the event that we don't want new users to be able to sign up (a config setting)"),
				jen.Comment("just decline the request from the get-go"),
				jen.If(jen.Op("!").ID("s").Dot("userCreationEnabled")).Block(
					jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("disallowing user creation")),
					utils.WriteXHeader("res", "StatusForbidden"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("fetch parsed input from request context"),
				jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.ID("UserCreationMiddlewareCtxKey")).Assert(jen.Op("*").Qual(pkg.ModelsV1Package(), "UserInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("valid input not attached to UsersService CreateHandler request")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID("span"), jen.ID("input").Dot("Username")),
				jen.Line(),
				jen.Comment("NOTE: I feel comfortable letting username be in the logger, since"),
				jen.Comment("the logging statements below are only in the event of errors. If"),
				jen.Comment("and when that changes, this can/should be removed."),
				jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("username"), jen.ID("input").Dot("Username")),
				jen.Line(),
				jen.Comment("hash the password"),
				jen.List(jen.ID("hp"), jen.Err()).Assign().ID("s").Dot("authenticator").Dot("HashPassword").Call(utils.CtxVar(), jen.ID("input").Dot("Password")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("valid input not attached to request")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.ID("input").Dot("Password").Equals().ID("hp"),
				jen.Line(),
				jen.Comment("generate a two factor secret"),
				jen.List(jen.ID("input").Dot("TwoFactorSecret"), jen.Err()).Equals().ID("randString").Call(),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error generating TOTP secret")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("create the user"),
				jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("database").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("input")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.If(jen.Err().Op("==").Qual(pkg.DatabaseV1Package("client"), "ErrUserExists")).Block(
						jen.ID("logger").Dot("Info").Call(jen.Lit("duplicate username attempted")),
						utils.WriteXHeader("res", "StatusBadRequest"),
						jen.Return(),
					),
					jen.Line(),
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error creating user")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("UserCreationResponse is a struct we can use to notify the user of"),
				jen.Comment("their two factor secret, but ideally just this once and then never again."),
				jen.ID("ucr").Assign().VarPointer().Qual(pkg.ModelsV1Package(), "UserCreationResponse").Valuesln(
					jen.ID("ID").MapAssign().ID("user").Dot("ID"),
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"),
					jen.ID("PasswordLastChangedOn").MapAssign().ID("user").Dot("PasswordLastChangedOn"),
					jen.ID("CreatedOn").MapAssign().ID("user").Dot("CreatedOn"),
					jen.ID("UpdatedOn").MapAssign().ID("user").Dot("UpdatedOn"),
					jen.ID("ArchivedOn").MapAssign().ID("user").Dot("ArchivedOn"),
					jen.ID("TwoFactorQRCode").MapAssign().ID("s").Dot("buildQRCode").Call(utils.CtxVar(), jen.ID("user").Dot("Username"), jen.ID("user").Dot("TwoFactorSecret")),
				),
				jen.Line(),
				jen.Comment("notify the relevant parties"),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("user").Dot("ID")),
				jen.ID("s").Dot("userCounter").Dot("Increment").Call(utils.CtxVar()),
				jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").MapAssign().ID("string").Call(jen.Qual(pkg.ModelsV1Package(), "Create")),
					jen.ID("Data").MapAssign().ID("ucr"),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.ID("topicName")),
				)),
				jen.Line(),
				jen.Comment("encode and peace"),
				utils.WriteXHeader("res", "StatusCreated"),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("ucr")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildQRCode builds a QR code for a given username and secret"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("buildQRCode").Params(utils.CtxParam(), jen.List(jen.ID("username"), jen.ID("twoFactorSecret")).ID("string")).Params(jen.ID("string")).Block(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(utils.CtxVar(), jen.Lit("buildQRCode")),
			jen.Defer().ID("span").Dot("End").Call(),
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
				jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("trying to encode secret to qr code")),
				jen.Return().Lit(""),
			),
			jen.Line(),
			jen.Comment("scale the QR code so that it's not a PNG for ants"),
			jen.List(jen.ID("qrcode"), jen.Err()).Equals().Qual("github.com/boombuler/barcode", "Scale").Call(jen.ID("qrcode"), jen.Lit(256), jen.Lit(256)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("trying to enlarge qr code")),
				jen.Return().Lit(""),
			),
			jen.Line(),
			jen.Comment("encode the QR code to PNG"),
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.If(jen.Err().Equals().Qual("image/png", "Encode").Call(jen.VarPointer().ID("b"), jen.ID("qrcode")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("trying to encode qr code to png")),
				jen.Return().Lit(""),
			),
			jen.Line(),
			jen.Comment("base64 encode the image for easy HTML use"),
			jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("data:image/jpeg;base64,%s"), jen.Qual("encoding/base64", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b").Dot("Bytes").Call())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler is our read route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ReadHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("figure out who this is all for"),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("document it for posterity"),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("fetch user data"),
				jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetUser").Call(utils.CtxVar(), jen.ID("userID")),
				jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("no such user")),
					utils.WriteXHeader("res", "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching user from database")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode response and peace"),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("NewTOTPSecretHandler fetches a user, and issues them a new TOTP secret, after validating"),
		jen.Line(),
		jen.Comment("that information received from TOTPSecretRefreshInputContextMiddleware is valid"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("NewTOTPSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("NewTOTPSecretHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("check request context for parsed input"),
				jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID("req").Dot("Context").Call().Dot("Value").Call(jen.ID("TOTPSecretRefreshMiddlewareCtxKey")).Assert(jen.Op("*").Qual(pkg.ModelsV1Package(),
					"TOTPSecretRefreshInput",
				)),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("no input found on TOTP secret refresh request")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("also check for the user's ID"),
				jen.List(jen.ID("userID"), jen.ID("ok")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.Qual(pkg.ModelsV1Package(), "UserIDKey")).Assert(jen.ID("uint64")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("no user ID attached to TOTP secret refresh request")),
					utils.WriteXHeader("res", "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("make sure this is all on the up-and-up"),
				jen.List(jen.ID("user"), jen.ID("sc")).Assign().ID("s").Dot("validateCredentialChangeRequest").Callln(
					utils.CtxVar(),
					jen.ID("userID"),
					jen.ID("input").Dot("CurrentPassword"),
					jen.ID("input").Dot("TOTPToken"),
				),
				jen.Line(),
				jen.Comment("if the above function returns something other than 200, it means some error occurred"),
				jen.If(jen.ID("sc").DoesNotEqual().Qual("net/http", "StatusOK")).Block(
					jen.ID("res").Dot("WriteHeader").Call(jen.ID("sc")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("document who this is for"),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID("span"), jen.ID("user").Dot("Username")),
				jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user"), jen.ID("user").Dot("ID")),
				jen.Line(),
				jen.Comment("set the two factor secret"),
				jen.List(jen.ID("tfs"), jen.Err()).Assign().ID("randString").Call(),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered generating random TOTP string")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.ID("user").Dot("TwoFactorSecret").Equals().ID("tfs"),
				jen.Line(),
				jen.Comment("update the user in the database"),
				jen.If(jen.Err().Assign().ID("s").Dot("database").Dot("UpdateUser").Call(utils.CtxVar(), jen.ID("user")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered updating TOTP token")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("let the requester know we're all good"),
				utils.WriteXHeader("res", "StatusAccepted"),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.VarPointer().Qual(pkg.ModelsV1Package(), "TOTPSecretRefreshResponse").Values(jen.ID("TwoFactorSecret").MapAssign().ID("user").Dot("TwoFactorSecret"))), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdatePasswordHandler updates a user's password, after validating that information received"),
		jen.Line(),
		jen.Comment("from PasswordUpdateInputContextMiddleware is valid"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UpdatePasswordHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UpdatePasswordHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("check request context for parsed value"),
				jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.ID("PasswordChangeMiddlewareCtxKey")).Assert(jen.Op("*").Qual(pkg.ModelsV1Package(), "PasswordUpdateInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("no input found on UpdatePasswordHandler request")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("check request context for user ID"),
				jen.List(jen.ID("userID"), jen.ID("ok")).Assign().ID(utils.ContextVarName).Dot("Value").Call(jen.Qual(pkg.ModelsV1Package(), "UserIDKey")).Assert(jen.ID("uint64")), jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("no user ID attached to UpdatePasswordHandler request")),
					utils.WriteXHeader("res", "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("make sure everything's on the up-and-up"),
				jen.List(jen.ID("user"), jen.ID("sc")).Assign().ID("s").Dot("validateCredentialChangeRequest").Callln(
					utils.CtxVar(),
					jen.ID("userID"),
					jen.ID("input").Dot("CurrentPassword"),
					jen.ID("input").Dot("TOTPToken"),
				),
				jen.Line(),
				jen.Comment("if the above function returns something other than 200, it means some error occurred"),
				jen.If(jen.ID("sc").DoesNotEqual().Qual("net/http", "StatusOK")).Block(
					jen.ID("res").Dot("WriteHeader").Call(jen.ID("sc")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("document who this is all for"),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID("span"), jen.ID("user").Dot("Username")),
				jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user"), jen.ID("user").Dot("ID")),
				jen.Line(),
				jen.Comment("hash the new password"),
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("user").Dot("HashedPassword"), jen.Err()).Equals().ID("s").Dot("authenticator").Dot("HashPassword").Call(utils.CtxVar(), jen.ID("input").Dot("NewPassword")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error hashing password")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("update the user"),
				jen.If(jen.Err().Equals().ID("s").Dot("database").Dot("UpdateUser").Call(utils.CtxVar(), jen.ID("user")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered updating user")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("we're all good"),
				utils.WriteXHeader("res", "StatusAccepted"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler is a handler for archiving a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ArchiveHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("figure out who this is for"),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.Qual(pkg.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("do the deed"),
				jen.If(jen.Err().Assign().ID("s").Dot("database").Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("userID")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("deleting user from database")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("inform the relatives"),
				jen.ID("s").Dot("userCounter").Dot("Decrement").Call(utils.CtxVar()),
				jen.ID("s").Dot("reporter").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").MapAssign().ID("string").Call(jen.Qual(pkg.ModelsV1Package(), "Archive")),
					jen.ID("Data").MapAssign().Qual(pkg.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().ID("userID")),
					jen.ID("Topics").MapAssign().Index().ID("string").Values(jen.ID("topicName")),
				)),
				jen.Line(),
				jen.Comment("we're all good"),
				utils.WriteXHeader("res", "StatusNoContent"),
			),
		),
		jen.Line(),
	)
	return ret
}
