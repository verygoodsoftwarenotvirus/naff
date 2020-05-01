package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	for _, typ := range proj.DataTypes {
		n := typ.Name

		if typ.OwnedByAUserAtSomeLevel(proj) {
			ret.Add(
				jen.Func().IDf("TestProvide%sServiceUserIDFetcher", n.Plural()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
					jen.ID("T").Dot("Parallel").Call(),
					jen.Line(),
					utils.BuildSubTestWithoutContext(
						"obligatory",
						jen.Underscore().Equals().IDf("Provide%sServiceUserIDFetcher", n.Plural()).Call(),
					),
				),
				jen.Line(),
			)
		}

		for _, ot := range proj.FindOwnerTypeChain(typ) {
			ret.Add(
				jen.Func().IDf("TestProvide%sService%sIDFetcher", n.Plural(), ot.Name.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
					jen.ID("T").Dot("Parallel").Call(),
					jen.Line(),
					utils.BuildSubTestWithoutContext(
						"obligatory",
						jen.Underscore().Equals().IDf("Provide%sService%sIDFetcher", n.Plural(), ot.Name.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
					),
				),
				jen.Line(),
			)
		}

		ret.Add(
			jen.Func().IDf("TestProvide%sService%sIDFetcher", n.Plural(), n.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"obligatory",
					jen.Underscore().Equals().IDf("Provide%sService%sIDFetcher", n.Plural(), n.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Func().ID("TestProvideUsersServiceUserIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideUsersServiceUserIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideAuthServiceUserIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideAuthServiceUserIDFetcher").Call(),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhooksServiceUserIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideWebhooksServiceUserIDFetcher").Call(),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideWebhooksServiceWebhookIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideWebhooksServiceWebhookIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideOAuth2ClientsServiceClientIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideOAuth2ClientsServiceClientIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_userIDFetcherFromRequestContext").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "UserIDKey"),
						jen.ID("expected"),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("userIDFetcherFromRequestContext").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without attached value",
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("actual").Assign().ID("userIDFetcherFromRequestContext").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertZero(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildRouteParamUserIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("fn").Assign().ID("buildRouteParamUserIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1UsersPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().String().Values(utils.FormatString("%d", jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid value somehow",
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Assign().ID("buildRouteParamUserIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Uint64().Call(jen.Zero()),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1UsersPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().String().Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	for _, typ := range proj.DataTypes {
		n := typ.Name
		ret.Add(
			jen.Func().IDf("Test_buildRouteParam%sIDFetcher", n.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"happy path",
					jen.ID("fn").Assign().IDf("buildRouteParam%sIDFetcher", n.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
					jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
					jen.Line(),
					jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
					jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.ID(constants.RequestVarName).Dot("Context").Call(),
							jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
							jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
								jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
									jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1Package(n.PackageName()), "URIParamKey")),
									jen.ID("Values").MapAssign().Index().String().Values(utils.FormatString("%d", jen.ID("expected"))),
								),
							),
						),
					),
					jen.Line(),
					jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
					utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"with invalid value somehow",
					jen.Comment("NOTE: This will probably never happen in dev or production"),
					jen.ID("fn").Assign().IDf("buildRouteParam%sIDFetcher", n.Singular()).Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
					jen.ID("expected").Assign().Uint64().Call(jen.Zero()),
					jen.Line(),
					jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
					jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
							jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
								jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
									jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1Package(n.PackageName()), "URIParamKey")),
									jen.ID("Values").MapAssign().Index().String().Values(jen.Lit("expected")),
								),
							),
						),
					),
					jen.Line(),
					jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
					utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Func().ID("Test_buildRouteParamWebhookIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("fn").Assign().ID("buildRouteParamWebhookIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1WebhooksPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().String().Values(utils.FormatString("%d", jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid value somehow",
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Assign().ID("buildRouteParamWebhookIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Uint64().Call(jen.Zero()),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1WebhooksPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().String().Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildRouteParamOAuth2ClientIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("fn").Assign().ID("buildRouteParamOAuth2ClientIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().String().Values(utils.FormatString("%d", jen.ID("expected"))),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid value somehow",
				jen.Comment("NOTE: This will probably never happen in dev or production"),
				jen.ID("fn").Assign().ID("buildRouteParamOAuth2ClientIDFetcher").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.ID("expected").Assign().Uint64().Call(jen.Zero()),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual("github.com/go-chi/chi", "RouteCtxKey"),
						jen.AddressOf().Qual("github.com/go-chi/chi", "Context").Valuesln(
							jen.ID("URLParams").MapAssign().Qual("github.com/go-chi/chi", "RouteParams").Valuesln(
								jen.ID("Keys").MapAssign().Index().String().Values(jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "URIParamKey")),
								jen.ID("Values").MapAssign().Index().String().Values(jen.Lit("expected")),
							),
						),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("fn").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	return ret
}
