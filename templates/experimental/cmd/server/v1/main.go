package v1

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mainDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("main").Params().Block(
		jen.ID("logger").Op(":=").ID("zerolog").Dot(
			"NewZeroLogger",
		).Call(),
		jen.ID("configFilepath").Op(":=").Qual("os", "Getenv").Call(jen.Lit("CONFIGURATION_FILEPATH")),
		jen.If(jen.ID("configFilepath").Op("==").Lit("")).Block(
			jen.ID("logger").Dot(
				"Fatal",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("no configuration file provided"))),
		),
		jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("config").Dot(
			"ParseConfigFile",
		).Call(jen.ID("configFilepath")),
		jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("cfg").Op("==").ID("nil")).Block(
			jen.ID("logger").Dot(
				"Fatal",
			).Call(jen.ID("err")),
		),
		jen.List(jen.ID("tctx"), jen.ID("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Qual("context", "Background").Call(), jen.ID("cfg").Dot(
			"Meta",
		).Dot(
			"StartupDeadline",
		)),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.Lit("initialization")),
		jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("cfg").Dot(
			"ProvideDatabase",
		).Call(jen.ID("ctx"), jen.ID("logger")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot(
				"Fatal",
			).Call(jen.ID("err")),
		),
		jen.List(jen.ID("server"), jen.ID("err")).Op(":=").ID("BuildServer").Call(jen.ID("ctx"), jen.ID("cfg"), jen.ID("logger"), jen.ID("db")),
		jen.ID("span").Dot(
			"End",
		).Call(),
		jen.ID("cancel").Call(),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Qual("log", "Fatal").Call(jen.ID("err")),
		),
		jen.ID("server").Dot(
			"Serve",
		).Call(),
	),

		jen.Line(),
	)
	return ret
}
