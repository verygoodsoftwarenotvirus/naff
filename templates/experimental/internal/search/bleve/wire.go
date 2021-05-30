package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(jen.ID("ProvideBleveIndexManagerProvider")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideBleveIndexManagerProvider is a wrapper around NewBleveIndexManager."),
		jen.Line(),
		jen.Func().ID("ProvideBleveIndexManagerProvider").Params().Params(jen.ID("search").Dot("IndexManagerProvider")).Body(
			jen.Return().ID("NewBleveIndexManager")),
		jen.Line(),
	)

	return code
}
