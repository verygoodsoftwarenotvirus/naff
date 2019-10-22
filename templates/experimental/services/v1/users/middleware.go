package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("UserCreationMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("user_creation_input"),
			jen.ID("PasswordChangeMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("user_password_change"),
			jen.ID("TOTPSecretRefreshMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("totp_refresh"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserInputMiddleware fetches user input from requests"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot(
					"UserInput",
				)),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("UserInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"DecodeRequest",
				).Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered decoding request body")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("UserCreationMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req").Dot(
					"WithContext",
				).Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("PasswordUpdateInputMiddleware fetches password update input from requests"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot(
					"PasswordUpdateInput",
				)),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("PasswordUpdateInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"DecodeRequest",
				).Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered decoding request body")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("PasswordChangeMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req").Dot(
					"WithContext",
				).Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("TOTPSecretRefreshInputMiddleware fetches 2FA update input from requests"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot(
					"TOTPSecretRefreshInput",
				)),
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("TOTPSecretRefreshInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"DecodeRequest",
				).Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered decoding request body")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("TOTPSecretRefreshMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req").Dot(
					"WithContext",
				).Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)
	return ret
}
