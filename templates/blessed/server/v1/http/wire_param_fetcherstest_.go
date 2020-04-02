package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	for _, typ := range proj.DataTypes {
		n := typ.Name

		if typ.BelongsToUser {
			ret.Add(
				jen.Func().IDf("TestProvide%sServiceUserIDFetcher", n.Singular()).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("T").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
						jen.ID("_").Equals().IDf("Provide%sServiceUserIDFetcher", n.Singular()).Call(),
					)),
				),
				jen.Line(),
			)
		}
		if typ.BelongsToStruct != nil {
			ret.Add(
				jen.Func().IDf("TestProvide%sService%sIDFetcher", n.Singular(), typ.BelongsToStruct.Singular()).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("T").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
						jen.ID("_").Equals().IDf("Provide%sService%sIDFetcher", n.Singular(), typ.BelongsToStruct.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
					)),
				),
				jen.Line(),
			)
		}

		ret.Add(
			jen.Func().IDf("TestProvide%sIDFetcher", n.Singular()).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("_").Equals().IDf("Provide%sIDFetcher", n.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				)),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Func().ID("TestProvideUsernameFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("_").Equals().ID("ProvideUsernameFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideAuthUserIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("_").Equals().ID("ProvideAuthUserIDFetcher").Call(),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhooksUserIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("_").Equals().ID("ProvideWebhooksUserIDFetcher").Call(),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhookIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("_").Equals().ID("ProvideWebhookIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideOAuth2ServiceClientIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("_").Equals().ID("ProvideOAuth2ServiceClientIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestUserIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "UserIDKey"),
						jen.ID("expected"),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("UserIDFetcher").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildChiUserIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("fn").Assign().ID("buildChiUserIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1UsersPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Assign().ID("buildChiUserIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().ID("uint64").Call(jen.Lit(0)),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1UsersPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	for _, typ := range proj.DataTypes {
		n := typ.Name
		ret.Add(
			jen.Func().IDf("Test_buildChi%sIDFetcher", n.Singular()).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("fn").Assign().IDf("buildChi%sIDFetcher", n.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
					jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
					jen.Line(),
					jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
					jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.ID("req").Dot("Context").Call(),
							jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
							jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
								jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
									jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1Package(n.PackageName()), "URIParamKey")),
									jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
								),
							),
						),
					),
					jen.Line(),
					jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
					jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.Comment("NOTE: This will probably never happen in dev or production"),
					jen.ID("fn").Assign().IDf("buildChi%sIDFetcher", n.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
					jen.ID("expected").Assign().ID("uint64").Call(jen.Lit(0)),
					jen.Line(),
					jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
					jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.ID("req").Dot("Context").Call(), jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
							jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
								jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
									jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1Package(n.PackageName()), "URIParamKey")),
									jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Lit("expected")),
								),
							),
						),
					),
					jen.Line(),
					jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
					jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				)),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Func().ID("Test_buildChiWebhookIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("fn").Assign().ID("buildChiWebhookIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1WebhooksPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Assign().ID("buildChiWebhookIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().ID("uint64").Call(jen.Lit(0)),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1WebhooksPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildChiOAuth2ClientIDFetcher").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("fn").Assign().ID("buildChiOAuth2ClientIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(), jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid value somehow"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Assign().ID("buildChiOAuth2ClientIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().ID("uint64").Call(jen.Lit(0)),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.VarPointer().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().ID("string").Values(jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().ID("string").Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)
	return ret
}
