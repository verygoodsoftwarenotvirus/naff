package logging

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func zerologLoggerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_buildZerologger").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("buildZerologger").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestNewLogger").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("NewZerologLogger").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_WithName").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("l").Dot("WithName").Call(jen.ID("t").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_SetLevel").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("l").Dot("SetLevel").Call(jen.ID("ErrorLevel")),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_SetRequestIDFunc").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("l").Dot("SetRequestIDFunc").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
						jen.Return().Lit(""))),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_Info").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("l").Dot("Info").Call(jen.ID("t").Dot("Name").Call()),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_Debug").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("l").Dot("Debug").Call(jen.ID("t").Dot("Name").Call()),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_Error").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("l").Dot("Error").Call(
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_Printf").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("l").Dot("Printf").Call(jen.ID("t").Dot("Name").Call()),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_Clone").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("l").Dot("Clone").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_WithValue").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("l").Dot("WithValue").Call(
							jen.Lit("name"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_WithValues").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("l").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("name").Op(":").ID("t").Dot("Name").Call())),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_WithError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("l").Dot("WithError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_WithRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("l"), jen.ID("ok")).Op(":=").ID("NewZerologLogger").Call().Assert(jen.Op("*").ID("zerologLogger")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.ID("l").Dot("requestIDFunc").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
						jen.Return().ID("t").Dot("Name").Call()),
					jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "ParseRequestURI").Call(jen.Lit("https://todo.verygoodsoftwarenotvirus.ru?things=stuff")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("l").Dot("WithRequest").Call(jen.Op("&").Qual("net/http", "Request").Valuesln(jen.ID("URL").Op(":").ID("u"))),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_zerologLogger_WithResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("NewZerologLogger").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("l").Dot("WithResponse").Call(jen.Op("&").Qual("net/http", "Response").Values()),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
