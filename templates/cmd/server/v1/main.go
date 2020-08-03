package v1

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("main")

	utils.AddImports(proj, code)

	code.Add(buildMain(proj)...)

	return code
}

func buildMain(proj *models.Project) []jen.Code {
	internalConfigImp := fmt.Sprintf("%s/internal/v1/config", proj.OutputPath)

	return []jen.Code{
		jen.Func().ID("main").Params().Body(
			jen.Comment("initialize our logger of choice."),
			jen.ID(constants.LoggerVarName).Assign().Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "NewZeroLogger").Call(),
			jen.Line(),
			jen.Comment("find and validate our configuration filepath."),
			jen.ID("configFilepath").Assign().Qual("os", "Getenv").Call(jen.Lit("CONFIGURATION_FILEPATH")),
			jen.If(jen.ID("configFilepath").IsEqualTo().EmptyString()).Body(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(utils.Error("no configuration file provided")),
			),
			jen.Line(),
			jen.Comment("parse our config file."),
			jen.List(jen.ID("cfg"), jen.Err()).Assign().Qual(internalConfigImp, "ParseConfigFile").Call(jen.ID("configFilepath")),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID("cfg").IsEqualTo().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Qual("fmt", "Errorf").Call(jen.Lit("error parsing configuration file: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("only allow initialization to take so long."),
			jen.List(constants.CtxVar(), jen.ID("cancel")).Assign().Qual("context", "WithTimeout").Call(jen.Qual("context", "Background").Call(), jen.ID("cfg").Dot("Meta").Dot("StartupDeadline")),
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("initialization")),
			jen.Line(),
			jen.Comment("connect to our database."),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("connecting to database")),
			jen.List(jen.ID("rawDB"), jen.Err()).Assign().ID("cfg").Dot("ProvideDatabaseConnection").Call(jen.ID(constants.LoggerVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Qual("fmt", "Errorf").Call(jen.Lit("error connecting to database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("establish the database client."),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("setting up database client")),
			jen.List(jen.ID("dbClient"), jen.Err()).Assign().ID("cfg").Dot("ProvideDatabaseClient").Call(constants.CtxVar(), jen.ID(constants.LoggerVarName), jen.ID("rawDB")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing database client: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("build our server struct."),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("building server")),
			jen.List(jen.ID("server"), jen.Err()).Assign().ID("BuildServer").Call(constants.CtxVar(), jen.ID("cfg"), jen.ID(constants.LoggerVarName), jen.ID("dbClient"), jen.ID("rawDB")),
			jen.ID(constants.SpanVarName).Dot("End").Call(),
			jen.ID("cancel").Call(),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(constants.LoggerVarName).Dot("Fatal").Call(jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing HTTP server: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("I slept and dreamt that life was joy."),
			jen.Comment("  I awoke and saw that life was service."),
			jen.Comment("  	I acted and behold, service deployed."),
			jen.ID("server").Dot("Serve").Call(),
		),
		jen.Line(),
	}
}
