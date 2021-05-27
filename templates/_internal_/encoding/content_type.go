package encoding

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func contentTypeDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ContentTypeJSON").ID("ContentType").Op("=").ID("buildContentType").Call(jen.ID("contentTypeJSON")),
			jen.ID("ContentTypeXML").ID("ContentType").Op("=").ID("buildContentType").Call(jen.ID("contentTypeXML")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("ContentType").Op("*").ID("contentType"),
			jen.ID("contentType").Op("*").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").Op("*").ID("clientEncoder")).ID("ContentType").Params().Params(jen.ID("string")).Body(
			jen.Return().ID("contentTypeToString").Call(jen.ID("e").Dot("contentType"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildContentType").Params(jen.ID("s").ID("string")).Params(jen.Op("*").ID("contentType")).Body(
			jen.ID("ct").Op(":=").ID("contentType").Call(jen.Op("&").ID("s")),
			jen.Return().Op("&").ID("ct"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("contentTypeToString").Params(jen.ID("c").Op("*").ID("contentType")).Params(jen.ID("string")).Body(
			jen.Switch(jen.ID("c")).Body(
				jen.Case(jen.ID("ContentTypeJSON")).Body(
					jen.Return().ID("contentTypeJSON")),
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.Return().ID("contentTypeXML")),
				jen.Default().Body(
					jen.Return().Lit("")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("contentTypeFromString").Params(jen.ID("val").ID("string")).Params(jen.ID("ContentType")).Body(
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("val")))).Body(
				jen.Case(jen.ID("contentTypeJSON")).Body(
					jen.Return().ID("ContentTypeJSON")),
				jen.Case(jen.ID("contentTypeXML")).Body(
					jen.Return().ID("ContentTypeXML")),
				jen.Default().Body(
					jen.Return().ID("defaultContentType")),
			)),
		jen.Line(),
	)

	return code
}
