package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2TestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("mustBuildCode").Params(jen.ID("t").ParamPointer().Qual("testing", "T"), jen.ID("totpSecret").ID("string")).Params(jen.ID("string")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("totpSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Return().ID("code"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyOAuth2ClientInput").Params(jen.ID("t").ParamPointer().Qual("testing", "T"), jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpToken")).ID("string")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("x").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
				jen.ID("UserLoginInput").MapAssign().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("username"),
					jen.ID("Password").MapAssign().ID("password"),
					jen.ID("TOTPToken").MapAssign().ID("mustBuildCode").Call(jen.ID("t"), jen.ID("totpToken")),
				),
				jen.ID("Scopes").MapAssign().Index().ID("string").Values(jen.Lit("*")),
				jen.ID("RedirectURI").MapAssign().Lit("http://localhost"),
			),
			jen.Line(),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("convertInputToClient").Params(jen.ID("input").Op("*").Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "OAuth2Client")).Block(
			jen.Return().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
				jen.ID("ClientID").MapAssign().ID("input").Dot("ClientID"),
				jen.ID("ClientSecret").MapAssign().ID("input").Dot("ClientSecret"),
				jen.ID("RedirectURI").MapAssign().ID("input").Dot("RedirectURI"),
				jen.ID("Scopes").MapAssign().ID("input").Dot("Scopes"),
				jen.ID("BelongsToUser").MapAssign().ID("input").Dot("BelongsToUser")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("checkOAuth2ClientEquality").Params(jen.ID("t").ParamPointer().Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").Qual(proj.ModelsV1Package(), "OAuth2Client")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ID")),
			jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("actual").Dot("ClientID")),
			jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("actual").Dot("ClientSecret")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("RedirectURI"), jen.ID("actual").Dot("RedirectURI")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Scopes"), jen.ID("actual").Dot("Scopes")),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("CreatedOn")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual").Dot("ArchivedOn")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Clients").Params(jen.ID("test").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.ID("_ctx").Assign().Add(utils.InlineCtx()),
			jen.Line(),
			jen.Comment("create user"),
			jen.List(jen.ID("x"), jen.ID("y"), jen.ID("cookie")).Assign().ID("buildDummyUser").Call(jen.ID("test")),
			jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
			jen.Line(),
			jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(
				jen.ID("test"),
				jen.ID("x").Dot("Username"),
				jen.ID("y").Dot("Password"),
				jen.ID("x").Dot("TwoFactorSecret"),
			),
			jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateOAuth2Client").Call(jen.ID("_ctx"), jen.ID("cookie"), jen.ID("input")),
			jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.Err()),
			jen.Line(),
			jen.List(jen.ID("testClient"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
				jen.ID("_ctx"),
				jen.ID("premade").Dot("ClientID"),
				jen.ID("premade").Dot("ClientSecret"),
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
				jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("PlainClient").Call(),
				jen.ID("premade").Dot("Scopes"),
				jen.ID("debug"),
			),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("test"), jen.Err(), jen.Lit("error setting up auxiliary client")),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be creatable"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					utils.CreateCtx(),
					jen.Line(),
					jen.Comment("Create oauth2Client"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert oauth2Client equality"),
					jen.ID("checkOAuth2ClientEquality").Call(jen.ID("t"), jen.ID("convertInputToClient").Call(jen.ID("input")), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to read one that doesn't exist"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					utils.CreateCtx(),
					jen.Line(),
					jen.Comment("Fetch oauth2Client"),
					jen.List(jen.ID("_"), jen.Err()).Assign().ID("testClient").Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("nonexistentID")),
					jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					utils.CreateCtx(),
					jen.Line(),
					jen.Comment("Create oauth2Client"),
					jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(jen.ID("t"), jen.ID("x").Dot("Username"), jen.ID("y").Dot("Password"), jen.ID("x").Dot("TwoFactorSecret")),
					jen.List(jen.ID("c"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("c"), jen.Err()),
					jen.Line(),
					jen.Comment("Fetch oauth2Client"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("testClient").Dot("GetOAuth2Client").Call(utils.CtxVar(), jen.ID("c").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert oauth2Client equality"),
					jen.ID("checkOAuth2ClientEquality").Call(jen.ID("t"), jen.ID("convertInputToClient").Call(jen.ID("input")), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					utils.CreateCtx(),
					jen.Line(),
					jen.Comment("Create oauth2Client"),
					jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(jen.ID("t"), jen.ID("x").Dot("Username"), jen.ID("y").Dot("Password"), jen.ID("x").Dot("TwoFactorSecret")),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("should be unable to authorize after being deleted"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					utils.CreateCtx(),
					jen.Line(),
					jen.Comment("create user"),
					jen.List(jen.ID("createdUser"), jen.ID("createdUserInput"), jen.ID("_")).Assign().ID("buildDummyUser").Call(jen.ID("test")),
					jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
					jen.Line(),
					jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(
						jen.ID("test"),
						jen.ID("createdUserInput").Dot("Username"),
						jen.ID("createdUserInput").Dot("Password"),
						jen.ID("createdUser").Dot("TwoFactorSecret"),
					),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("ArchiveHandler oauth2Client"),
					jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
					jen.Line(),
					jen.List(jen.ID("c2"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
						utils.CtxVar(),
						jen.ID("premade").Dot("ClientID"),
						jen.ID("premade").Dot("ClientSecret"),
						jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
						jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
						jen.ID("buildHTTPClient").Call(),
						jen.ID("premade").Dot("Scopes"),
						jen.ID("true"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("c2"), jen.Err()),
					jen.Line(),
					jen.List(jen.ID("_"), jen.Err()).Equals().ID("c2").Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Nil()),
					jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err(), jen.Lit("expected error from what should be an unauthorized client")),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					utils.CreateCtx(),
					jen.Line(),
					jen.Comment("Create oauth2Clients"),
					jen.Var().ID("expected").Index().Op("*").Qual(proj.ModelsV1Package(), "OAuth2Client"),
					jen.For(jen.ID("i").Assign().Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
						jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(
							jen.ID("t"),
							jen.ID("x").Dot("Username"),
							jen.ID("y").Dot("Password"),
							jen.ID("x").Dot("TwoFactorSecret"),
						),
						jen.List(jen.ID("oac"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("cookie"), jen.ID("input")),
						jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("oac"), jen.Err()),
						jen.ID("expected").Equals().ID("append").Call(jen.ID("expected"), jen.ID("oac")),
					),
					jen.Line(),
					jen.Comment("Assert oauth2Client list equality"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("testClient").Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Nil()),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Qual("github.com/stretchr/testify/assert", "True").Callln(
						jen.ID("t"),
						jen.ID("len").Call(jen.ID("actual").Dot("Clients")).Op("-").ID("len").Call(jen.ID("expected")).Op(">").Lit(0),
						jen.Lit("expected %d - %d to be > 0"),
						jen.ID("len").Call(jen.ID("actual").Dot("Clients")),
						jen.ID("len").Call(jen.ID("expected")),
					),
					jen.Line(),
					jen.For(jen.List(jen.ID("_"), jen.ID("oAuth2Client")).Assign().Range().ID("expected")).Block(
						jen.ID("clientFound").Assign().ID("false"),
						jen.For(jen.List(jen.ID("_"), jen.ID("c")).Assign().Range().ID("actual").Dot("Clients")).Block(
							jen.If(jen.ID("c").Dot("ID").Op("==").ID("oAuth2Client").Dot("ID")).Block(
								jen.ID("clientFound").Equals().ID("true"),
								jen.Break(),
							),
						),
						utils.AssertTrue(jen.ID("clientFound"), jen.Lit("expected oAuth2Client ID %d to be present in results"), jen.ID("oAuth2Client").Dot("ID")),
					),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.For(jen.List(jen.ID("_"), jen.ID("oa2c")).Assign().Range().ID("expected")).Block(
						jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("oa2c").Dot("ID")),
						jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(
							jen.ID("t"),
							jen.Err(),
							jen.Lit("error deleting client %d: %v"),
							jen.ID("oa2c").Dot("ID"),
							jen.Err(),
						),
					),
				)),
			)),
		),
		jen.Line(),
	)
	return ret
}
