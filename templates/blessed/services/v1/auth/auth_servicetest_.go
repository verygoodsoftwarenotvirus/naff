package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(pkg, ret)

	ret.Add(utils.FakeSeedFunc())

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("logger").Assign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("cfg").Assign().VarPointer().Qual(pkg.InternalConfigV1Package(), "ServerConfig").Valuesln(
				jen.ID("Auth").MapAssign().Qual(pkg.InternalConfigV1Package(), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").MapAssign().Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
				),
			),
			jen.ID("auth").Assign().VarPointer().Qual(pkg.InternalAuthV1Package("mock"), "Authenticator").Values(),
			jen.ID("userDB").Assign().VarPointer().Qual(pkg.ModelsV1Package("mock"), "UserDataManager").Values(),
			jen.ID("oauth").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
			jen.ID("userIDFetcher").Assign().Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Add(utils.FakeUint64Func()),
			),
			jen.ID("ed").Assign().Qual(pkg.InternalEncodingV1Package(), "ProvideResponseEncoder").Call(),
			jen.Line(),
			jen.ID("service").Assign().ID("ProvideAuthService").Callln(
				jen.ID("logger"),
				jen.ID("cfg"),
				jen.ID("auth"),
				jen.ID("userDB"),
				jen.ID("oauth"),
				jen.ID("userIDFetcher"),
				jen.ID("ed"),
			),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	)
	return ret
}
