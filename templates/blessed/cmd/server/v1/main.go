package v1

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("main")

	internalConfigImp := fmt.Sprintf("%s/internal/v1/config", proj.OutputPath)
	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("main").Params().Block(
			jen.Comment("initialize our logger of choice"),
			jen.ID("logger").Assign().Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.Line(),
			jen.Comment("find and validate our configuration filepath"),
			jen.ID("configFilepath").Assign().Qual("os", "Getenv").Call(jen.Lit("CONFIGURATION_FILEPATH")),
			jen.If(jen.ID("configFilepath").Op("==").EmptyString()).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Qual("errors", "New").Call(jen.Lit("no configuration file provided"))),
			),
			jen.Line(),
			jen.Comment("parse our config file"),
			jen.List(jen.ID("cfg"), jen.Err()).Assign().Qual(internalConfigImp, "ParseConfigFile").Call(jen.ID("configFilepath")),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID("cfg").Op("==").ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.Comment("only allow initialization to take so long"),
			jen.List(utils.CtxVar(), jen.ID("cancel")).Assign().Qual("context", "WithTimeout").Call(jen.Qual("context", "Background").Call(), jen.ID("cfg").Dot("Meta").Dot("StartupDeadline")),
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("initialization")),
			jen.Line(),
			jen.Comment("connect to our database"),
			jen.List(jen.ID("db"), jen.Err()).Assign().ID("cfg").Dot("ProvideDatabase").Call(utils.CtxVar(), jen.ID("logger")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.Comment("build our server struct"),
			jen.List(jen.ID("server"), jen.Err()).Assign().ID("BuildServer").Call(utils.CtxVar(), jen.ID("cfg"), jen.ID("logger"), jen.ID("db")),
			jen.ID("span").Dot("End").Call(),
			jen.ID("cancel").Call(),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Qual("log", "Fatal").Call(jen.Err()),
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
