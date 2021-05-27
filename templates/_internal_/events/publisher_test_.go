package events

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func publisherTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildTestPublisher").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.ID("Publisher")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderMemory"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
			jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("ProvidePublisher").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("cfg"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("p"),
			),
			jen.Return().ID("p"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvidePublisher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderMemory"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("ProvidePublisher").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("p"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("ProvidePublisher").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("p"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with disabled config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderMemory"), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("false")),
					jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("ProvidePublisher").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("p"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error initializing topic"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").Lit(""), jen.ID("Topic").Op(":").ID("t").Dot("Name").Call(), jen.ID("Enabled").Op(":").ID("true")),
					jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("ProvidePublisher").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("p"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_publisher_PublishEvent").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("p").Op(":=").ID("buildTestPublisher").Call(jen.ID("t")),
					jen.ID("x").Op(":=").Struct(jen.ID("Name").ID("string")).Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("err").Op(":=").ID("p").Dot("PublishEvent").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("x"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("p").Op(":=").ID("buildTestPublisher").Call(jen.ID("t")),
					jen.ID("x").Op(":=").Struct(jen.ID("Name").Qual("encoding/json", "Number")).Valuesln(jen.ID("Name").Op(":").Qual("encoding/json", "Number").Call(jen.ID("t").Dot("Name").Call())),
					jen.ID("err").Op(":=").ID("p").Dot("PublishEvent").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("x"),
						jen.ID("nil"),
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
		jen.Func().ID("TestNoopEventPublisher_PublishEvent").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call()),
			),
		),
		jen.Line(),
	)

	return code
}
