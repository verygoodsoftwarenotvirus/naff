package model

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersDotGo() *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.ID("fake").Dot("Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("mustBuildCode").Params(jen.ID("totpSecret").ID("string")).Params(jen.ID("string")).Block(
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(jen.ID("totpSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Return().ID("code"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("RandomUserInput creates a random UserInput"),
		jen.Line(),
		jen.Func().ID("RandomUserInput").Params().Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserInput")).Block(
			jen.Comment("I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,"),
			jen.Comment("could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit."),
			jen.ID("username").Op(":=").ID("fake").Dot("UserName").Call().Op("+").ID("fake").Dot("HexColor").Call().Op("+").ID("fake").Dot("Country").Call(),
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserInput").Valuesln(
				jen.ID("Username").Op(":").ID("username"),
				jen.ID("Password").Op(":").ID("fake").Dot("Password").Call(jen.Lit(64), jen.Lit(128), jen.ID("true"), jen.ID("true"), jen.ID("true")),
			),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)
	return ret
}
