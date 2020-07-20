package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("auth")

	utils.AddImports(proj, code)

	code.ImportName("github.com/alexedwards/scs/v2/memstore", "memstore")

	code.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("cfg").Assign().AddressOf().Qual(proj.InternalConfigV1Package(), "ServerConfig").Valuesln(
				jen.ID("Auth").MapAssign().Qual(proj.InternalConfigV1Package(), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").MapAssign().Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
				),
			),
			jen.ID("auth").Assign().AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
			jen.ID("userDB").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
			jen.ID("oauth").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
			jen.ID("ed").Assign().Qual(proj.InternalEncodingV1Package(), "ProvideResponseEncoder").Call(),
			jen.Line(),
			jen.ID("sm").Assign().Qual(constants.SessionManagerLibrary, "New").Call(),
			jen.Comment("this is currently the default, but in case that changes"),
			jen.ID("sm").Dot("Store").Equals().Qual("github.com/alexedwards/scs/v2/memstore", "New").Call(),
			jen.Line(),
			jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideAuthService").Callln(
				jen.ID(constants.LoggerVarName),
				jen.ID("cfg"),
				jen.ID("auth"),
				jen.ID("userDB"),
				jen.ID("oauth"),
				jen.ID("sm"),
				jen.ID("ed"),
			),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideAuthService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.Comment("T.Parallel() // NOTE: undo when data race is fixed"),
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
				jen.ID("ed").Assign().Qual(proj.InternalEncodingV1Package(), "ProvideResponseEncoder").Call(),
				jen.ID("sm").Assign().Qual(constants.SessionManagerLibrary, "New").Call(),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideAuthService").Callln(
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("cfg"),
					jen.ID("auth"),
					jen.ID("userDB"),
					jen.ID("oauth"),
					jen.ID("sm"),
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
				jen.ID("ed").Assign().Qual(proj.InternalEncodingV1Package(), "ProvideResponseEncoder").Call(),
				jen.ID("sm").Assign().Qual(constants.SessionManagerLibrary, "New").Call(),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideAuthService").Callln(
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Nil(),
					jen.ID("auth"),
					jen.ID("userDB"),
					jen.ID("oauth"),
					jen.ID("sm"),
					jen.ID("ed"),
				),
				utils.AssertNil(jen.ID("service"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
	)

	return code
}
