package http

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireParamFetchersDotGo() *jen.File {
	ret := jen.NewFile("httpserver")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("paramFetcherProviders").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideUserIDFetcher"), jen.ID("ProvideUsernameFetcher"), jen.ID("ProvideOAuth2ServiceClientIDFetcher"), jen.ID("ProvideAuthUserIDFetcher"), jen.ID("ProvideItemIDFetcher"), jen.ID("ProvideWebhooksUserIDFetcher"), jen.ID("ProvideWebhookIDFetcher")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideUserIDFetcher provides a UserIDFetcher").ID("ProvideUserIDFetcher").Params().Params(jen.ID("items").Dot(
		"UserIDFetcher",
	)).Block(
		jen.Return().ID("UserIDFetcher"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideItemIDFetcher provides an ItemIDFetcher").ID("ProvideItemIDFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("items").Dot(
		"ItemIDFetcher",
	)).Block(
		jen.Return().ID("buildChiItemIDFetcher").Call(jen.ID("logger")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideUsernameFetcher provides a UsernameFetcher").ID("ProvideUsernameFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("users").Dot(
		"UserIDFetcher",
	)).Block(
		jen.Return().ID("buildChiUserIDFetcher").Call(jen.ID("logger")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideAuthUserIDFetcher provides a UsernameFetcher").ID("ProvideAuthUserIDFetcher").Params().Params(jen.ID("auth").Dot(
		"UserIDFetcher",
	)).Block(
		jen.Return().ID("UserIDFetcher"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideWebhooksUserIDFetcher provides a UserIDFetcher").ID("ProvideWebhooksUserIDFetcher").Params().Params(jen.ID("webhooks").Dot(
		"UserIDFetcher",
	)).Block(
		jen.Return().ID("UserIDFetcher"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideWebhookIDFetcher provides an WebhookIDFetcher").ID("ProvideWebhookIDFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("webhooks").Dot(
		"WebhookIDFetcher",
	)).Block(
		jen.Return().ID("buildChiWebhookIDFetcher").Call(jen.ID("logger")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideOAuth2ServiceClientIDFetcher provides a ClientIDFetcher").ID("ProvideOAuth2ServiceClientIDFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("oauth2clients").Dot(
		"ClientIDFetcher",
	)).Block(
		jen.Return().ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID("logger")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UserIDFetcher fetches a user ID from a request routed by chi.").ID("UserIDFetcher").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
		jen.Return().ID("req").Dot(
			"Context",
		).Call().Dot(
			"Value",
		).Call(jen.ID("models").Dot(
			"UserIDKey",
		)).Assert(jen.ID("uint64")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildChiUserIDFetcher builds a function that fetches a Username from a request routed by chi.").ID("buildChiUserIDFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("users").Dot(
		"UserIDFetcher",
	)).Block(
		jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("chi").Dot(
				"URLParam",
			).Call(jen.ID("req"), jen.ID("users").Dot(
				"URIParamKey",
			)), jen.Lit(10), jen.Lit(64)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("fetching user ID from request")),
			),
			jen.Return().ID("u"),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// chiItemIDFetcher fetches a ItemID from a request routed by chi.").ID("buildChiItemIDFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
		jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("chi").Dot(
				"URLParam",
			).Call(jen.ID("req"), jen.ID("items").Dot(
				"URIParamKey",
			)), jen.Lit(10), jen.Lit(64)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("fetching ItemID from request")),
			),
			jen.Return().ID("u"),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// chiWebhookIDFetcher fetches a WebhookID from a request routed by chi.").ID("buildChiWebhookIDFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
		jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("chi").Dot(
				"URLParam",
			).Call(jen.ID("req"), jen.ID("webhooks").Dot(
				"URIParamKey",
			)), jen.Lit(10), jen.Lit(64)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("fetching WebhookID from request")),
			),
			jen.Return().ID("u"),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// chiOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi.").ID("buildChiOAuth2ClientIDFetcher").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
		jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("chi").Dot(
				"URLParam",
			).Call(jen.ID("req"), jen.ID("oauth2clients").Dot(
				"URIParamKey",
			)), jen.Lit(10), jen.Lit(64)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("fetching OAuth2ClientID from request")),
			),
			jen.Return().ID("u"),
		),
	),

		jen.Line(),
	)
	return ret
}
