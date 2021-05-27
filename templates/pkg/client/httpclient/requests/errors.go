package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func errorsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ErrNoURLProvided").Op("=").Qual("errors", "New").Call(jen.Lit("no URL provided")),
			jen.ID("ErrNilEncoderProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil encoder provided")),
			jen.ID("ErrNilInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil input provided")),
			jen.ID("ErrInvalidIDProvided").Op("=").Qual("errors", "New").Call(jen.Lit("required ID provided is zero")),
			jen.ID("ErrEmptyUsernameProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty username provided")),
			jen.ID("ErrCookieRequired").Op("=").Qual("errors", "New").Call(jen.Lit("cookie required for request")),
			jen.ID("ErrInvalidPhotoEncodingForUpload").Op("=").Qual("errors", "New").Call(jen.Lit("invalid photo encoding")),
			jen.ID("ErrInvalidSecretKeyLength").Op("=").Qual("errors", "New").Call(jen.Lit("invalid secret key length")),
		),
		jen.Line(),
	)

	return code
}
