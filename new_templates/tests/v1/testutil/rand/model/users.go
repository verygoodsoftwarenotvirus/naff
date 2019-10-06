package randmodel

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Func().ID("init").Params().Block(
		jen.ID("fake").Dot(
			"Seed",
		).Call(jen.Qual("time", "Now").Call().Dot(
			"UnixNano",
		).Call()),
	),
	)
	ret.Add(jen.Func().ID("mustBuildCode").Params(jen.ID("totpSecret").ID("string")).Params(jen.ID("string")).Block(
		jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot(
			"GenerateCode",
		).Call(jen.ID("totpSecret"), jen.Qual("time", "Now").Call().Dot(
			"UTC",
		).Call()),
		jen.If(
			jen.ID("err").Op("!=").ID("nil"),
		).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
		jen.Return().ID("code"),
	),
	)
	ret.Add(jen.Func())
	return ret
}
