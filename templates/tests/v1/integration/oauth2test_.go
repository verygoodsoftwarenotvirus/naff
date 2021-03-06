package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2TestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("integration")

	utils.AddImports(proj, code)

	code.Add(buildOAuth2TestMustBuildCode()...)
	code.Add(buildOAuth2TestBuildDummyOAuth2ClientInput(proj)...)
	code.Add(buildOAuth2TestConvertInputToClient(proj)...)
	code.Add(buildOAuth2TestCheckOAuth2ClientEquality(proj)...)
	code.Add(buildOAuth2TestTestOAuth2Clients(proj)...)

	return code
}

func buildOAuth2TestMustBuildCode() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("mustBuildCode").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.ID("totpSecret").String()).Params(jen.String()).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("totpSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			utils.RequireNoError(jen.Err(), nil),
			jen.Return().ID("code"),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2TestBuildDummyOAuth2ClientInput(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildDummyOAuth2ClientInput").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpToken")).String()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
				jen.ID("UserLoginInput").MapAssign().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("username"),
					jen.ID("Password").MapAssign().ID("password"),
					jen.ID("TOTPToken").MapAssign().ID("mustBuildCode").Call(jen.ID("t"), jen.ID("totpToken")),
				),
				jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("*")),
				jen.ID("RedirectURI").MapAssign().Lit("http://localhost"),
			),
			jen.Line(),
			jen.Return().ID("x"),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2TestConvertInputToClient(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("convertInputToClient").Params(jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Body(
			jen.Return().AddressOf().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
				jen.ID("ClientID").MapAssign().ID("input").Dot("ClientID"),
				jen.ID("ClientSecret").MapAssign().ID("input").Dot("ClientSecret"),
				jen.ID("RedirectURI").MapAssign().ID("input").Dot("RedirectURI"),
				jen.ID("Scopes").MapAssign().ID("input").Dot("Scopes"),
				jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("input").Dot(constants.UserOwnershipFieldName)),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2TestCheckOAuth2ClientEquality(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("checkOAuth2ClientEquality").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			utils.AssertNotZero(jen.ID("actual").Dot("ID"), nil),
			utils.AssertNotEmpty(jen.ID("actual").Dot("ClientID"), nil),
			utils.AssertNotEmpty(jen.ID("actual").Dot("ClientSecret"), nil),
			utils.AssertEqual(jen.ID("expected").Dot("RedirectURI"), jen.ID("actual").Dot("RedirectURI"), nil),
			utils.AssertEqual(jen.ID("expected").Dot("Scopes"), jen.ID("actual").Dot("Scopes"), nil),
			utils.AssertNotZero(jen.ID("actual").Dot("CreatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("ArchivedOn"), nil),
		),
		jen.Line(),
	}

	return lines
}

func buildOAuth2TestTestOAuth2Clients(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestOAuth2Clients").Params(jen.ID("test").PointerTo().Qual("testing", "T")).Body(
			jen.ID("_ctx").Assign().Add(constants.InlineCtx()),
			jen.Line(),
			jen.Comment("create user."),
			jen.List(jen.ID("x"), jen.ID("y"), jen.ID("cookie")).Assign().ID("buildDummyUser").Call(
				jen.ID("_ctx"),
				jen.ID("test"),
			),
			jen.Qual(constants.AssertPkg, "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
			jen.Line(),
			jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(
				jen.ID("test"),
				jen.ID("x").Dot("Username"),
				jen.ID("y").Dot("Password"),
				jen.ID("x").Dot("TwoFactorSecret"),
			),
			jen.List(jen.ID("premade"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateOAuth2Client").Call(jen.ID("_ctx"), jen.ID("cookie"), jen.ID("input")),
			jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.Err()),
			jen.Line(),
			jen.List(jen.ID("testClient"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
				jen.ID("_ctx"),
				jen.ID("premade").Dot("ClientID"),
				jen.ID("premade").Dot("ClientSecret"),
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
				jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("PlainClient").Call(),
				jen.ID("premade").Dot("Scopes"),
				jen.ID("debug"),
			),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("test"), jen.Err(), jen.Lit("error setting up auxiliary client")),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"should be creatable",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create oauth2Client."),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert oauth2Client equality."),
					jen.ID("checkOAuth2ClientEquality").Call(jen.ID("t"), jen.ID("convertInputToClient").Call(jen.ID("input")), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up."),
					jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(constants.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"it should return an error when trying to read one that doesn't exist",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Fetch oauth2Client."),
					jen.List(jen.Underscore(), jen.Err()).Assign().ID("testClient").Dot("GetOAuth2Client").Call(constants.CtxVar(), jen.ID("nonexistentID")),
					utils.AssertError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"it should be readable",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create oauth2Client."),
					jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(jen.ID("t"), jen.ID("x").Dot("Username"), jen.ID("y").Dot("Password"), jen.ID("x").Dot("TwoFactorSecret")),
					jen.List(jen.ID("c"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("c"), jen.Err()),
					jen.Line(),
					jen.Comment("Fetch oauth2Client."),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("testClient").Dot("GetOAuth2Client").Call(constants.CtxVar(), jen.ID("c").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert oauth2Client equality."),
					jen.ID("checkOAuth2ClientEquality").Call(jen.ID("t"), jen.ID("convertInputToClient").Call(jen.ID("input")), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up."),
					jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(constants.CtxVar(), jen.ID("actual").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"should be able to be deleted",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create oauth2Client."),
					jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(jen.ID("t"), jen.ID("x").Dot("Username"), jen.ID("y").Dot("Password"), jen.ID("x").Dot("TwoFactorSecret")),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Clean up."),
					jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(constants.CtxVar(), jen.ID("premade").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"should be unable to authorize after being deleted",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("create user."),
					jen.List(jen.ID("createdUser"), jen.ID("createdUserInput"), jen.Underscore()).Assign().ID("buildDummyUser").Call(
						constants.CtxVar(),
						jen.ID("test"),
					),
					jen.Qual(constants.AssertPkg, "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
					jen.Line(),
					jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(
						jen.ID("test"),
						jen.ID("createdUserInput").Dot("Username"),
						jen.ID("createdUserInput").Dot("Password"),
						jen.ID("createdUser").Dot("TwoFactorSecret"),
					),
					jen.List(jen.ID("premade"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("archive oauth2Client."),
					utils.RequireNoError(jen.ID("testClient").Dot("ArchiveOAuth2Client").Call(constants.CtxVar(), jen.ID("premade").Dot("ID")), nil),
					jen.Line(),
					jen.List(jen.ID("c2"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
						constants.CtxVar(),
						jen.ID("premade").Dot("ClientID"),
						jen.ID("premade").Dot("ClientSecret"),
						jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
						jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
						jen.ID("buildHTTPClient").Call(),
						jen.ID("premade").Dot("Scopes"),
						jen.True(),
					),
					jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("c2"), jen.Err()),
					jen.Line(),
					jen.List(jen.Underscore(), jen.Err()).Equals().ID("c2").Dot("GetOAuth2Clients").Call(constants.CtxVar(), jen.Nil()),
					utils.AssertError(jen.Err(), jen.Lit("expected error from what should be an unauthorized client")),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"should be able to be read in a list",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create oauth2Clients."),
					jen.Var().ID("expected").Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
					jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Body(
						jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(
							jen.ID("t"),
							jen.ID("x").Dot("Username"),
							jen.ID("y").Dot("Password"),
							jen.ID("x").Dot("TwoFactorSecret"),
						),
						jen.List(jen.ID("oac"), jen.Err()).Assign().ID("testClient").Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID("cookie"), jen.ID("input")),
						jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("oac"), jen.Err()),
						jen.ID("expected").Equals().ID("append").Call(jen.ID("expected"), jen.ID("oac")),
					),
					jen.Line(),
					jen.Comment("Assert oauth2Client list equality."),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("testClient").Dot("GetOAuth2Clients").Call(constants.CtxVar(), jen.Nil()),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Qual(constants.AssertPkg, "True").Callln(
						jen.ID("t"),
						jen.Len(jen.ID("actual").Dot("Clients")).Minus().ID("len").Call(jen.ID("expected")).GreaterThan().Zero(),
						jen.Lit("expected %d - %d to be > 0"),
						jen.Len(jen.ID("actual").Dot("Clients")),
						jen.Len(jen.ID("expected")),
					),
					jen.Line(),
					jen.For(jen.List(jen.Underscore(), jen.ID("oAuth2Client")).Assign().Range().ID("expected")).Body(
						jen.ID("clientFound").Assign().False(),
						jen.For(jen.List(jen.Underscore(), jen.ID("c")).Assign().Range().ID("actual").Dot("Clients")).Body(
							jen.If(jen.ID("c").Dot("ID").IsEqualTo().ID("oAuth2Client").Dot("ID")).Body(
								jen.ID("clientFound").Equals().True(),
								jen.Break(),
							),
						),
						utils.AssertTrue(jen.ID("clientFound"), jen.Lit("expected oAuth2Client ID %d to be present in results"), jen.ID("oAuth2Client").Dot("ID")),
					),
					jen.Line(),
					jen.Comment("Clean up."),
					jen.For(jen.List(jen.Underscore(), jen.ID("oa2c")).Assign().Range().ID("expected")).Body(
						jen.Err().Equals().ID("testClient").Dot("ArchiveOAuth2Client").Call(constants.CtxVar(), jen.ID("oa2c").Dot("ID")),
						utils.AssertNoError(
							jen.Err(),
							jen.Lit("error deleting client %d: %v"),
							jen.ID("oa2c").Dot("ID"),
							jen.Err(),
						),
					),
				),
			)),
		),
		jen.Line(),
	}

	return lines
}
