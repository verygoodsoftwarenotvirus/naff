package httpserver

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	buildServerLines := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("DebugMode").MapAssign().ID("true"),
			jen.ID("db").MapAssign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
			jen.ID("config").MapAssign().VarPointer().Qual(proj.InternalConfigV1Package(), "ServerConfig").Values(),
			jen.ID("encoder").MapAssign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
			jen.ID("httpServer").MapAssign().ID("provideHTTPServer").Call(),
			jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("frontendService").MapAssign().Qual(proj.ServiceV1FrontendPackage(), "ProvideFrontendService").Callln(
				jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.Qual(proj.InternalConfigV1Package(), "FrontendSettings").Values(),
			),
			jen.ID("webhooksService").MapAssign().VarPointer().Qual(proj.ModelsV1Package("mock"), "WebhookDataServer").Values(),
			jen.ID("usersService").MapAssign().VarPointer().Qual(proj.ModelsV1Package("mock"), "UserDataServer").Values(),
			jen.ID("authService").MapAssign().VarPointer().Qual(proj.ServiceV1AuthPackage(), "Service").Values(),
		}
		for _, typ := range proj.DataTypes {
			tpuvn := typ.Name.PluralUnexportedVarName()
			tsn := typ.Name.Singular()
			lines = append(lines,
				jen.IDf("%sService", tpuvn).MapAssign().VarPointer().Qual(proj.ModelsV1Package("mock"), fmt.Sprintf("%sDataServer", tsn)).Values(),
			)
		}

		lines = append(lines,
			jen.ID("oauth2ClientsService").MapAssign().VarPointer().Qual(proj.ModelsV1Package("mock"), "OAuth2ClientDataServer").Values(),
		)

		return lines
	}

	ret.Add(
		jen.Func().ID("buildTestServer").Params().Params(jen.PointerTo().ID("Server")).Block(
			jen.ID("s").Assign().VarPointer().ID("Server").Valuesln(
				buildServerLines()...,
			),
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	buildProvideServerArgs := func(cookieSecret string) []jen.Code {
		args := []jen.Code{
			utils.CtxVar(),
			jen.VarPointer().Qual(proj.InternalConfigV1Package(), "ServerConfig").Valuesln(
				jen.ID("Auth").MapAssign().Qual(proj.InternalConfigV1Package(), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").MapAssign().Lit(cookieSecret),
				),
			),
			jen.VarPointer().Qual(proj.ServiceV1AuthPackage(), "Service").Values(),
			jen.VarPointer().Qual(proj.ServiceV1FrontendPackage(), "Service").Values(),
		}

		for _, typ := range proj.DataTypes {
			pn := typ.Name.PackageName()
			args = append(args, jen.VarPointer().Qual(proj.ServiceV1Package(pn), "Service").Values())
		}

		args = append(args,
			jen.VarPointer().Qual(proj.ServiceV1UsersPackage(), "Service").Values(),
			jen.VarPointer().Qual(proj.ServiceV1OAuth2ClientsPackage(), "Service").Values(),
			jen.VarPointer().Qual(proj.ServiceV1WebhooksPackage(), "Service").Values(),
			jen.ID("mockDB"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		)

		// if proj.EnableNewsman {
		args = append(args,
			jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
		)
		// }

		return args
	}

	ret.Add(
		jen.Func().ID("TestProvideServer").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleWebhookList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookList").Call(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleWebhookList"),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideServer").Callln(
					buildProvideServerArgs("THISISAVERYLONGSTRINGFORTESTPURPOSES")...,
				),
				jen.Line(),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid cookie secret",
				jen.ID("exampleWebhookList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookList").Call(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleWebhookList"),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideServer").Callln(
					buildProvideServerArgs("THISSTRINGISNTLONGENOUGH:(")...,
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error fetching webhooks",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList")).Call(jen.Nil()),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideServer").Callln(
					buildProvideServerArgs("THISISAVERYLONGSTRINGFORTESTPURPOSES")...,
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
