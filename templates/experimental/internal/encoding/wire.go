package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(
			jen.ID("ProvideServerEncoderDecoder"),
			jen.ID("ProvideContentType"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideContentType provides a ContentType from a Config."),
		jen.Line(),
		jen.Func().ID("ProvideContentType").Params(jen.ID("cfg").ID("Config")).Params(jen.ID("ContentType")).Body(
			jen.Return().ID("contentTypeFromString").Call(jen.ID("cfg").Dot("ContentType"))),
		jen.Line(),
	)

	return code
}
