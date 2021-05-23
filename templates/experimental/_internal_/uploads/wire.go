package uploads

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
			jen.ID("ProvideUploadManager"),
			jen.ID("wire").Dot("FieldsOf").Call(
				jen.ID("new").Call(jen.Op("*").ID("Config")),
				jen.Lit("Storage"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUploadManager transforms a *storage.Uploader into an UploadManager."),
		jen.Func().ID("ProvideUploadManager").Params(jen.ID("u").Op("*").ID("storage").Dot("Uploader")).Params(jen.ID("UploadManager")).Body(
			jen.Return().ID("u")),
		jen.Line(),
	)

	return code
}
