package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Comment("fetchRandomOAuth2Client retrieves a random client from the list of available clients"),
		jen.Line(),
		jen.Func().ID("fetchRandomOAuth2Client").Params(jen.ID("c").PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(),
			"OAuth2Client",
		)).Block(
			jen.List(jen.ID("clientsRes"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(utils.InlineCtx(), jen.Nil()),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID("clientsRes").IsEqualTo().ID("nil").Or().ID("len").Call(jen.ID("clientsRes").Dot("Clients")).Op("<=").One()).Block(jen.Return().ID("nil")),
			jen.Line(),
			jen.Var().ID("selectedClient").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.For(jen.ID("selectedClient").IsEqualTo().ID("nil")).Block(
				jen.ID("ri").Assign().Qual("math/rand", "Intn").Call(jen.Len(jen.ID("clientsRes").Dot("Clients"))),
				jen.ID("c").Assign().AddressOf().ID("clientsRes").Dot("Clients").Index(jen.ID("ri")),
				jen.If(jen.ID("c").Dot("ClientID").DoesNotEqual().Lit("FIXME")).Block(jen.ID("selectedClient").Equals().ID("c")),
			),
			jen.Line(),
			jen.Return().ID("selectedClient"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("mustBuildCode").Params(jen.ID("totpSecret").String()).Params(jen.String()).Block(
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
				jen.ID("totpSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Block(
				jen.Panic(jen.Err()),
			),
			jen.Return(jen.ID("code")),
		),
	)

	ret.Add(
		jen.Func().ID("buildOAuth2ClientActions").Params(jen.ID("c").PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Params(jen.Map(jen.String()).PointerTo().ID("Action")).Block(
			jen.Return().Map(jen.String()).PointerTo().ID("Action").Valuesln(
				jen.Lit("CreateOAuth2Client").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("CreateOAuth2Client"), jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						utils.CreateCtx(),
						jen.ID("ui").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserCreationInput").Call(),
						jen.List(jen.ID("u"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("ui")),
						jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildHealthCheckRequest").Call(utils.CtxVar()),
						),
						jen.Line(),
						jen.ID("uli").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
							jen.ID("Username").MapAssign().ID("ui").Dot("Username"),
							jen.ID("Password").MapAssign().ID("ui").Dot("Password"),
							jen.ID("TOTPToken").MapAssign().ID("mustBuildCode").Call(jen.ID("u").Dot("TwoFactorSecret")),
						),
						jen.Line(),
						jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("c").Dot("Login").Call(utils.CtxVar(), jen.ID("uli")),
						jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildHealthCheckRequest").Call(utils.CtxVar()),
						),
						jen.Line(),
						jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildCreateOAuth2ClientRequest").Callln(
							utils.CtxVar(),
							jen.ID("cookie"),
							jen.AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
								jen.ID("UserLoginInput").MapAssign().PointerTo().ID("uli"),
							),
						),
						jen.Return().List(jen.ID("req"), jen.Err()),
					),
					jen.ID("Weight").MapAssign().Lit(100)), jen.Lit("GetOAuth2Client").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetOAuth2Client"), jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						jen.If(jen.ID("randomOAuth2Client").Assign().ID("fetchRandomOAuth2Client").Call(jen.ID("c")), jen.ID("randomOAuth2Client").DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildGetOAuth2ClientRequest").Call(utils.InlineCtx(), jen.ID("randomOAuth2Client").Dot("ID")),
						),
						jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").MapAssign().Lit(100)), jen.Lit("GetOAuth2Clients").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetOAuth2Clients"), jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						jen.Return().ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(utils.InlineCtx(), jen.Nil()),
					),
					jen.ID("Weight").MapAssign().Lit(100))),
		),
		jen.Line(),
	)
	return ret
}
