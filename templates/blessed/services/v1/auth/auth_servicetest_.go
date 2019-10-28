package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func authServiceTestDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("logger").Op(":=").ID("noop").Dot("ProvideNoopLogger").Call(),
			jen.ID("cfg").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "ServerConfig").Valuesln(
				jen.ID("Auth").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "AuthSettings").Valuesln(
					jen.ID("CookieSecret").Op(":").Lit("BLAHBLAHBLAHPRETENDTHISISSECRET!"),
				),
			),
			jen.ID("auth").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(),
			jen.ID("userDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager").Values(),
			jen.ID("oauth").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Values(),
			jen.ID("userIDFetcher").Op(":=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(1),
			),
			jen.ID("ed").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding", "ProvideResponseEncoder").Call(),
			jen.Line(),
			jen.ID("service").Op(":=").ID("ProvideAuthService").Callln(
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
