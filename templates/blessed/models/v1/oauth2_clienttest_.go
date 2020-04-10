package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_GetID").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().Lit("123"),
				jen.ID("oac").Assign().AddressOf().ID("OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().ID("expected"),
				),
				utils.AssertEqual(jen.ID("expected"), jen.ID("oac").Dot("GetID").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_GetSecret").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().Lit("123"),
				jen.ID("oac").Assign().AddressOf().ID("OAuth2Client").Valuesln(
					jen.ID("ClientSecret").MapAssign().ID("expected"),
				),
				utils.AssertEqual(jen.ID("expected"), jen.ID("oac").Dot("GetSecret").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_GetDomain").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().Lit("123"),
				jen.ID("oac").Assign().AddressOf().ID("OAuth2Client").Valuesln(
					jen.ID("RedirectURI").MapAssign().ID("expected"),
				),
				utils.AssertEqual(jen.ID("expected"), jen.ID("oac").Dot("GetDomain").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_GetUserID").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expectation").Assign().Uint64().Call(jen.Lit(123)),
				jen.ID("expected").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expectation")),
				jen.ID("oac").Assign().AddressOf().ID("OAuth2Client").Valuesln(
					jen.ID("BelongsToUser").MapAssign().ID("expectation"),
				),
				utils.AssertEqual(jen.ID("expected"), jen.ID("oac").Dot("GetUserID").Call(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_HasScope").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("oac").Assign().AddressOf().ID("OAuth2Client").Valuesln(
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff")),
				),
				jen.Line(),
				utils.AssertTrue(jen.ID("oac").Dot("HasScope").Call(jen.ID("oac").Dot("Scopes").Index(jen.Zero())), nil),
				utils.AssertFalse(jen.ID("oac").Dot("HasScope").Call(jen.Lit("blah")), nil),
				utils.AssertFalse(jen.ID("oac").Dot("HasScope").Call(jen.EmptyString()), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
