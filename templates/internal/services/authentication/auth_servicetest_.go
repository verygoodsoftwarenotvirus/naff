package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.ImportName("github.com/alexedwards/scs/v2/memstore", "memstore")

	code.Add(buildBuildTestService(proj)...)
	code.Add(buildTestProvideAuthService(proj)...)

	return code
}

func buildBuildTestService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID("Service")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
			jen.ID("cfg").Assign().Qual(proj.InternalConfigPackage(), "AuthSettings").Valuesln(
				jen.ID("CookieSecret").MapAssign().Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
			),
			jen.ID("auth").Assign().AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
			jen.ID("userDB").Assign().AddressOf().Qual(proj.TypesPackage("mock"), "UserDataManager").Values(),
			jen.ID("oauth").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
			jen.ID("ed").Assign().Qual(proj.InternalEncodingPackage(), "ProvideResponseEncoder").Call(),
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
	}

	return lines
}

func buildTestProvideAuthService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideAuthService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("cfg").Assign().Qual(proj.InternalConfigPackage(), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").MapAssign().Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
				),
				jen.ID("auth").Assign().AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
				jen.ID("userDB").Assign().AddressOf().Qual(proj.TypesPackage("mock"), "UserDataManager").Values(),
				jen.ID("oauth").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ed").Assign().Qual(proj.InternalEncodingPackage(), "ProvideResponseEncoder").Call(),
				jen.ID("sm").Assign().Qual(constants.SessionManagerLibrary, "New").Call(),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideAuthService").Callln(
					jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
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
		),
	}

	return lines
}
