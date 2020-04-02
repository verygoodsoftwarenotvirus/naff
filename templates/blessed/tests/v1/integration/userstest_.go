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
			jen.ID("b").Assign().ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
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
		jen.Func().ID("randString").Params().Params(jen.ID("string"), jen.ID("error")).Block(
			jen.ID("b").Assign().ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.Comment("Note that err == nil only if we read len(b) bytes"),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Lit(""), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyUserInput").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "UserInput")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
			jen.ID("userInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
				jen.ID("Username").MapAssign().Qual(utils.FakeLibrary, "Username").Call(),
				jen.ID("Password").MapAssign().Qual(utils.FakeLibrary, "Password").Call(jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.Lit(64)),
			),
			jen.Line(),
			jen.Return().ID("userInput"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyUser").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "UserCreationResponse"), jen.Op("*").Qual(proj.ModelsV1Package(), "UserInput"), jen.ParamPointer().Qual("net/http", "Cookie")).Block(
			jen.ID("t").Dot("Helper").Call(),
			utils.CreateCtx(),
			jen.Line(),
			jen.Comment("build user creation route input"),
			jen.ID("userInput").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("userInput")),
			jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("user")),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Line(),
			jen.If(jen.ID("user").Op("==").ID("nil").Op("||").ID("err").DoesNotEqual().ID("nil")).Block(
				jen.ID("t").Dot("FailNow").Call(),
			),
			jen.ID("cookie").Assign().ID("loginUser").Call(
				jen.ID("t"),
				jen.ID("userInput").Dot("Username"),
				jen.ID("userInput").Dot("Password"),
				jen.ID("user").Dot("TwoFactorSecret"),
			),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("cookie")),
			jen.Line(),
			jen.Return().List(jen.ID("user"), jen.ID("userInput"), jen.ID("cookie")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("checkUserCreationEquality").Params(jen.ID("t").ParamPointer().Qual("testing", "T"), jen.ID("expected").Op("*").Qual(proj.ModelsV1Package(), "UserInput"), jen.ID("actual").Op("*").Qual(proj.ModelsV1Package(), "UserCreationResponse")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ID")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Username"),
				jen.ID("actual").Dot("Username"),
			),
			jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("actual").Dot("TwoFactorSecret")),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("CreatedOn")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual").Dot("UpdatedOn")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual").Dot("ArchivedOn")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("checkUserEquality").Params(jen.ID("t").ParamPointer().Qual("testing", "T"), jen.ID("expected").Op("*").Qual(proj.ModelsV1Package(), "UserInput"), jen.ID("actual").Op("*").Qual(proj.ModelsV1Package(), "User")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ID")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Username"),
				jen.ID("actual").Dot("Username"),
			),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("CreatedOn")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual").Dot("UpdatedOn")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual").Dot("ArchivedOn")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestUsers").Params(jen.ID("test").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be creatable"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("tctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create user"),
					jen.ID("expected").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(jen.ID("tctx"), jen.VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
						jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
						jen.ID("Password").MapAssign().ID("expected").Dot("Password"))),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert user equality"),
					jen.ID("checkUserCreationEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveUser").Call(jen.ID("tctx"), jen.ID("actual").Dot("ID"))),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to read something that doesn't exist"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("tctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Fetch user"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetUser").Call(jen.ID("tctx"), jen.ID("nonexistentID")),
					jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
					jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("tctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create user"),
					jen.ID("expected").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
					jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(jen.ID("tctx"), jen.VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
						jen.ID("Username").MapAssign().ID("expected").Dot("Username"),
						jen.ID("Password").MapAssign().ID("expected").Dot("Password"))),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("premade").Dot("TwoFactorSecret")),
					jen.Line(),
					jen.Comment("Fetch user"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetUser").Call(jen.ID("tctx"), jen.ID("premade").Dot("ID")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID("t").Dot("Logf").Call(jen.Lit("error encountered trying to fetch user %q: %v\n"), jen.ID("premade").Dot("Username"), jen.Err()),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert user equality"),
					jen.ID("checkUserEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveUser").Call(jen.ID("tctx"), jen.ID("actual").Dot("ID"))),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("tctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create user"),
					jen.ID("y").Assign().ID("buildDummyUserInput").Call(jen.ID("t")),
					jen.List(jen.ID("u"), jen.Err()).Assign().ID("todoClient").Dot("CreateUser").Call(jen.ID("tctx"), jen.ID("y")),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
					jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("u")),
					jen.Line(),
					jen.If(jen.ID("u").Op("==").ID("nil").Op("||").ID("err").DoesNotEqual().ID("nil")).Block(
						jen.ID("t").Dot("Log").Call(jen.Lit("something has gone awry, user returned is nil")),
						jen.ID("t").Dot("FailNow").Call(),
					),
					jen.Line(),
					jen.Comment("Execute"),
					jen.Err().Equals().ID("todoClient").Dot("ArchiveUser").Call(jen.ID("tctx"), jen.ID("u").Dot("ID")),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
					jen.ID("tctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create users"),
					jen.Var().ID("expected").Index().Op("*").Qual(proj.ModelsV1Package(), "UserCreationResponse"),
					jen.For(jen.ID("i").Assign().Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
						jen.List(jen.ID("user"), jen.ID("_"), jen.ID("c")).Assign().ID("buildDummyUser").Call(jen.ID("t")),
						jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("c")),
						jen.ID("expected").Equals().ID("append").Call(jen.ID("expected"), jen.ID("user")),
					),
					jen.Line(),
					jen.Comment("Assert user list equality"),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dot("GetUsers").Call(jen.ID("tctx"), jen.Nil()),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Users"))),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.For(jen.List(jen.ID("_"), jen.ID("user")).Assign().Range().ID("actual").Dot("Users")).Block(
						jen.Err().Equals().ID("todoClient").Dot("ArchiveUser").Call(jen.ID("tctx"), jen.ID("user").Dot("ID")),
						jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
					),
				)),
			)),
		),
		jen.Line(),
	)
	return ret
}
