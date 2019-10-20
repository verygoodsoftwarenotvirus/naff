package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func httpRoutesDotGo() *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("URIParamKey").Op("=").Lit("userID"),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Comment("// this function tests that we have appropriate access to crypto/rand").ID("init").Params().Block(
			jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachUsernameToSpan provides a consistent way to attach a username to a span"),
		jen.Line(),
		jen.Func().ID("attachUsernameToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("username").ID("string")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot(
					"AddAttributes",
				).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("username"), jen.ID("username"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachUserIDToSpan provides a consistent way to attach a user ID to a span"),
		jen.Line(),
		jen.Func().ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot(
					"AddAttributes",
				).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("randString produces a random string"),
		jen.Line(),
		jen.Func().Comment("// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/").ID("randString").Params().Params(jen.ID("string"), jen.ID("error")).Block(
			jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.Lit(""), jen.ID("err")),
			),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot(
				"EncodeToString",
			).Call(jen.ID("b")), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("validateCredentialChangeRequest takes a user's credentials and determines"),
		jen.Line(),
		jen.Comment("if they match what is on record"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("validateCredentialChangeRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.List(jen.ID("password"), jen.ID("totpToken")).ID("string")).Params(jen.ID("user").Op("*").ID("models").Dot(
			"User",
		),
			jen.ID("httpStatus").ID("int")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("validateCredentialChangeRequest")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithValue",
			).Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("Database").Dot(
				"GetUser",
			).Call(jen.ID("ctx"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusNotFound")),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered fetching user")),
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusInternalServerError")),
			),
			jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("s").Dot(
				"authenticator",
			).Dot(
				"ValidateLogin",
			).Call(jen.ID("ctx"), jen.ID("user").Dot(
				"HashedPassword",
			),
				jen.ID("password"), jen.ID("user").Dot(
					"TwoFactorSecret",
				),
				jen.ID("totpToken"), jen.ID("user").Dot(
					"Salt",
				)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered generating random TOTP string")),
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusInternalServerError")),
			).Else().If(jen.Op("!").ID("valid")).Block(
				jen.ID("logger").Dot(
					"WithValue",
				).Call(jen.Lit("valid"), jen.ID("valid")).Dot("Error").Call(jen.ID("err"), jen.Lit("invalid attempt to cycle TOTP token")),
				jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "StatusUnauthorized")),
			),
			jen.Return().List(jen.ID("user"), jen.Qual("net/http", "StatusOK")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is a handler for responding with a list of users"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("ListHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("qf").Op(":=").ID("models").Dot(
					"ExtractQueryFilter",
				).Call(jen.ID("req")),
				jen.List(jen.ID("users"), jen.ID("err")).Op(":=").ID("s").Dot("Database").Dot(
					"GetUsers",
				).Call(jen.ID("ctx"), jen.ID("qf")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching users for ListHandler route")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.If(jen.ID("err").Op("=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"EncodeResponse",
				).Call(jen.ID("res"), jen.ID("users")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler is our user creation route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("CreateHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.If(jen.Op("!").ID("s").Dot(
					"userCreationEnabled",
				)).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Info",
					).Call(jen.Lit("disallowing user creation")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusForbidden")),
					jen.Return(),
				),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot(
					"Value",
				).Call(jen.ID("UserCreationMiddlewareCtxKey")).Assert(jen.Op("*").ID("models").Dot(
					"UserInput",
				)),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Info",
					).Call(jen.Lit("valid input not attached to UsersService CreateHandler request")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("input").Dot(
					"Username",
				)),
				jen.ID("logger").Op(":=").ID("s").Dot(
					"logger",
				).Dot(
					"WithValue",
				).Call(jen.Lit("username"), jen.ID("input").Dot(
					"Username",
				)),
				jen.List(jen.ID("hp"), jen.ID("err")).Op(":=").ID("s").Dot(
					"authenticator",
				).Dot(
					"HashPassword",
				).Call(jen.ID("ctx"), jen.ID("input").Dot(
					"Password",
				)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("valid input not attached to request")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("input").Dot(
					"Password",
				).Op("=").ID("hp"),
				jen.List(jen.ID("input").Dot(
					"TwoFactorSecret",
				),
					jen.ID("err")).Op("=").ID("randString").Call(),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error generating TOTP secret")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("Database").Dot(
					"CreateUser",
				).Call(jen.ID("ctx"), jen.ID("input")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.If(jen.ID("err").Op("==").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client", "ErrUserExists")).Block(
						jen.ID("logger").Dot(
							"Info",
						).Call(jen.Lit("duplicate username attempted")),
						jen.ID("res").Dot(
							"WriteHeader",
						).Call(jen.Qual("net/http", "StatusBadRequest")),
						jen.Return(),
					),
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error creating user")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("ucr").Op(":=").Op("&").ID("models").Dot(
					"UserCreationResponse",
				).Valuesln(
	jen.ID("ID").Op(":").ID("user").Dot(
					"ID",
				),
					jen.ID("Username").Op(":").ID("user").Dot(
						"Username",
					),
					jen.ID("TwoFactorSecret").Op(":").ID("user").Dot(
						"TwoFactorSecret",
					),
					jen.ID("PasswordLastChangedOn").Op(":").ID("user").Dot(
						"PasswordLastChangedOn",
					),
					jen.ID("CreatedOn").Op(":").ID("user").Dot(
						"CreatedOn",
					),
					jen.ID("UpdatedOn").Op(":").ID("user").Dot(
						"UpdatedOn",
					),
					jen.ID("ArchivedOn").Op(":").ID("user").Dot(
						"ArchivedOn",
					),
					jen.ID("TwoFactorQRCode").Op(":").ID("s").Dot(
						"buildQRCode",
					).Call(jen.ID("ctx"), jen.ID("user").Dot(
						"Username",
					),
						jen.ID("user").Dot(
							"TwoFactorSecret",
						))),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("user").Dot(
					"ID",
				)),
				jen.ID("s").Dot(
					"userCounter",
				).Dot(
					"Increment",
				).Call(jen.ID("ctx")),
				jen.ID("s").Dot(
					"reporter",
				).Dot(
					"Report",
				).Call(jen.ID("newsman").Dot(
					"Event",
				).Valuesln(
	jen.ID("EventType").Op(":").ID("string").Call(jen.ID("models").Dot(
					"Create",
				)), jen.ID("Data").Op(":").ID("ucr"), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(
	jen.ID("topicName")))),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusCreated")),
				jen.If(jen.ID("err").Op("=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"EncodeResponse",
				).Call(jen.ID("res"), jen.ID("ucr")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildQRCode builds a QR code for a given username and secret"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("buildQRCode").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("username"), jen.ID("twoFactorSecret")).ID("string")).Params(jen.ID("string")).Block(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("buildQRCode")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("qrcode"), jen.ID("err")).Op(":=").ID("qr").Dot(
				"Encode",
			).Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("otpauth://totp/%s:%s?secret=%s&issuer=%s"), jen.Lit("todoservice"), jen.ID("username"), jen.ID("twoFactorSecret"), jen.Lit("todoService")), jen.ID("qr").Dot(
				"L",
			),
				jen.ID("qr").Dot(
					"Auto",
				)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot("Error").Call(jen.ID("err"), jen.Lit("trying to encode secret to qr code")),
				jen.Return().Lit(""),
			),
			jen.List(jen.ID("qrcode"), jen.ID("err")).Op("=").ID("barcode").Dot(
				"Scale",
			).Call(jen.ID("qrcode"), jen.Lit(256), jen.Lit(256)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot("Error").Call(jen.ID("err"), jen.Lit("trying to enlarge qr code")),
				jen.Return().Lit(""),
			),

			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.If(jen.ID("err").Op("=").Qual("image/png", "Encode").Call(jen.Op("&").ID("b"), jen.ID("qrcode")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot("Error").Call(jen.ID("err"), jen.Lit("trying to encode qr code to png")),
				jen.Return().Lit(""),
			),
			jen.Return().Qual("fmt", "Sprintf").Call(jen.Lit("data:image/jpeg;base64,%s"), jen.Qual("encoding/base64", "StdEncoding").Dot(
				"EncodeToString",
			).Call(jen.ID("b").Dot(
				"Bytes",
			).Call())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler is our read route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("ReadHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("userID").Op(":=").ID("s").Dot(
					"userIDFetcher",
				).Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot(
					"logger",
				).Dot(
					"WithValue",
				).Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("Database").Dot(
					"GetUser",
				).Call(jen.ID("ctx"), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot(
						"Debug",
					).Call(jen.Lit("no such user")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusNotFound")),
					jen.Return(),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error fetching user from database")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.If(jen.ID("err").Op("=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"EncodeResponse",
				).Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
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
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("NewTOTPSecretHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("req").Dot(
					"Context",
				).Call().Dot(
					"Value",
				).Call(jen.ID("TOTPSecretRefreshMiddlewareCtxKey")).Assert(jen.Op("*").ID("models").Dot(
					"TOTPSecretRefreshInput",
				)),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Debug",
					).Call(jen.Lit("no input found on TOTP secret refresh request")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.List(jen.ID("userID"), jen.ID("ok")).Op(":=").ID("ctx").Dot(
					"Value",
				).Call(jen.ID("models").Dot(
					"UserIDKey",
				)).Assert(jen.ID("uint64")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Debug",
					).Call(jen.Lit("no user ID attached to TOTP secret refresh request")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.List(jen.ID("user"), jen.ID("sc")).Op(":=").ID("s").Dot(
					"validateCredentialChangeRequest",
				).Call(jen.ID("ctx"), jen.ID("userID"), jen.ID("input").Dot(
					"CurrentPassword",
				),
					jen.ID("input").Dot(
						"TOTPToken",
					)),
				jen.If(jen.ID("sc").Op("!=").Qual("net/http", "StatusOK")).Block(
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.ID("sc")),
					jen.Return(),
				),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("user").Dot(
					"Username",
				)),
				jen.ID("logger").Op(":=").ID("s").Dot(
					"logger",
				).Dot(
					"WithValue",
				).Call(jen.Lit("user"), jen.ID("user").Dot(
					"ID",
				)),
				jen.List(jen.ID("tfs"), jen.ID("err")).Op(":=").ID("randString").Call(),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered generating random TOTP string")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("user").Dot(
					"TwoFactorSecret",
				).Op("=").ID("tfs"),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot("Database").Dot(
					"UpdateUser",
				).Call(jen.ID("ctx"), jen.ID("user")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered updating TOTP token")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusAccepted")),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"EncodeResponse",
				).Call(jen.ID("res"), jen.Op("&").ID("models").Dot(
					"TOTPSecretRefreshResponse",
				).Valuesln(
	jen.ID("TwoFactorSecret").Op(":").ID("user").Dot(
					"TwoFactorSecret",
				))), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
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
		jen.Params(jen.ID("s").Op("*").ID("Service")).ID("UpdatePasswordHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("UpdatePasswordHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot(
					"Value",
				).Call(jen.ID("PasswordChangeMiddlewareCtxKey")).Assert(jen.Op("*").ID("models").Dot(
					"PasswordUpdateInput",
				)),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Debug",
					).Call(jen.Lit("no input found on UpdatePasswordHandler request")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.List(jen.ID("userID"), jen.ID("ok")).Op(":=").ID("ctx").Dot(
					"Value",
				).Call(jen.ID("models").Dot(
					"UserIDKey",
				)).Assert(jen.ID("uint64")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Debug",
					).Call(jen.Lit("no user ID attached to UpdatePasswordHandler request")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.List(jen.ID("user"), jen.ID("sc")).Op(":=").ID("s").Dot(
					"validateCredentialChangeRequest",
				).Call(jen.ID("ctx"), jen.ID("userID"), jen.ID("input").Dot(
					"CurrentPassword",
				),
					jen.ID("input").Dot(
						"TOTPToken",
					)),
				jen.If(jen.ID("sc").Op("!=").Qual("net/http", "StatusOK")).Block(
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.ID("sc")),
					jen.Return(),
				),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("user").Dot(
					"Username",
				)),
				jen.ID("logger").Op(":=").ID("s").Dot(
					"logger",
				).Dot(
					"WithValue",
				).Call(jen.Lit("user"), jen.ID("user").Dot(
					"ID",
				)),

				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("user").Dot(
					"HashedPassword",
				),
					jen.ID("err")).Op("=").ID("s").Dot(
					"authenticator",
				).Dot(
					"HashPassword",
				).Call(jen.ID("ctx"), jen.ID("input").Dot(
					"NewPassword",
				)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error hashing password")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("Database").Dot(
					"UpdateUser",
				).Call(jen.ID("ctx"), jen.ID("user")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered updating user")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusAccepted")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler is a handler for archiving a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("ArchiveHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("userID").Op(":=").ID("s").Dot(
					"userIDFetcher",
				).Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot(
					"logger",
				).Dot(
					"WithValue",
				).Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot("Database").Dot(
					"ArchiveUser",
				).Call(jen.ID("ctx"), jen.ID("userID")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("deleting user from database")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("s").Dot(
					"userCounter",
				).Dot(
					"Decrement",
				).Call(jen.ID("ctx")),
				jen.ID("s").Dot(
					"reporter",
				).Dot(
					"Report",
				).Call(jen.ID("newsman").Dot(
					"Event",
				).Valuesln(
	jen.ID("EventType").Op(":").ID("string").Call(jen.ID("models").Dot(
					"Archive",
				)), jen.ID("Data").Op(":").ID("models").Dot(
					"User",
				).Valuesln(
	jen.ID("ID").Op(":").ID("userID")), jen.ID("Topics").Op(":").Index().ID("string").Valuesln(
	jen.ID("topicName")))),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusNoContent")),
			),
		),
		jen.Line(),
	)
	return ret
}
