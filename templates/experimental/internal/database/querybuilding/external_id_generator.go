package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func externalIDGeneratorDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Null().Type().ID("ExternalIDGenerator").Interface(jen.ID("NewExternalID").Params().Params(jen.ID("string"))),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("UUIDExternalIDGenerator").Struct(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewExternalID implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("g").ID("UUIDExternalIDGenerator")).ID("NewExternalID").Params().Params(jen.ID("string")).Body(
			jen.Return().ID("uuid").Dot("NewString").Call()),
		jen.Line(),
	)

	return code
}
