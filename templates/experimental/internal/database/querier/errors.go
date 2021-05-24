package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func errorsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("ErrNilInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil input provided")).Var().ID("ErrNilTransactionProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty input provided")).Var().ID("ErrEmptyInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty input provided")).Var().ID("ErrInvalidIDProvided").Op("=").Qual("errors", "New").Call(jen.Lit("required ID provided is zero")),
		jen.Line(),
	)

	return code
}
