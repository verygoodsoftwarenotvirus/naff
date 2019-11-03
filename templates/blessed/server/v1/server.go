package v1

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("server")

	httpPackage := fmt.Sprintf("%s/server/v1/http", pkgRoot)
	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Server is the structure responsible for hosting all available protocols"),
			jen.Comment("In the events we adopted a gRPC implementation of the surface, this is"),
			jen.Comment("the structure that would contain it and be responsible for calling its"),
			jen.Comment("serve method"),
			jen.ID("Server").Struct(
				jen.ID("config").Op("*").Qual(filepath.Join(pkgRoot, "internal/v1/config"), "ServerConfig"),
				jen.ID("httpServer").Op("*").Qual(httpPackage, "Server"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our wire superset of providers this package offers"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(jen.ID("ProvideServer")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideServer builds a new Server instance"),
		jen.Line(),
		jen.Func().ID("ProvideServer").Params(jen.ID("cfg").Op("*").Qual(filepath.Join(pkgRoot, "internal/v1/config"), "ServerConfig"), jen.ID("httpServer").Op("*").Qual(httpPackage, "Server")).Params(jen.Op("*").ID("Server"), jen.ID("error")).Block(
			jen.ID("srv").Op(":=").Op("&").ID("Server").Valuesln(
				jen.ID("config").Op(":").ID("cfg"),
				jen.ID("httpServer").Op(":").ID("httpServer"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("srv"), jen.ID("nil")),
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
