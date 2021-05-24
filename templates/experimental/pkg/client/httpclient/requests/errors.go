package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func errorsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("ErrNoURLProvided").Op("=").Qual("errors", "New").Call(jen.Lit("no URL provided")).Var().ID("ErrNilEncoderProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil encoder provided")).Var().ID("ErrNilInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil input provided")).Var().ID("ErrInvalidIDProvided").Op("=").Qual("errors", "New").Call(jen.Lit("required ID provided is zero")).Var().ID("ErrEmptyUsernameProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty username provided")).Var().ID("ErrCookieRequired").Op("=").Qual("errors", "New").Call(jen.Lit("cookie required for request")).Var().ID("ErrInvalidPhotoEncodingForUpload").Op("=").Qual("errors", "New").Call(jen.Lit("invalid photo encoding")).Var().ID("ErrInvalidSecretKeyLength").Op("=").Qual("errors", "New").Call(jen.Lit("invalid secret key length")),
		jen.Line(),
	)

	return code
}
