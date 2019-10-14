package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func httpRoutesDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("URIParamKey").Op("=").Lit("oauth2ClientID").Var().ID("oauth2ClientIDURIParamKey").Op("=").Lit("client_id").Var().ID("clientIDKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("client_id"),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachUserIDToSpan provides a consistent way of attaching an user ID to a given span").ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachOAuth2ClientDatabaseIDToSpan provides a consistent way of attaching an oauth2 client ID to a given span").ID("attachOAuth2ClientDatabaseIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("clientID").ID("uint64")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("oauth2client_db_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("clientID"), jen.Lit(10)))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachOAuth2ClientIDToSpan provides a consistent way of attaching a client ID to a given span").ID("attachOAuth2ClientIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("clientID").ID("string")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("client_id"), jen.ID("clientID"))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// randString produces a random string").Comment("// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/").ID("randString").Params().Params(jen.ID("string")).Block(
		jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(32)),
		jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
		jen.Return().Qual("encoding/base32", "StdEncoding").Dot(
			"WithPadding",
		).Call(jen.Qual("encoding/base32", "NoPadding")).Dot(
			"EncodeToString",
		).Call(jen.ID("b")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// fetchUserID grabs a userID out of the request context").Params(jen.ID("s").Op("*").ID("Service")).ID("fetchUserID").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
		jen.If(jen.List(jen.ID("id"), jen.ID("ok")).Op(":=").ID("req").Dot(
			"Context",
		).Call().Dot(
			"Value",
		).Call(jen.ID("models").Dot(
			"UserIDKey",
		)).Assert(jen.ID("uint64")), jen.ID("ok")).Block(
			jen.Return().ID("id"),
		),
		jen.Return().Lit(0),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ListHandler is a handler that returns a list of OAuth2 clients").Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("ListHandler")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.ID("qf").Op(":=").ID("models").Dot(
				"ExtractQueryFilter",
			).Call(jen.ID("req")),
			jen.ID("userID").Op(":=").ID("s").Dot(
				"fetchUserID",
			).Call(jen.ID("req")),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithValue",
			).Call(jen.Lit("user_id"), jen.ID("userID")),
			jen.List(jen.ID("oauth2Clients"), jen.ID("err")).Op(":=").ID("s").Dot(
				"database",
			).Dot(
				"GetOAuth2Clients",
			).Call(jen.ID("ctx"), jen.ID("qf"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.ID("oauth2Clients").Op("=").Op("&").ID("models").Dot(
					"OAuth2ClientList",
				).Valuesln(jen.ID("Clients").Op(":").Index().ID("models").Dot(
					"OAuth2Client",
				).Valuesln()),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("encountered error getting list of oauth2 clients from database")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot(
				"encoderDecoder",
			).Dot(
				"EncodeResponse",
			).Call(jen.ID("res"), jen.ID("oauth2Clients")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("encoding response")),
			),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateHandler is our OAuth2 client creation route").Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("CreateHandler")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot(
				"Value",
			).Call(jen.ID("CreationMiddlewareCtxKey")).Assert(jen.Op("*").ID("models").Dot(
				"OAuth2ClientCreationInput",
			)),
			jen.If(jen.Op("!").ID("ok")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"Info",
				).Call(jen.Lit("valid input not attached to request")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusBadRequest")),
				jen.Return(),
			),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithValues",
			).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("username").Op(":").ID("input").Dot(
				"Username",
			), jen.Lit("scopes").Op(":").ID("input").Dot(
				"Scopes",
			), jen.Lit("redirect_uri").Op(":").ID("input").Dot(
				"RedirectURI",
			))),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot(
				"database",
			).Dot(
				"GetUserByUsername",
			).Call(jen.ID("ctx"), jen.ID("input").Dot(
				"Username",
			)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("fetching user by username")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("input").Dot(
				"BelongsTo",
			).Op("=").ID("user").Dot(
				"ID",
			),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("user").Dot(
				"ID",
			)),
			jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("s").Dot(
				"authenticator",
			).Dot(
				"ValidateLogin",
			).Call(jen.ID("ctx"), jen.ID("user").Dot(
				"HashedPassword",
			), jen.ID("input").Dot(
				"Password",
			), jen.ID("user").Dot(
				"TwoFactorSecret",
			), jen.ID("input").Dot(
				"TOTPToken",
			), jen.ID("user").Dot(
				"Salt",
			)),
			jen.If(jen.Op("!").ID("valid")).Block(
				jen.ID("logger").Dot(
					"Debug",
				).Call(jen.Lit("invalid credentials provided")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("validating user credentials")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("input").Dot(
				"ClientID",
			).Op("=").ID("randString").Call(),
			jen.ID("input").Dot(
				"ClientSecret",
			).Op("=").ID("randString").Call(),
			jen.ID("input").Dot(
				"BelongsTo",
			).Op("=").ID("s").Dot(
				"fetchUserID",
			).Call(jen.ID("req")),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("s").Dot(
				"database",
			).Dot(
				"CreateOAuth2Client",
			).Call(jen.ID("ctx"), jen.ID("input")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("creating oauth2Client in the database")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot(
				"ID",
			)),
			jen.ID("s").Dot(
				"oauth2ClientCounter",
			).Dot(
				"Increment",
			).Call(jen.ID("ctx")),
			jen.ID("res").Dot(
				"WriteHeader",
			).Call(jen.Qual("net/http", "StatusCreated")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot(
				"encoderDecoder",
			).Dot(
				"EncodeResponse",
			).Call(jen.ID("res"), jen.ID("client")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("encoding response")),
			),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ReadHandler is a route handler for retrieving an OAuth2 client").Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("ReadHandler")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.ID("userID").Op(":=").ID("s").Dot(
				"fetchUserID",
			).Call(jen.ID("req")),
			jen.ID("oauth2ClientID").Op(":=").ID("s").Dot(
				"urlClientIDExtractor",
			).Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithValues",
			).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("oauth2_client_id").Op(":").ID("oauth2ClientID"), jen.Lit("user_id").Op(":").ID("userID"))),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("oauth2ClientID")),
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot(
				"database",
			).Dot(
				"GetOAuth2Client",
			).Call(jen.ID("ctx"), jen.ID("oauth2ClientID"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.ID("logger").Dot(
					"Debug",
				).Call(jen.Lit("ReadHandler called on nonexistent client")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusNotFound")),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("error fetching oauth2Client from database")),
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
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("encoding response")),
			),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveHandler is a route handler for archiving an OAuth2 client").Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("ArchiveHandler")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.ID("userID").Op(":=").ID("s").Dot(
				"fetchUserID",
			).Call(jen.ID("req")),
			jen.ID("oauth2ClientID").Op(":=").ID("s").Dot(
				"urlClientIDExtractor",
			).Call(jen.ID("req")),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithValues",
			).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("oauth2_client_id").Op(":").ID("oauth2ClientID"), jen.Lit("user_id").Op(":").ID("userID"))),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("oauth2ClientID")),
			jen.ID("err").Op(":=").ID("s").Dot(
				"database",
			).Dot(
				"ArchiveOAuth2Client",
			).Call(jen.ID("ctx"), jen.ID("oauth2ClientID"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusNotFound")),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("encountered error deleting oauth2 client")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.Return(),
			),
			jen.ID("s").Dot(
				"oauth2ClientCounter",
			).Dot(
				"Decrement",
			).Call(jen.ID("ctx")),
			jen.ID("res").Dot(
				"WriteHeader",
			).Call(jen.Qual("net/http", "StatusNoContent")),
		),
	),

		jen.Line(),
	)
	return ret
}
