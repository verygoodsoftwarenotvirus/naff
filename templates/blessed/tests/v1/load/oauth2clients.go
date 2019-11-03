package load

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Comment("fetchRandomOAuth2Client retrieves a random client from the list of available clients"),
		jen.Line(),
		jen.Func().ID("fetchRandomOAuth2Client").Params(jen.ID("c").Op("*").Qual(filepath.Join(pkgRoot, "client/v1/http"), "V1Client")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"),
			"OAuth2Client",
		)).Block(
			jen.List(jen.ID("clientsRes"), jen.ID("err")).Op(":=").ID("c").Dot("GetOAuth2Clients").Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("clientsRes").Op("==").ID("nil").Op("||").ID("len").Call(jen.ID("clientsRes").Dot("Clients")).Op("<=").Lit(1)).Block(jen.Return().ID("nil")),
			jen.Line(),
			jen.Var().ID("selectedClient").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
			jen.For(jen.ID("selectedClient").Op("==").ID("nil")).Block(
				jen.ID("ri").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.ID("clientsRes").Dot("Clients"))),
				jen.ID("c").Op(":=").Op("&").ID("clientsRes").Dot("Clients").Index(jen.ID("ri")),
				jen.If(jen.ID("c").Dot("ClientID").Op("!=").Lit("FIXME")).Block(jen.ID("selectedClient").Op("=").ID("c")),
			),
			jen.Line(),
			jen.Return().ID("selectedClient"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildOAuth2ClientActions").Params(jen.ID("c").Op("*").Qual(filepath.Join(pkgRoot, "client/v1/http"), "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
			jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(
				jen.Lit("CreateOAuth2Client").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("CreateOAuth2Client"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.ID("ui").Op(":=").Qual(filepath.Join(pkgRoot, "tests/v1/testutil/rand/model"), "RandomUserInput").Call(),
						jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("c").Dot("CreateUser").Call(jen.Qual("context", "Background").Call(), jen.ID("ui")),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildHealthCheckRequest").Call(),
						),
						jen.Line(),
						jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("c").Dot("Login").Call(jen.Qual("context", "Background").Call(), jen.ID("u").Dot("Username"),
							jen.ID("ui").Dot("Password"),
							jen.ID("u").Dot("TwoFactorSecret"),
						),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildHealthCheckRequest").Call(),
						),
						jen.Line(),
						jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("BuildCreateOAuth2ClientRequest").Callln(
							jen.Qual("context", "Background").Call(),
							jen.ID("cookie"),
							jen.Qual(filepath.Join(pkgRoot, "tests/v1/testutil/rand/model"), "RandomOAuth2ClientInput").Callln(
								jen.ID("u").Dot("Username"),
								jen.ID("ui").Dot("Password"),
								jen.ID("u").Dot("TwoFactorSecret"),
							),
						),
						jen.Return().List(jen.ID("req"), jen.ID("err")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetOAuth2Client").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("GetOAuth2Client"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomOAuth2Client").Op(":=").ID("fetchRandomOAuth2Client").Call(jen.ID("c")), jen.ID("randomOAuth2Client").Op("!=").ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildGetOAuth2ClientRequest").Call(jen.Qual("context", "Background").Call(), jen.ID("randomOAuth2Client").Dot("ID")),
						),
						jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetOAuth2Clients").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("GetOAuth2Clients"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
					),
					jen.ID("Weight").Op(":").Lit(100))),
		),
		jen.Line(),
	)
	return ret
}
