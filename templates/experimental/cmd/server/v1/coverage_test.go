package v1

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func coverageTestDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestRunMain").Params(jen.ID("_").Op("*").Qual("testing", "T")).Block(
		jen.List(jen.ID("d"), jen.ID("err")).Op(":=").Qual("time", "ParseDuration").Call(jen.Qual("os", "Getenv").Call(jen.Lit("RUNTIME_DURATION"))),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Qual("log", "Fatal").Call(jen.ID("err")),
		),
		jen.Go().ID("main").Call(),
		jen.Qual("time", "Sleep").Call(jen.ID("d")),
	),

		jen.Line(),
	)
	return ret
}
