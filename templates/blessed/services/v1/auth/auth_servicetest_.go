package auth

import (
	"path/filepath"

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
			jen.ID("cfg").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "ServerConfig").Valuesln(
				jen.ID("Auth").MapAssign().Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "AuthSettings").Valuesln(
					jen.ID("CookieSecret").MapAssign().Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
				),
			),
			jen.ID("auth").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth/mock"), "Authenticator").Values(),
			jen.ID("userDB").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1/mock"), "UserDataManager").Values(),
			jen.ID("oauth").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
			jen.ID("userIDFetcher").Assign().Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Add(utils.FakeUint64Func()),
			),
			jen.ID("ed").Assign().Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "ProvideResponseEncoder").Call(),
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
