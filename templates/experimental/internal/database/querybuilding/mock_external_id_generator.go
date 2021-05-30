package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockExternalIDGeneratorDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockExternalIDGenerator").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewExternalID implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockExternalIDGenerator")).ID("NewExternalID").Params().Params(jen.ID("string")).Body(
			jen.Return().ID("m").Dot("Called").Call().Dot("String").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
