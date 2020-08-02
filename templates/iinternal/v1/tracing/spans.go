package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spansDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildStartSpan()...)

	return code
}

func buildStartSpan() []jen.Code {
	lines := []jen.Code{
		jen.Comment("StartSpan starts a span."),
		jen.Line(),
		jen.Func().ID("StartSpan").Params(
			constants.CtxParam(),
			jen.ID("funcName").String(),
		).Params(
			jen.Qual("context", "Context"),
			jen.PointerTo().Qual("go.opencensus.io/trace", "Span"),
		).Block(
			jen.Return(jen.Qual("go.opencensus.io/trace", "StartSpan").Call(
				constants.CtxVar(),
				jen.ID("funcName"),
			)),
		),
		jen.Line(),
	}

	return lines
}
