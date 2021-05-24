package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func errorsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("ErrNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("404: not found")).Var().ID("ErrInvalidRequestInput").Op("=").Qual("errors", "New").Call(jen.Lit("400: bad request")).Var().ID("ErrBanned").Op("=").Qual("errors", "New").Call(jen.Lit("403: banned")).Var().ID("ErrInternalServerError").Op("=").Qual("errors", "New").Call(jen.Lit("500: internal server error")).Var().ID("ErrUnauthorized").Op("=").Qual("errors", "New").Call(jen.Lit("401: not authorized")).Var().ID("ErrNoURLProvided").Op("=").Qual("errors", "New").Call(jen.Lit("no URL provided")).Var().ID("ErrInvalidTOTPToken").Op("=").Qual("errors", "New").Call(jen.Lit("invalid TOTP token")).Var().ID("ErrNilInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil input provided")).Var().ID("ErrEmptyInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty input provided")).Var().ID("ErrInvalidIDProvided").Op("=").Qual("errors", "New").Call(jen.Lit("required ID provided is zero")).Var().ID("ErrEmptyQueryProvided").Op("=").Qual("errors", "New").Call(jen.Lit("query provided was empty")).Var().ID("ErrEmptyUsernameProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty username provided")).Var().ID("ErrCookieRequired").Op("=").Qual("errors", "New").Call(jen.Lit("cookie required for request")).Var().ID("ErrNoCookiesReturned").Op("=").Qual("errors", "New").Call(jen.Lit("no cookies returned from request")).Var().ID("ErrInvalidAvatarSize").Op("=").Qual("errors", "New").Call(jen.Lit("invalid avatar size")).Var().ID("ErrInvalidImageExtension").Op("=").Qual("errors", "New").Call(jen.Lit("invalid image extension")).Var().ID("ErrNilResponse").Op("=").Qual("errors", "New").Call(jen.Lit("nil response")).Var().ID("ErrArgumentIsNotPointer").Op("=").Qual("errors", "New").Call(jen.Lit("value is not a pointer")),
		jen.Line(),
	)

	return code
}
