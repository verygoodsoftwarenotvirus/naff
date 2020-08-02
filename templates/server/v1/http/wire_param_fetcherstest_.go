package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	for _, typ := range proj.DataTypes {
		if typ.OwnedByAUserAtSomeLevel(proj) {
			code.Add(buildTestProvideSomethingServiceUserIDFetcher(typ)...)
		}

		for _, ot := range proj.FindOwnerTypeChain(typ) {
			code.Add(buildTestProvideSomethingServiceOwnerTypeIDFetcher(typ, ot)...)
		}

		code.Add(buildTestProvideSomethingServiceSomethingIDFetcher(typ)...)
	}

	code.Add(buildTestProvideUsersServiceUserIDFetcher()...)
	code.Add(buildTestProvideWebhooksServiceUserIDFetcher()...)
	code.Add(buildTestProvideWebhooksServiceWebhookIDFetcher()...)
	code.Add(buildTestProvideOAuth2ClientsServiceClientIDFetcher()...)
	code.Add(buildTest_userIDFetcherFromRequestContext(proj)...)
	code.Add(buildTest_buildRouteParamUserIDFetcher(proj)...)

	for _, typ := range proj.DataTypes {
		code.Add(buildTest_buildRouteParamSomethingIDFetcher(proj, typ)...)
	}

	code.Add(buildTest_buildRouteParamWebhookIDFetcher(proj)...)
	code.Add(buildTest_buildRouteParamOAuth2ClientIDFetcher(proj)...)

	return code
}

func buildTestProvideSomethingServiceUserIDFetcher(typ models.DataType) []jen.Code {
	n := typ.Name

	lines := []jen.Code{
		jen.Func().IDf("TestProvide%sServiceUserIDFetcher", n.Plural()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().IDf("Provide%sServiceUserIDFetcher", n.Plural()).Call(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideSomethingServiceSomethingIDFetcher(typ models.DataType) []jen.Code {
	n := typ.Name

	lines := []jen.Code{
		jen.Func().IDf("TestProvide%sService%sIDFetcher", n.Plural(), n.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().IDf("Provide%sService%sIDFetcher", n.Plural(), n.Singular()).Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideSomethingServiceOwnerTypeIDFetcher(typ models.DataType, ot models.DataType) []jen.Code {
	n := typ.Name

	lines := []jen.Code{
		jen.Func().IDf("TestProvide%sService%sIDFetcher", n.Plural(), ot.Name.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().IDf("Provide%sService%sIDFetcher", n.Plural(), ot.Name.Singular()).Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideUsersServiceUserIDFetcher() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideUsersServiceUserIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideUsersServiceUserIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideWebhooksServiceUserIDFetcher() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideWebhooksServiceUserIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideWebhooksServiceUserIDFetcher").Call(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideWebhooksServiceWebhookIDFetcher() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideWebhooksServiceWebhookIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideWebhooksServiceWebhookIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideOAuth2ClientsServiceClientIDFetcher() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideOAuth2ClientsServiceClientIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().ID("ProvideOAuth2ClientsServiceClientIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTest_userIDFetcherFromRequestContext(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_userIDFetcherFromRequestContext").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				utils.BuildFakeVar(proj, "User"),
				jen.ID("expected").Assign().ID("exampleUser").Dot("ToSessionInfo").Call(),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "SessionInfoKey"),
						jen.ID("expected"),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("userIDFetcherFromRequestContext").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected").Dot(constants.UserIDFieldName), jen.ID("actual"), nil),
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
	}

	return lines
}

func buildTest_buildRouteParamUserIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_buildRouteParamUserIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("fn").Assign().ID("buildRouteParamUserIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
				jen.ID("fn").Assign().ID("buildRouteParamUserIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
	}

	return lines
}

func buildTest_buildRouteParamSomethingIDFetcher(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	lines := []jen.Code{
		jen.Func().IDf("Test_buildRouteParam%sIDFetcher", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("fn").Assign().IDf("buildRouteParam%sIDFetcher", sn).Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
				jen.ID("fn").Assign().IDf("buildRouteParam%sIDFetcher", sn).Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
	}

	return lines
}

func buildTest_buildRouteParamWebhookIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_buildRouteParamWebhookIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("fn").Assign().ID("buildRouteParamWebhookIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
				jen.ID("fn").Assign().ID("buildRouteParamWebhookIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
	}

	return lines
}

func buildTest_buildRouteParamOAuth2ClientIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_buildRouteParamOAuth2ClientIDFetcher").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("fn").Assign().ID("buildRouteParamOAuth2ClientIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
				jen.ID("fn").Assign().ID("buildRouteParamOAuth2ClientIDFetcher").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
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
	}

	return lines
}
