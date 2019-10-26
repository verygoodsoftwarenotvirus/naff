package http

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersTestDotGo(types []models.DataType) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("TestProvideUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("_").Op("=").ID("ProvideUserIDFetcher").Call(),
			)),
		),
		jen.Line(),
	)

	for _, typ := range types {
		n := typ.Name
		ret.Add(
			jen.Func().IDf("TestProvide%sIDFetcher", n.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("_").Op("=").IDf("Provide%sIDFetcher", n.Singular()).Call(jen.ID("noop").Dot(
						"ProvideNoopLogger",
					).Call()),
				)),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Func().ID("TestProvideUsernameFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("_").Op("=").ID("ProvideUsernameFetcher").Call(jen.ID("noop").Dot(
					"ProvideNoopLogger",
				).Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideAuthUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("_").Op("=").ID("ProvideAuthUserIDFetcher").Call(),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhooksUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("_").Op("=").ID("ProvideWebhooksUserIDFetcher").Call(),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhookIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("_").Op("=").ID("ProvideWebhookIDFetcher").Call(jen.ID("noop").Dot(
					"ProvideNoopLogger",
				).Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideOAuth2ServiceClientIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("_").Op("=").ID("ProvideOAuth2ServiceClientIDFetcher").Call(jen.ID("noop").Dot(
					"ProvideNoopLogger",
				).Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserIDKey"),
						jen.ID("expected"),
					),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("UserIDFetcher").Call(jen.ID("req")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildChiUserIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("fn").Op(":=").ID("buildChiUserIDFetcher").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID("users").Dot("URIParamKey")),
								jen.ID("Values").Op(":").Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Op(":=").ID("buildChiUserIDFetcher").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID("users").Dot("URIParamKey")),
								jen.ID("Values").Op(":").Index().ID("string").Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	for _, typ := range types {
		n := typ.Name
		ret.Add(
			jen.Func().IDf("Test_buildChi%sIDFetcher", n.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("fn").Op(":=").IDf("buildChi%sIDFetcher", n.Singular()).Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
					jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
					jen.Line(),
					jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
					jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.ID("req").Dot("Context").Call(),
							jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
							jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
								jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
									jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID(n.PackageName()).Dot("URIParamKey")),
									jen.ID("Values").Op(":").Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
								),
							),
						),
					),
					jen.Line(),
					jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.Comment("NOTE: This will probably never happen in dev or production"),
					jen.ID("fn").Op(":=").IDf("buildChi%sIDFetcher", n.Singular()).Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
					jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
					jen.Line(),
					jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
					jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.ID("req").Dot("Context").Call(), jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
							jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
								jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
									jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID(n.PackageName()).Dot("URIParamKey")),
									jen.ID("Values").Op(":").Index().ID("string").Values(jen.Lit("expected")),
								),
							),
						),
					),
					jen.Line(),
					jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				)),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Func().ID("Test_buildChiWebhookIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("fn").Op(":=").ID("buildChiWebhookIDFetcher").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID("webhooks").Dot("URIParamKey")),
								jen.ID("Values").Op(":").Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Op(":=").ID("buildChiWebhookIDFetcher").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID("webhooks").Dot("URIParamKey")),
								jen.ID("Values").Op(":").Index().ID("string").Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildChiOAuth2ClientIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("fn").Op(":=").ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(), jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID("oauth2clients").Dot("URIParamKey")),
								jen.ID("Values").Op(":").Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Op(":=").ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call()),
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
				jen.Line(),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.Op("&").Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").Op(":").Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Values(jen.ID("oauth2clients").Dot("URIParamKey")),
								jen.ID("Values").Op(":").Index().ID("string").Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)
	return ret
}
