package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("Providers is our collection of what we provide to other services."),
		jen.Newline(),
		jen.Var().ID("Providers").Op("=").Qual(constants.DependencyInjectionPkg, "NewSet").Callln(jen.ID("ProvideService")),
		jen.Newline(),
	)

	return code
}
