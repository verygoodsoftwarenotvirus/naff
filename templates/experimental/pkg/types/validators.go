package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func validatorsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidType").Op("=").Qual("errors", "New").Call(jen.Lit("unexpected type received")),
			jen.ID("errDurationTooLong").Op("=").Qual("errors", "New").Call(jen.Lit("duration too long")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "Rule").Op("=").Parens(jen.Op("*").ID("urlValidator")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("urlValidator").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("v").Op("*").ID("stringDurationValidator")).ID("Validate").Params(jen.ID("value").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("raw"), jen.ID("ok")).Op(":=").ID("value").Assert(jen.ID("string")),
			jen.If(jen.Op("!").ID("ok")).Body(
				jen.Return().ID("errInvalidType")),
			jen.List(jen.ID("d"), jen.ID("err")).Op(":=").Qual("time", "ParseDuration").Call(jen.ID("raw")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("err")),
			jen.If(jen.ID("d").Op(">").ID("v").Dot("maxDuration")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("%w: %v"),
					jen.ID("errDurationTooLong"),
					jen.ID("d"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "Rule").Op("=").Parens(jen.Op("*").ID("stringDurationValidator")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("stringDurationValidator").Struct(jen.ID("maxDuration").Qual("time", "Duration")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("v").Op("*").ID("stringDurationValidator")).ID("Validate").Params(jen.ID("value").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("raw"), jen.ID("ok")).Op(":=").ID("value").Assert(jen.ID("string")),
			jen.If(jen.Op("!").ID("ok")).Body(
				jen.Return().ID("errInvalidType")),
			jen.List(jen.ID("d"), jen.ID("err")).Op(":=").Qual("time", "ParseDuration").Call(jen.ID("raw")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("err")),
			jen.If(jen.ID("d").Op(">").ID("v").Dot("maxDuration")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("%w: %v"),
					jen.ID("errDurationTooLong"),
					jen.ID("d"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
