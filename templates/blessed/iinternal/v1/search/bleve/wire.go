package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what this library offers to external users in the form of dependencies."),
			jen.Line(),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideBleveIndexManagerProvider"),
			),
		),
		jen.Line(),
		jen.Comment("ProvideBleveIndexManagerProvider is a wrapper around NewBleveIndexManager"),
		jen.Line(),
		jen.Func().ID("ProvideBleveIndexManagerProvider").Params().Params(jen.Qual(proj.InternalSearchV1Package(), "IndexManagerProvider")).Block(
			jen.Return().ID("NewBleveIndexManager"),
		),
	)

	return code
}
