package v1

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("server")

	utils.AddImports(proj, code)

	code.Add(buildTypeDefs(proj)...)
	code.Add(buildVarDefs()...)
	code.Add(buildProvideServer(proj)...)
	code.Add(buildServe()...)

	return code
}

func buildTypeDefs(proj *models.Project) []jen.Code {
	httpPackage := filepath.Join(proj.OutputPath, "server", "v1", "http")

	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("Server is the structure responsible for hosting all available protocols."),
			jen.Comment("In the events we adopted a gRPC implementation of the surface, this is"),
			jen.Comment("the structure that would contain it and be responsible for calling its"),
			jen.Comment("serve method."),
			jen.ID("Server").Struct(
				jen.ID("config").PointerTo().Qual(proj.InternalConfigV1Package(), "ServerConfig"),
				jen.ID("httpServer").PointerTo().Qual(httpPackage, "Server"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildVarDefs() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(jen.ID("ProvideServer")),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideServer(proj *models.Project) []jen.Code {
	httpPackage := filepath.Join(proj.OutputPath, "server", "v1", "http")

	lines := []jen.Code{
		jen.Comment("ProvideServer builds a new Server instance."),
		jen.Line(),
		jen.Func().ID("ProvideServer").Params(jen.ID("cfg").PointerTo().Qual(proj.InternalConfigV1Package(), "ServerConfig"), jen.ID("httpServer").PointerTo().Qual(httpPackage, "Server")).Params(jen.PointerTo().ID("Server"), jen.Error()).Body(
			jen.ID("srv").Assign().AddressOf().ID("Server").Valuesln(
				jen.ID("config").MapAssign().ID("cfg"),
				jen.ID("httpServer").MapAssign().ID("httpServer"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("srv"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildServe() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Serve serves HTTP traffic."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Server")).ID("Serve").Params().Body(
			jen.ID("s").Dot("httpServer").Dot("Serve").Call(),
		),
		jen.Line(),
	}

	return lines
}
