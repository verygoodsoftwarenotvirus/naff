package model

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("mustBuildCode").Params(jen.ID("totpSecret").ID("string")).Params(jen.ID("string")).Block(
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("totpSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Return().ID("code"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("RandomUserInput creates a random UserInput"),
		jen.Line(),
		jen.Func().ID("RandomUserInput").Params().Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput")).Block(
			jen.Comment("I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,"),
			jen.Comment("could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit."),
			jen.ID("username").Assign().Qual(utils.FakeLibrary, "Username").Call().Op("+").Qual(utils.FakeLibrary, "HexColor").Call().Op("+").Qual(utils.FakeLibrary, "Country").Call(),
			jen.ID("x").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput").Valuesln(
				jen.ID("Username").MapAssign().ID("username"),
				jen.ID("Password").MapAssign().Qual(utils.FakeLibrary, "Password").Call(jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.Lit(64)),
			),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)
	return ret
}
