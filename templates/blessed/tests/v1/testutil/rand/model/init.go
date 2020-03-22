package model

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func initDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(proj, ret)
	ret.Add(utils.FakeSeedFunc())

	return ret
}
