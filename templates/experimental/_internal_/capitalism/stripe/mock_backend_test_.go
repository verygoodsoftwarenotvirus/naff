package stripe

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockBackendTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("stripe").Dot("Backend").Op("=").Parens(jen.Op("*").ID("mockBackend")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("mockBackend").Struct(
			jen.ID("mock").Dot("Mock"),
			jen.ID("anticipatedReturns").Index().Index().ID("byte"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockBackend")).ID("AnticipateCall").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("v").Interface()).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("b"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("v")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("m").Dot("anticipatedReturns").Op("=").ID("append").Call(
				jen.ID("m").Dot("anticipatedReturns"),
				jen.ID("b"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockBackend")).ID("Call").Params(jen.List(jen.ID("method"), jen.ID("path"), jen.ID("key")).ID("string"), jen.ID("params").ID("stripe").Dot("ParamsContainer"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.ID("b").Op(":=").ID("m").Dot("anticipatedReturns").Index(jen.Lit(0)),
			jen.ID("m").Dot("anticipatedReturns").Op("=").ID("append").Call(
				jen.ID("m").Dot("anticipatedReturns").Index(jen.Empty(), jen.Lit(0)),
				jen.ID("m").Dot("anticipatedReturns").Index(jen.Lit(1), jen.Empty()).Op("..."),
			),
			jen.If(jen.ID("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
				jen.ID("b"),
				jen.ID("v"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("method"),
				jen.ID("path"),
				jen.ID("key"),
				jen.ID("params"),
				jen.ID("v"),
			).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockBackend")).ID("CallRaw").Params(jen.List(jen.ID("method"), jen.ID("path"), jen.ID("key")).ID("string"), jen.ID("body").Op("*").ID("form").Dot("Values"), jen.ID("params").Op("*").ID("stripe").Dot("Params"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.ID("b").Op(":=").ID("m").Dot("anticipatedReturns").Index(jen.Lit(0)),
			jen.ID("m").Dot("anticipatedReturns").Op("=").ID("append").Call(
				jen.ID("m").Dot("anticipatedReturns").Index(jen.Empty(), jen.Lit(0)),
				jen.ID("m").Dot("anticipatedReturns").Index(jen.Lit(1), jen.Empty()).Op("..."),
			),
			jen.If(jen.ID("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
				jen.ID("b"),
				jen.ID("v"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("method"),
				jen.ID("path"),
				jen.ID("key"),
				jen.ID("body"),
				jen.ID("params"),
				jen.ID("v"),
			).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockBackend")).ID("CallMultipart").Params(jen.List(jen.ID("method"), jen.ID("path"), jen.ID("key"), jen.ID("boundary")).ID("string"), jen.ID("body").Op("*").Qual("bytes", "Buffer"), jen.ID("params").Op("*").ID("stripe").Dot("Params"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.ID("b").Op(":=").ID("m").Dot("anticipatedReturns").Index(jen.Lit(0)),
			jen.ID("m").Dot("anticipatedReturns").Op("=").ID("append").Call(
				jen.ID("m").Dot("anticipatedReturns").Index(jen.Empty(), jen.Lit(0)),
				jen.ID("m").Dot("anticipatedReturns").Index(jen.Lit(1), jen.Empty()).Op("..."),
			),
			jen.If(jen.ID("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
				jen.ID("b"),
				jen.ID("v"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("method"),
				jen.ID("path"),
				jen.ID("key"),
				jen.ID("boundary"),
				jen.ID("body"),
				jen.ID("params"),
				jen.ID("v"),
			).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockBackend")).ID("SetMaxNetworkRetries").Params(jen.ID("maxNetworkRetries").ID("int")).Body(
			jen.ID("m").Dot("Called").Call(jen.ID("maxNetworkRetries"))),
		jen.Line(),
	)

	return code
}
