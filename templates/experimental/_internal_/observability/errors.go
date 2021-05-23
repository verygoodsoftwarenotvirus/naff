package observability

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func errorsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("PrepareError standardizes our error handling by logging, tracing, and formatting an error consistently."),
		jen.Line(),
		jen.Func().ID("PrepareError").Params(jen.ID("err").ID("error"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("span").ID("tracing").Dot("Span"), jen.ID("descriptionFmt").ID("string"), jen.ID("descriptionArgs").Op("...").Interface()).Params(jen.ID("error")).Body(
			jen.If(jen.ID("err").Op("==").ID("nil")).Body(
				jen.Return().ID("nil")),
			jen.ID("desc").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("descriptionFmt"),
				jen.ID("descriptionArgs").Op("..."),
			),
			jen.ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("Error").Call(
				jen.ID("err"),
				jen.ID("desc"),
			),
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachErrorToSpan").Call(
					jen.ID("span"),
					jen.ID("desc"),
					jen.ID("err"),
				)),
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("%s: %w"),
				jen.ID("desc"),
				jen.ID("err"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AcknowledgeError standardizes our error handling by logging and tracing consistently."),
		jen.Line(),
		jen.Func().ID("AcknowledgeError").Params(jen.ID("err").ID("error"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("span").ID("tracing").Dot("Span"), jen.ID("descriptionFmt").ID("string"), jen.ID("descriptionArgs").Op("...").Interface()).Body(
			jen.ID("desc").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("descriptionFmt"),
				jen.ID("descriptionArgs").Op("..."),
			),
			jen.ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("Error").Call(
				jen.ID("err"),
				jen.ID("desc"),
			),
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachErrorToSpan").Call(
					jen.ID("span"),
					jen.ID("desc"),
					jen.ID("err"),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NoteEvent standardizes our logging and tracing notifications."),
		jen.Line(),
		jen.Func().ID("NoteEvent").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("span").ID("tracing").Dot("Span"), jen.ID("descriptionFmt").ID("string"), jen.ID("descriptionArgs").Op("...").Interface()).Body(
			jen.ID("desc").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("descriptionFmt"),
				jen.ID("descriptionArgs").Op("..."),
			),
			jen.ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("Debug").Call(jen.ID("desc")),
			jen.ID("span").Dot("AddEvent").Call(
				jen.ID("desc"),
				jen.ID("trace").Dot("WithTimestamp").Call(jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			),
		),
		jen.Line(),
	)

	return code
}
