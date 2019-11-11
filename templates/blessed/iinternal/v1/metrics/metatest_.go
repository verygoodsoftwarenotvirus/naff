package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metaTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Func().ID("TestRegisterDefaultViews").Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Comment("obligatory"),
			jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("RegisterDefaultViews").Call()),
		),
		jen.Line(),
	)
	return ret
}
