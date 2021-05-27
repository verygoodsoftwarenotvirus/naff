package random

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func randTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("erroneousReader").Struct(),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("r").Op("*").ID("erroneousReader")).ID("Read").Params(jen.ID("p").Index().ID("byte")).Params(jen.ID("n").ID("int"), jen.ID("err").ID("error")).Body(
			jen.Return().List(jen.Op("-").Lit(1), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestGenerateBase32EncodedString").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("GenerateBase32EncodedString").Call(
						jen.ID("ctx"),
						jen.Lit(32),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestGenerateBase64EncodedString").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("GenerateBase64EncodedString").Call(
						jen.ID("ctx"),
						jen.Lit(32),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestGenerateRawBytes").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("GenerateRawBytes").Call(
						jen.ID("ctx"),
						jen.Lit(32),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestStandardSecretGenerator_GenerateBase32EncodedString").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleLength").Op(":=").Lit(123),
					jen.ID("s").Op(":=").ID("NewGenerator").Call(jen.ID("nil")),
					jen.List(jen.ID("value"), jen.ID("err")).Op(":=").ID("s").Dot("GenerateBase32EncodedString").Call(
						jen.ID("ctx"),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("value"),
					),
					jen.ID("assert").Dot("Greater").Call(
						jen.ID("t"),
						jen.ID("len").Call(jen.ID("value")),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from secure PRNG"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleLength").Op(":=").Lit(123),
					jen.List(jen.ID("s"), jen.ID("ok")).Op(":=").ID("NewGenerator").Call(jen.ID("nil")).Assert(jen.Op("*").ID("standardGenerator")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.ID("s").Dot("randReader").Op("=").Op("&").ID("erroneousReader").Values(),
					jen.List(jen.ID("value"), jen.ID("err")).Op(":=").ID("s").Dot("GenerateBase32EncodedString").Call(
						jen.ID("ctx"),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("value"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestStandardSecretGenerator_GenerateBase64EncodedString").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleLength").Op(":=").Lit(123),
					jen.ID("s").Op(":=").ID("NewGenerator").Call(jen.ID("nil")),
					jen.List(jen.ID("value"), jen.ID("err")).Op(":=").ID("s").Dot("GenerateBase64EncodedString").Call(
						jen.ID("ctx"),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("value"),
					),
					jen.ID("assert").Dot("Greater").Call(
						jen.ID("t"),
						jen.ID("len").Call(jen.ID("value")),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from secure PRNG"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleLength").Op(":=").Lit(123),
					jen.List(jen.ID("s"), jen.ID("ok")).Op(":=").ID("NewGenerator").Call(jen.ID("nil")).Assert(jen.Op("*").ID("standardGenerator")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.ID("s").Dot("randReader").Op("=").Op("&").ID("erroneousReader").Values(),
					jen.List(jen.ID("value"), jen.ID("err")).Op(":=").ID("s").Dot("GenerateBase64EncodedString").Call(
						jen.ID("ctx"),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("value"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestStandardSecretGenerator_GenerateRawBytes").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleLength").Op(":=").Lit(123),
					jen.ID("s").Op(":=").ID("NewGenerator").Call(jen.ID("nil")),
					jen.List(jen.ID("value"), jen.ID("err")).Op(":=").ID("s").Dot("GenerateRawBytes").Call(
						jen.ID("ctx"),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("value"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("len").Call(jen.ID("value")),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from secure PRNG"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleLength").Op(":=").Lit(123),
					jen.List(jen.ID("s"), jen.ID("ok")).Op(":=").ID("NewGenerator").Call(jen.ID("nil")).Assert(jen.Op("*").ID("standardGenerator")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.ID("s").Dot("randReader").Op("=").Op("&").ID("erroneousReader").Values(),
					jen.List(jen.ID("value"), jen.ID("err")).Op(":=").ID("s").Dot("GenerateRawBytes").Call(
						jen.ID("ctx"),
						jen.ID("exampleLength"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("value"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
