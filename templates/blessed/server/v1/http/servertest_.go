package httpserver

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(pkg, ret)

	buildServerLines := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("DebugMode").MapAssign().ID("true"),
			jen.ID("db").MapAssign().Qual(pkg.DatabaseV1Package(), "BuildMockDatabase").Call(),
			jen.ID("config").MapAssign().VarPointer().Qual(pkg.InternalConfigV1Package(), "ServerConfig").Values(),
			jen.ID("encoder").MapAssign().VarPointer().Qual(pkg.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
			jen.ID("httpServer").MapAssign().ID("provideHTTPServer").Call(),
			jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("frontendService").MapAssign().Qual(pkg.ServiceV1FrontendPackage(), "ProvideFrontendService").Callln(
				jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.Qual(pkg.InternalConfigV1Package(), "FrontendSettings").Values(),
			),
			jen.ID("webhooksService").MapAssign().VarPointer().Qual(pkg.ModelsV1Package("mock"), "WebhookDataServer").Values(),
			jen.ID("usersService").MapAssign().VarPointer().Qual(pkg.ModelsV1Package("mock"), "UserDataServer").Values(),
			jen.ID("authService").MapAssign().VarPointer().Qual(pkg.ServiceV1AuthPackage(), "Service").Values(),
		}
		for _, typ := range pkg.DataTypes {
			tpuvn := typ.Name.PluralUnexportedVarName()
			tsn := typ.Name.Singular()
			lines = append(lines,
				jen.IDf("%sService", tpuvn).MapAssign().VarPointer().Qual(pkg.ModelsV1Package("mock"), fmt.Sprintf("%sDataServer", tsn)).Values(),
			)
		}

		lines = append(lines,
			jen.ID("oauth2ClientsService").MapAssign().VarPointer().Qual(pkg.ModelsV1Package("mock"), "OAuth2ClientDataServer").Values(),
		)

		return lines
	}

	ret.Add(
		utils.FakeSeedFunc(),
	)

	ret.Add(
		jen.Func().ID("buildTestServer").Params().Params(jen.Op("*").ID("Server")).Block(
			jen.ID("s").Assign().VarPointer().ID("Server").Valuesln(
				buildServerLines()...,
			),
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	buildProvideServerArgs := func() []jen.Code {
		args := []jen.Code{
			utils.CtxVar(),
			jen.VarPointer().Qual(pkg.InternalConfigV1Package(), "ServerConfig").Valuesln(
				jen.ID("Auth").MapAssign().Qual(pkg.InternalConfigV1Package(), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").MapAssign().Lit("THISISAVERYLONGSTRINGFORTESTPURPOSES"),
				),
			),
			jen.VarPointer().Qual(pkg.ServiceV1AuthPackage(), "Service").Values(),
			jen.VarPointer().Qual(pkg.ServiceV1FrontendPackage(), "Service").Values(),
		}

		for _, typ := range pkg.DataTypes {
			pn := typ.Name.PackageName()
			args = append(args, jen.VarPointer().Qual(pkg.ServiceV1Package(pn), "Service").Values())
		}

		args = append(args,
			jen.VarPointer().Qual(pkg.ServiceV1UsersPackage(), "Service").Values(),
			jen.VarPointer().Qual(pkg.ServiceV1OAuth2ClientsPackage(), "Service").Values(),
			jen.VarPointer().Qual(pkg.ServiceV1WebhooksPackage(), "Service").Values(),
			jen.ID("mockDB"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.VarPointer().Qual(pkg.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
		)

		// if pkg.EnableNewsman {
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
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("mockDB").Assign().Qual(pkg.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(jen.Lit("GetAllWebhooks"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.VarPointer().Qual(pkg.ModelsV1Package(), "WebhookList").Values(), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("ProvideServer").Callln(
					buildProvideServerArgs()...,
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
			)),
		),
		jen.Line(),
	)
	return ret
}
