package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildUsersTestsInit()...)
	code.Add(buildUsersTestsRandString()...)
	code.Add(buildUsersTestsBuildDummyUser(proj)...)
	code.Add(buildUsersTestsCheckUserCreationEquality(proj)...)
	code.Add(buildUsersTestsCheckUserEquality(proj)...)
	code.Add(buildUsersTestsTestUsers(proj)...)

	return code
}

func buildUsersTestsInit() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersTestsRandString() []jen.Code {
	lines := []jen.Code{
		jen.Comment("randString produces a random string."),
		jen.Line(),
		jen.Comment("https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/"),
		jen.Line(),
		jen.Func().ID("randString").Params().Params(jen.String(), jen.Error()).Body(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.Comment("Note that err == nil only if we read len(b) bytes"),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.EmptyString(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersTestsBuildDummyUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildDummyUser").Params(
			constants.CtxParam(),
			jen.ID("t").PointerTo().Qual("testing", "T"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "UserCreationResponse"),
			jen.PointerTo().Qual(proj.TypesPackage(), "UserCreationInput"),
			jen.PointerTo().Qual("net/http", "Cookie"),
		).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.Comment("build user creation route input."),
			jen.ID("userInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserCreationInput").Call(),
			jen.List(jen.ID("user"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateUser").Call(
				constants.CtxVar(),
				jen.ID("userInput"),
			),
			utils.RequireNotNil(jen.ID("user"), nil),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			utils.RequireNoError(jen.Err(), nil),
			utils.RequireNoError(
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("VerifyTOTPSecret").Call(
					constants.CtxVar(),
					jen.ID("user").Dot("ID"),
					jen.ID("token"),
				),
				nil,
			),
			jen.Line(),
			jen.ID("cookie").Assign().ID("loginUser").Call(
				constants.CtxVar(),
				jen.ID("t"),
				jen.ID("userInput").Dot("Username"),
				jen.ID("userInput").Dot("Password"),
				jen.ID("user").Dot("TwoFactorSecret"),
			),
			jen.Line(),
			utils.RequireNoError(jen.Err(), nil),
			utils.RequireNotNil(jen.ID("cookie"), nil),
			jen.Line(),
			jen.Return().List(jen.ID("user"), jen.ID("userInput"), jen.ID("cookie")),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersTestsCheckUserCreationEquality(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("checkUserCreationEquality").Params(
			jen.ID("t").PointerTo().Qual("testing", "T"),
			jen.ID("expected").PointerTo().Qual(proj.TypesPackage(), "UserCreationInput"),
			jen.ID("actual").PointerTo().Qual(proj.TypesPackage(), "UserCreationResponse"),
		).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			utils.AssertNotZero(jen.ID("actual").Dot("ID"), nil),
			utils.AssertEqual(
				jen.ID("expected").Dot("Username"),
				jen.ID("actual").Dot("Username"),
				nil,
			),
			utils.AssertNotEmpty(jen.ID("actual").Dot("TwoFactorSecret"), nil),
			utils.AssertNotZero(jen.ID("actual").Dot("CreatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("LastUpdatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("ArchivedOn"), nil),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersTestsCheckUserEquality(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("checkUserEquality").Params(
			jen.ID("t").PointerTo().Qual("testing", "T"),
			jen.ID("expected").PointerTo().Qual(proj.TypesPackage(), "UserCreationInput"),
			jen.ID("actual").PointerTo().Qual(proj.TypesPackage(), "User"),
		).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			utils.AssertNotZero(jen.ID("actual").Dot("ID"), nil),
			utils.AssertEqual(
				jen.ID("expected").Dot("Username"),
				jen.ID("actual").Dot("Username"),
				nil,
			),
			utils.AssertNotZero(jen.ID("actual").Dot("CreatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("LastUpdatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("ArchivedOn"), nil),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersTestsTestUsers(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestUsers").Params(jen.ID("test").PointerTo().Qual("testing", "T")).Body(
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"should be creatable",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create user."),
					utils.BuildFakeVarWithCustomName(proj, "exampleUserInput", "BuildFakeUserCreationInput"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateUser").Call(
						constants.CtxVar(),
						jen.ID("exampleUserInput"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert user equality."),
					jen.ID("checkUserCreationEquality").Call(jen.ID("t"), jen.ID("exampleUserInput"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up."),
					utils.AssertNoError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveUser").Call(
						constants.CtxVar(),
						jen.ID("actual").Dot("ID")), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"it should return an error when trying to read something that doesn't exist",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Fetch user."),
					jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("GetUser").Call(
						constants.CtxVar(),
						jen.ID("nonexistentID"),
					),
					utils.AssertNil(jen.ID("actual"), nil),
					utils.AssertError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					"it should be readable",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create user."),
					utils.BuildFakeVarWithCustomName(proj, "exampleUserInput", "BuildFakeUserCreationInput"),
					jen.List(jen.ID("premade"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateUser").Call(
						constants.CtxVar(),
						jen.ID("exampleUserInput"),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					utils.AssertNotEmpty(jen.ID("premade").Dot("TwoFactorSecret"), nil),
					jen.Line(),
					jen.List(jen.ID("secretVerificationToken"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
						jen.ID("premade").Dot("TwoFactorSecret"),
						jen.Qual("time", "Now").Call().Dot("UTC").Call(),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("secretVerificationToken"), jen.Err()),
					jen.Line(),
					utils.AssertNoError(
						jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("VerifyTOTPSecret").Call(
							constants.CtxVar(),
							jen.ID("premade").Dot("ID"),
							jen.ID("secretVerificationToken"),
						),
						nil,
					),
					jen.Line(),
					jen.Comment("Fetch user."),
					jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("GetUser").Call(
						constants.CtxVar(),
						jen.ID("premade").Dot("ID"),
					),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("error encountered trying to fetch user %q: %v\n"),
							jen.ID("premade").Dot("Username"),
							jen.Err(),
						),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert user equality."),
					jen.ID("checkUserEquality").Call(jen.ID("t"), jen.ID("exampleUserInput"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up."),
					utils.AssertNoError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveUser").Call(
						constants.CtxVar(), jen.ID("actual").Dot("ID")), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"should be able to be deleted",
					utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
					jen.Line(),
					jen.Comment("Create user."),
					utils.BuildFakeVarWithCustomName(proj, "exampleUserInput", "BuildFakeUserCreationInput"),
					jen.List(jen.ID("u"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("CreateUser").Call(
						constants.CtxVar(),
						jen.ID("exampleUserInput"),
					),
					utils.AssertNoError(jen.Err(), nil),
					utils.AssertNotNil(jen.ID("u"), nil),
					jen.Line(),
					jen.If(jen.ID("u").IsEqualTo().ID("nil").Or().Err().DoesNotEqual().ID("nil")).Body(
						jen.ID("t").Dot("Log").Call(jen.Lit("something has gone awry, user returned is nil")),
						jen.ID("t").Dot("FailNow").Call(),
					),
					jen.Line(),
					jen.Comment("Execute."),
					jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveUser").Call(constants.CtxVar(), jen.ID("u").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
		),
		jen.Line(),
	}

	return lines
}
