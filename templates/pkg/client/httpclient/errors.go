package httpclient

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
			jen.ID("ErrNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("404: not found")),
			jen.ID("ErrInvalidRequestInput").Op("=").Qual("errors", "New").Call(jen.Lit("400: bad request")),
			jen.ID("ErrBanned").Op("=").Qual("errors", "New").Call(jen.Lit("403: banned")),
			jen.ID("ErrInternalServerError").Op("=").Qual("errors", "New").Call(jen.Lit("500: internal server error")),
			jen.ID("ErrUnauthorized").Op("=").Qual("errors", "New").Call(jen.Lit("401: not authorized")),
			jen.ID("ErrNoURLProvided").Op("=").Qual("errors", "New").Call(jen.Lit("no URL provided")),
			jen.ID("ErrInvalidTOTPToken").Op("=").Qual("errors", "New").Call(jen.Lit("invalid TOTP token")),
			jen.ID("ErrNilInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil input provided")),
			jen.ID("ErrEmptyInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty input provided")),
			jen.ID("ErrInvalidIDProvided").Op("=").Qual("errors", "New").Call(jen.Lit("required ID provided is zero")),
			jen.ID("ErrEmptyQueryProvided").Op("=").Qual("errors", "New").Call(jen.Lit("query provided was empty")),
			jen.ID("ErrEmptyUsernameProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty username provided")),
			jen.ID("ErrCookieRequired").Op("=").Qual("errors", "New").Call(jen.Lit("cookie required for request")),
			jen.ID("ErrNoCookiesReturned").Op("=").Qual("errors", "New").Call(jen.Lit("no cookies returned from request")),
			jen.ID("ErrInvalidAvatarSize").Op("=").Qual("errors", "New").Call(jen.Lit("invalid avatar size")),
			jen.ID("ErrInvalidImageExtension").Op("=").Qual("errors", "New").Call(jen.Lit("invalid image extension")),
			jen.ID("ErrNilResponse").Op("=").Qual("errors", "New").Call(jen.Lit("nil response")),
			jen.ID("ErrArgumentIsNotPointer").Op("=").Qual("errors", "New").Call(jen.Lit("value is not a pointer")),
		),
		jen.Line(),
	)

	return code
}
