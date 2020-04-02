package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Comment("fetchRandomOAuth2Client retrieves a random client from the list of available clients"),
		jen.Line(),
		jen.Func().ID("fetchRandomOAuth2Client").Params(jen.ID("c").Op("*").Qual(pkg.HTTPClientV1Package(), "V1Client")).Params(jen.Op("*").Qual(pkg.ModelsV1Package(),
			"OAuth2Client",
		)).Block(
			jen.List(jen.ID("clientsRes"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(utils.InlineCtx(), jen.Nil()),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Op("||").ID("clientsRes").Op("==").ID("nil").Op("||").ID("len").Call(jen.ID("clientsRes").Dot("Clients")).Op("<=").Lit(1)).Block(jen.Return().ID("nil")),
			jen.Line(),
			jen.Var().ID("selectedClient").Op("*").Qual(pkg.ModelsV1Package(), "OAuth2Client"),
			jen.For(jen.ID("selectedClient").Op("==").ID("nil")).Block(
				jen.ID("ri").Assign().Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.ID("clientsRes").Dot("Clients"))),
				jen.ID("c").Assign().VarPointer().ID("clientsRes").Dot("Clients").Index(jen.ID("ri")),
				jen.If(jen.ID("c").Dot("ClientID").DoesNotEqual().Lit("FIXME")).Block(jen.ID("selectedClient").Equals().ID("c")),
			),
			jen.Line(),
			jen.Return().ID("selectedClient"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildOAuth2ClientActions").Params(jen.ID("c").Op("*").Qual(pkg.HTTPClientV1Package(), "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
			jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(
				jen.Lit("CreateOAuth2Client").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("CreateOAuth2Client"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.ID("ui").Assign().Qual(pkg.FakeModelsPackage(), "RandomUserInput").Call(),
						jen.List(jen.ID("u"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(utils.InlineCtx(), jen.ID("ui")),
						jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildHealthCheckRequest").Call(),
						),
						jen.Line(),
						jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("c").Dot("Login").Call(utils.InlineCtx(), jen.ID("u").Dot("Username"),
							jen.ID("ui").Dot("Password"),
							jen.ID("u").Dot("TwoFactorSecret"),
						),
						jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildHealthCheckRequest").Call(),
						),
						jen.Line(),
						jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildCreateOAuth2ClientRequest").Callln(
							utils.InlineCtx(),
							jen.ID("cookie"),
							jen.Qual(pkg.FakeModelsPackage(), "RandomOAuth2ClientInput").Callln(
								jen.ID("u").Dot("Username"),
								jen.ID("ui").Dot("Password"),
								jen.ID("u").Dot("TwoFactorSecret"),
							),
						),
						jen.Return().List(jen.ID("req"), jen.Err()),
					),
					jen.ID("Weight").MapAssign().Lit(100)), jen.Lit("GetOAuth2Client").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetOAuth2Client"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomOAuth2Client").Assign().ID("fetchRandomOAuth2Client").Call(jen.ID("c")), jen.ID("randomOAuth2Client").DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildGetOAuth2ClientRequest").Call(utils.InlineCtx(), jen.ID("randomOAuth2Client").Dot("ID")),
						),
						jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").MapAssign().Lit(100)), jen.Lit("GetOAuth2Clients").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetOAuth2Clients"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(utils.InlineCtx(), jen.Nil()),
					),
					jen.ID("Weight").MapAssign().Lit(100))),
		),
		jen.Line(),
	)
	return ret
}
