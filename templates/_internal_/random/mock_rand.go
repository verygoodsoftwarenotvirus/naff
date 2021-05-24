package random

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockRandDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("Generator").Op("=").Parens(jen.Op("*").ID("MockGenerator")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("MockGenerator").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("MockGenerator")).ID("GenerateBase32EncodedString").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("length"),
			),
			jen.Return().List(jen.ID("args").Dot("String").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateBase64EncodedString implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockGenerator")).ID("GenerateBase64EncodedString").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("length"),
			),
			jen.Return().List(jen.ID("args").Dot("String").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateRawBytes implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockGenerator")).ID("GenerateRawBytes").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.Index().ID("byte"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("length"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("byte")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
