package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func metaTestDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

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
