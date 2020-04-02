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
		jen.Func().ID("TestOAuth2Client_GetID").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Lit("uint64(123)"),
				jen.ID("oac").Assign().VarPointer().ID("OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().ID("expected"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot("GetID").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_GetSecret").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Lit("uint64(123)"),
				jen.ID("oac").Assign().VarPointer().ID("OAuth2Client").Valuesln(
					jen.ID("ClientSecret").MapAssign().ID("expected"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot("GetSecret").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_GetDomain").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Lit("uint64(123)"),
				jen.ID("oac").Assign().VarPointer().ID("OAuth2Client").Valuesln(
					jen.ID("RedirectURI").MapAssign().ID("expected"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot("GetDomain").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_GetUserID").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expectation").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("expectation")),
				jen.ID("oac").Assign().VarPointer().ID("OAuth2Client").Valuesln(
					jen.ID("BelongsToUser").MapAssign().ID("expectation"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot("GetUserID").Call()),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Client_HasScope").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("oac").Assign().VarPointer().ID("OAuth2Client").Valuesln(
					jen.ID("Scopes").MapAssign().Index().ID("string").Values(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff")),
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("oac").Dot("HasScope").Call(jen.ID("oac").Dot("Scopes").Index(jen.Lit(0)))),
				jen.Qual("github.com/stretchr/testify/assert", "False").Call(jen.ID("t"), jen.ID("oac").Dot("HasScope").Call(jen.Lit("blah"))),
				jen.Qual("github.com/stretchr/testify/assert", "False").Call(jen.ID("t"), jen.ID("oac").Dot("HasScope").Call(jen.Lit(""))),
			)),
		),
		jen.Line(),
	)
	return ret
}
