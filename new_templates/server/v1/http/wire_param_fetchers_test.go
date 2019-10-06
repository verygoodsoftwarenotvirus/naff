package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireParamFetchersTestDotGo() *jen.File {
	ret := jen.NewFile("httpserver")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("TestProvideUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").ID("ProvideUserIDFetcher").Call(),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestProvideItemIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").ID("ProvideItemIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestProvideUsernameFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").ID("ProvideUsernameFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestProvideAuthUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").ID("ProvideAuthUserIDFetcher").Call(),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestProvideWebhooksUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").ID("ProvideWebhooksUserIDFetcher").Call(),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestProvideWebhookIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").ID("ProvideWebhookIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestProvideOAuth2ServiceClientIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").ID("ProvideOAuth2ServiceClientIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserIDKey",
			), jen.ID("expected"))),
			jen.ID("actual").Op(":=").ID("UserIDFetcher").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_buildChiUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiUserIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("users").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiUserIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("users").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Lit("expected")))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_buildChiItemIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiItemIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("items").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiItemIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("items").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Lit("expected")))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_buildChiWebhookIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiWebhookIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("webhooks").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiWebhookIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("webhooks").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Lit("expected")))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_buildChiOAuth2ClientIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("oauth2clients").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("fn").Op(":=").ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("chi").Dot(
				"RouteCtxKey",
			), jen.Op("&").ID("chi").Dot(
				"Context",
			).Valuesln(jen.ID("URLParams").Op(":").ID("chi").Dot(
				"RouteParams",
			).Valuesln(jen.ID("Keys").Op(":").Index().ID("string").Valuesln(jen.ID("oauth2clients").Dot(
				"URIParamKey",
			)), jen.ID("Values").Op(":").Index().ID("string").Valuesln(jen.Lit("expected")))))),
			jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),
	)
	return ret
}
