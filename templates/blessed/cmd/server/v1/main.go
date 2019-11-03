package v1

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("main")

	internalConfigImp := fmt.Sprintf("%s/internal/v1/config", pkgRoot)
	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("main").Params().Block(
			jen.Comment("initialize our logger of choice"),
			jen.ID("logger").Op(":=").ID("zerolog").Dot("NewZeroLogger").Call(),
			jen.Line(),
			jen.Comment("find and validate our configuration filepath"),
			jen.ID("configFilepath").Op(":=").Qual("os", "Getenv").Call(jen.Lit("CONFIGURATION_FILEPATH")),
			jen.If(jen.ID("configFilepath").Op("==").Lit("")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Qual("errors", "New").Call(jen.Lit("no configuration file provided"))),
			),
			jen.Line(),
			jen.Comment("parse our config file"),
			jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").Qual(internalConfigImp, "ParseConfigFile").Call(jen.ID("configFilepath")),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("cfg").Op("==").ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Comment("only allow initialization to take so long"),
			jen.List(jen.ID("tctx"), jen.ID("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Qual("context", "Background").Call(), jen.ID("cfg").Dot("Meta").Dot("StartupDeadline")),
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.Lit("initialization")),
			jen.Line(),
			jen.Comment("connect to our database"),
			jen.List(jen.ID("db"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideDatabase").Call(jen.ID("ctx"), jen.ID("logger")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Comment("build our server struct"),
			jen.List(jen.ID("server"), jen.ID("err")).Op(":=").ID("BuildServer").Call(jen.ID("ctx"), jen.ID("cfg"), jen.ID("logger"), jen.ID("db")),
			jen.ID("span").Dot("End").Call(),
			jen.ID("cancel").Call(),
			jen.Line(),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Qual("log", "Fatal").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Comment("I slept and dreamt that life was joy."),
			jen.Comment("  I awoke and saw that life was service."),
			jen.Comment("  	I acted and behold, service deployed."),
			jen.ID("server").Dot("Serve").Call(),
		),
		jen.Line(),
	)
	return ret
}
