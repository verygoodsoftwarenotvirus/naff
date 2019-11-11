package frontend

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func frontendServiceDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("serviceName").Op("=").Lit("frontend_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Service is responsible for serving HTML (and other static resources)"),
			jen.ID("Service").Struct(
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("config").Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "FrontendSettings"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideFrontendService provides the frontend service to dependency injection"),
		jen.Line(),
		jen.Func().ID("ProvideFrontendService").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("cfg").Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "FrontendSettings")).Params(jen.Op("*").ID("Service")).Block(
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(
				jen.ID("config").Op(":").ID("cfg"),
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
			),
			jen.Return().ID("svc"),
		),
		jen.Line(),
	)
	return ret
}
