package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what we provide to dependency injectors"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID("ProvideDatabaseClient"),
			),
		),
		jen.Line(),
	)
	return ret
}
