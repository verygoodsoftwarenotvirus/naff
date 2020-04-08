package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.PointerTo().ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("logger").Assign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("cfg").Assign().AddressOf().Qual(proj.InternalConfigV1Package(), "ServerConfig").Valuesln(
				jen.ID("Auth").MapAssign().Qual(proj.InternalConfigV1Package(), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").MapAssign().Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
				),
			),
			jen.ID("auth").Assign().AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
			jen.ID("userDB").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
			jen.ID("oauth").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
			jen.ID("userIDFetcher").Assign().Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.Return().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call().Dot("ID"),
			),
			jen.ID("ed").Assign().Qual(proj.InternalEncodingV1Package(), "ProvideResponseEncoder").Call(),
			jen.Line(),
			jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideAuthService").Callln(
				jen.ID("logger"),
				jen.ID("cfg"),
				jen.ID("auth"),
				jen.ID("userDB"),
				jen.ID("oauth"),
				jen.ID("userIDFetcher"),
				jen.ID("ed"),
			),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideAuthService").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("cfg").Assign().AddressOf().Qual(proj.InternalConfigV1Package(), "ServerConfig").Valuesln(
					jen.ID("Auth").MapAssign().Qual(proj.InternalConfigV1Package(), "AuthSettings").Valuesln(
						jen.ID("CookieSecret").MapAssign().Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
					),
				),
				jen.ID("auth").Assign().AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
				jen.ID("userDB").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
				jen.ID("oauth").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("userIDFetcher").Assign().Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return(jen.Qual(proj.FakeModelsPackage(), "BuildFakeUser")).Call().Dot("ID"),
				),
				jen.ID("ed").Assign().Qual(proj.InternalEncodingV1Package(), "ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideAuthService").Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("cfg"),
					jen.ID("auth"),
					jen.ID("userDB"),
					jen.ID("oauth"),
					jen.ID("userIDFetcher"),
					jen.ID("ed"),
				),
				utils.AssertNotNil(jen.ID("service"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil config",
				jen.ID("auth").Assign().AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
				jen.ID("userDB").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
				jen.ID("oauth").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("userIDFetcher").Assign().Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Return(jen.Qual(proj.FakeModelsPackage(), "BuildFakeUser")).Call().Dot("ID"),
				),
				jen.ID("ed").Assign().Qual(proj.InternalEncodingV1Package(), "ProvideResponseEncoder").Call(),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideAuthService").Callln(
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.ID("auth"),
					jen.ID("userDB"),
					jen.ID("oauth"),
					jen.ID("userIDFetcher"),
					jen.ID("ed"),
				),
				utils.AssertNil(jen.ID("service"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
	)

	return ret
}
