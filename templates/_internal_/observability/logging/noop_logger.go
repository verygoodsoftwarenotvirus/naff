package logging

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func noopLoggerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("noopLogger").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("logger").Op("=").ID("new").Call(jen.ID("noopLogger")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewNoopLogger provides our noop zerologLogger to dependency managers."),
		jen.Line(),
		jen.Func().ID("NewNoopLogger").Params().Params(jen.ID("Logger")).Body(
			jen.Return().ID("logger")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Info satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("Info").Params(jen.ID("string")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Debug satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("Debug").Params(jen.ID("string")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Error satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("Error").Params(jen.ID("error"), jen.ID("string")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Fatal satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("Fatal").Params(jen.ID("error")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Printf satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("Printf").Params(jen.ID("string"), jen.Op("...").Interface()).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetLevel satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("SetLevel").Params(jen.ID("Level")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetRequestIDFunc satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("SetRequestIDFunc").Params(jen.ID("RequestIDFunc")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithName satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("WithName").Params(jen.ID("string")).Params(jen.ID("Logger")).Body(
			jen.Return().ID("l")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Clone satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("Clone").Params().Params(jen.ID("Logger")).Body(
			jen.Return().ID("l")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithValues satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("WithValues").Params(jen.Map(jen.ID("string")).Interface()).Params(jen.ID("Logger")).Body(
			jen.Return().ID("l")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithValue satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("WithValue").Params(jen.ID("string"), jen.Interface()).Params(jen.ID("Logger")).Body(
			jen.Return().ID("l")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithRequest satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("WithRequest").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("Logger")).Body(
			jen.Return().ID("l")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithResponse satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("WithResponse").Params(jen.Op("*").Qual("net/http", "Response")).Params(jen.ID("Logger")).Body(
			jen.Return().ID("l")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WithError satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("l").Op("*").ID("noopLogger")).ID("WithError").Params(jen.ID("error")).Params(jen.ID("Logger")).Body(
			jen.Return().ID("l")),
		jen.Line(),
	)

	return code
}
