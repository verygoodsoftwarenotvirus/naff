package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			utils.InlineFakeSeedFunc(),
			jen.Line(),
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("randString produces a random string"),
		jen.Line(),
		jen.Comment("https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/"),
		jen.Line(),
		jen.Func().ID("randString").Params().Params(jen.String(), jen.Error()).Block(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.Comment("Note that err == nil only if we read len(b) bytes"),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.EmptyString(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyUserInput").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserInput")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
			jen.ID("userInput").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
				jen.ID("Username").MapAssign().Qual(utils.FakeLibrary, "Username").Call(),
				jen.ID("Password").MapAssign().Qual(utils.FakeLibrary, "Password").Call(jen.True(), jen.True(), jen.True(), jen.True(), jen.True(), jen.Lit(64)),
			),
			jen.Line(),
			jen.Return().ID("userInput"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyUser").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserCreationResponse"), jen.PointerTo().Qual(proj.ModelsV1Package(), "UserInput"), jen.PointerTo().Qual("net/http", "Cookie")).Block(
			jen.ID("t").Dot("Helper").Call(),
			utils.CreateCtx(),
			jen.Line(),
			jen.Comment("build user creation route input"),
			jen.ID("userInput").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("userInput")),
			utils.AssertNotNil(jen.ID("user"), nil),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			jen.If(jen.ID("user").Op("==").ID("nil").Or().Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("t").Dot("FailNow").Call(),
			),
			jen.ID("cookie").Assign().ID("loginUser").Call(
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
	)

	ret.Add(
		jen.Func().ID("checkUserCreationEquality").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.ID("expected").PointerTo().Qual(proj.ModelsV1Package(), "UserInput"), jen.ID("actual").PointerTo().Qual(proj.ModelsV1Package(), "UserCreationResponse")).Block(
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
			utils.AssertNil(jen.ID("actual").Dot("UpdatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("ArchivedOn"), nil),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("checkUserEquality").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.ID("expected").PointerTo().Qual(proj.ModelsV1Package(), "UserInput"), jen.ID("actual").PointerTo().Qual(proj.ModelsV1Package(), "User")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			utils.AssertNotZero(jen.ID("actual").Dot("ID"), nil),
			utils.AssertEqual(
				jen.ID("expected").Dot("Username"),
				jen.ID("actual").Dot("Username"),
				nil,
			),
			utils.AssertNotZero(jen.ID("actual").Dot("CreatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("UpdatedOn"), nil),
			utils.AssertNil(jen.ID("actual").Dot("ArchivedOn"), nil),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestUsers").Params(jen.ID("test").PointerTo().Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be creatable",
					jen.Line(),
					jen.Comment("Create user"),
					jen.ID("expected").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(utils.CtxVar(), jen.AddressOf().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
						jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
						jen.ID("Password").MapAssign().ID("expected").Dot("Password"))),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert user equality"),
					jen.ID("checkUserCreationEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					utils.AssertNoError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("actual").Dot("ID")), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to read something that doesn't exist",
					jen.Line(),
					jen.Comment("Fetch user"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetUser").Call(utils.CtxVar(), jen.ID("nonexistentID")),
					utils.AssertNil(jen.ID("actual"), nil),
					utils.AssertError(jen.Err(), nil),
				),
				jen.Line(),
				utils.BuildSubTest(
					"it should be readable",
					jen.Line(),
					jen.Comment("Create user"),
					jen.ID("expected").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(utils.CtxVar(), jen.AddressOf().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
						jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
						jen.ID("Password").MapAssign().ID("expected").Dot("Password"))),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					utils.AssertNotEmpty(jen.ID("premade").Dot("TwoFactorSecret"), nil),
					jen.Line(),
					jen.Comment("Fetch user"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetUser").Call(utils.CtxVar(), jen.ID("premade").Dot("ID")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID("t").Dot("Logf").Call(jen.Lit("error encountered trying to fetch user %q: %v\n"), jen.ID("premade").Dot("Username"), jen.Err()),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert user equality"),
					jen.ID("checkUserEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					utils.AssertNoError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("actual").Dot("ID")), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be deleted",
					jen.Line(),
					jen.Comment("Create user"),
					jen.ID("y").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
					jen.List(jen.ID("u"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("y")),
					utils.AssertNoError(jen.Err(), nil),
					utils.AssertNotNil(jen.ID("u"), nil),
					jen.Line(),
					jen.If(jen.ID("u").Op("==").ID("nil").Or().Err().DoesNotEqual().ID("nil")).Block(
						jen.ID("t").Dot("Log").Call(jen.Lit("something has gone awry, user returned is nil")),
						jen.ID("t").Dot("FailNow").Call(),
					),
					jen.Line(),
					jen.Comment("Execute"),
					jen.Err().Equals().ID("todoClient").Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("u").Dot("ID")),
					utils.AssertNoError(jen.Err(), nil),
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be read in a list",
					jen.Line(),
					jen.Comment("Create users"),
					jen.Var().ID("expected").Index().PointerTo().Qual(proj.ModelsV1Package(), "UserCreationResponse"),
					jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
						jen.List(jen.ID("user"), jen.Underscore(), jen.ID("c")).Assign().ID("buildDummyUser").Call(jen.ID("t")),
						utils.AssertNotNil(jen.ID("c"), nil),
						jen.ID("expected").Equals().ID("append").Call(jen.ID("expected"), jen.ID("user")),
					),
					jen.Line(),
					jen.Comment("Assert user list equality"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetUsers").Call(utils.CtxVar(), jen.Nil()),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					utils.AssertTrue(jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Users")), nil),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.For(jen.List(jen.Underscore(), jen.ID("user")).Assign().Range().ID("actual").Dot("Users")).Block(
						jen.Err().Equals().ID("todoClient").Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("user").Dot("ID")),
						utils.AssertNoError(jen.Err(), nil),
					),
				),
			)),
		),
		jen.Line(),
	)
	return ret
}
