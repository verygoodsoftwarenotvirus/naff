package server

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildBuildTestServer(proj)...)
	code.Add(buildTestProvideServer(proj)...)

	return code
}

func buildProvideServerArgs(proj *models.Project, cookieSecret string) []jen.Code {
	provideServerArgs := []jen.Code{
		constants.CtxVar(),
		jen.AddressOf().Qual(proj.InternalConfigPackage(), "ServerConfig").Valuesln(
			jen.ID("Auth").MapAssign().Qual(proj.InternalConfigPackage(), "AuthSettings").Valuesln(
				jen.ID("CookieSecret").MapAssign().Lit(cookieSecret),
			),
		),
		jen.AddressOf().Qual(proj.ServiceAuthPackage(), "Service").Values(),
		jen.AddressOf().Qual(proj.ServiceFrontendPackage(), "Service").Values(),
	}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.PackageName()
		provideServerArgs = append(provideServerArgs, jen.AddressOf().Qual(proj.ServicePackage(pn), "Service").Values())
	}

	provideServerArgs = append(provideServerArgs,
		jen.AddressOf().Qual(proj.ServiceUsersPackage(), "Service").Values(),
		jen.AddressOf().Qual(proj.ServiceOAuth2ClientsPackage(), "Service").Values(),
		jen.AddressOf().Qual(proj.ServiceWebhooksPackage(), "Service").Values(),
		jen.ID("mockDB"), jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
		jen.AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
	)

	// if proj.EnableNewsman {
	provideServerArgs = append(provideServerArgs,
		jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
	)
	// }

	return provideServerArgs
}

func buildBuildTestServer(proj *models.Project) []jen.Code {
	buildServerLines := []jen.Code{
		jen.ID("DebugMode").MapAssign().True(),
		jen.ID("db").MapAssign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("config").MapAssign().AddressOf().Qual(proj.InternalConfigPackage(), "ServerConfig").Values(),
		jen.ID("encoder").MapAssign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
		jen.ID("httpServer").MapAssign().ID("provideHTTPServer").Call(),
		jen.ID(constants.LoggerVarName).MapAssign().Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
		jen.ID("frontendService").MapAssign().Qual(proj.ServiceFrontendPackage(), "ProvideFrontendService").Callln(
			jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
			jen.Qual(proj.InternalConfigPackage(), "FrontendSettings").Values(),
		),
		jen.ID("webhooksService").MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), "WebhookDataServer").Values(),
		jen.ID("usersService").MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), "UserDataServer").Values(),
		jen.ID("authService").MapAssign().AddressOf().Qual(proj.ServiceAuthPackage(), "Service").Values(),
	}
	for _, typ := range proj.DataTypes {
		tpuvn := typ.Name.PluralUnexportedVarName()
		tsn := typ.Name.Singular()
		buildServerLines = append(buildServerLines,
			jen.IDf("%sService", tpuvn).MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataServer", tsn)).Values(),
		)
	}

	buildServerLines = append(buildServerLines,
		jen.ID("oauth2ClientsService").MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), "OAuth2ClientDataServer").Values(),
	)

	lines := []jen.Code{
		jen.Func().ID("buildTestServer").Params().Params(jen.PointerTo().ID("Server")).Body(
			jen.ID("s").Assign().AddressOf().ID("Server").Valuesln(
				buildServerLines...,
			),
			jen.Line(),
			jen.Return().ID("s"),
		),
		jen.Line(),
	}

	return lines
}
func buildTestProvideServer(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideServer").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.Line(),
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("WebhookList")),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideServer").Callln(
					buildProvideServerArgs(proj, "THISISAVERYLONGSTRINGFORTESTPURPOSES")...,
				),
				jen.Line(),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid cookie secret",
				jen.Line(),
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("WebhookList")),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideServer").Callln(
					buildProvideServerArgs(proj, "THISSTRINGISNTLONGENOUGH:(")...,
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error fetching webhooks",
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), "WebhookList")).Call(jen.Nil()),
					constants.ObligatoryError(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideServer").Callln(
					buildProvideServerArgs(proj, "THISISAVERYLONGSTRINGFORTESTPURPOSES")...,
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}
