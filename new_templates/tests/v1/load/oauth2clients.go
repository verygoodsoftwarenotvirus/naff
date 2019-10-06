package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Func())
	ret.Add(jen.Func().ID("buildOAuth2ClientActions").Params(jen.ID("c").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
		jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(jen.Lit("CreateOAuth2Client").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("CreateOAuth2Client"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.ID("ui").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomUserInput").Call(),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("c").Dot(
				"CreateUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("ui")),
			jen.If(
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.Return().ID("c").Dot(
					"BuildHealthCheckRequest",
				).Call(),
			),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("c").Dot(
				"Login",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("u").Dot(
				"Username",
			), jen.ID("ui").Dot(
				"Password",
			), jen.ID("u").Dot(
				"TwoFactorSecret",
			)),
			jen.If(
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.Return().ID("c").Dot(
					"BuildHealthCheckRequest",
				).Call(),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot(
				"BuildCreateOAuth2ClientRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("cookie"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomOAuth2ClientInput").Call(jen.ID("u").Dot(
				"Username",
			), jen.ID("ui").Dot(
				"Password",
			), jen.ID("u").Dot(
				"TwoFactorSecret",
			))),
			jen.Return().List(jen.ID("req"), jen.ID("err")),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetOAuth2Client").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetOAuth2Client"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.If(
				jen.ID("randomOAuth2Client").Op(":=").ID("fetchRandomOAuth2Client").Call(jen.ID("c")),
				jen.ID("randomOAuth2Client").Op("!=").ID("nil"),
			).Block(
				jen.Return().ID("c").Dot(
					"BuildGetOAuth2ClientRequest",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("randomOAuth2Client").Dot(
					"ID",
				)),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetOAuth2Clients").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetOAuth2Clients"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.Return().ID("c").Dot(
				"BuildGetOAuth2ClientsRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
		), jen.ID("Weight").Op(":").Lit(100))),
	),
	)
	return ret
}
