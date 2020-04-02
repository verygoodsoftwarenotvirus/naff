package v1

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("server")

	httpPackage := fmt.Sprintf("%s/server/v1/http", proj.OutputPath)
	utils.AddImports(proj, ret)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Server is the structure responsible for hosting all available protocols"),
			jen.Comment("In the events we adopted a gRPC implementation of the surface, this is"),
			jen.Comment("the structure that would contain it and be responsible for calling its"),
			jen.Comment("serve method"),
			jen.ID("Server").Struct(
				jen.ID("config").Op("*").Qual(proj.InternalConfigV1Package(), "ServerConfig"),
				jen.ID("httpServer").Op("*").Qual(httpPackage, "Server"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(jen.ID("ProvideServer")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideServer builds a new Server instance"),
		jen.Line(),
		jen.Func().ID("ProvideServer").Params(jen.ID("cfg").Op("*").Qual(proj.InternalConfigV1Package(), "ServerConfig"), jen.ID("httpServer").Op("*").Qual(httpPackage, "Server")).Params(jen.Op("*").ID("Server"), jen.ID("error")).Block(
			jen.ID("srv").Assign().VarPointer().ID("Server").Valuesln(
				jen.ID("config").MapAssign().ID("cfg"),
				jen.ID("httpServer").MapAssign().ID("httpServer"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("srv"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Serve serves HTTP traffic"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("Serve").Params().Block(
			jen.ID("s").Dot("httpServer").Dot("Serve").Call(),
		),
		jen.Line(),
	)
	return ret
}
