package querier

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
			jen.ID("ErrNilInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("nil input provided")),
			jen.ID("ErrNilTransactionProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty input provided")),
			jen.ID("ErrEmptyInputProvided").Op("=").Qual("errors", "New").Call(jen.Lit("empty input provided")),
			jen.ID("ErrInvalidIDProvided").Op("=").Qual("errors", "New").Call(jen.Lit("required ID provided is zero")),
		),
		jen.Line(),
	)

	return code
}
