package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("metrics")

	const (
		typeName       = "opencensusCounter"
		pointerVarName = "c"
	)

	utils.AddImports(proj, code)

	code.Add(
		jen.Comment("opencensusCounter is a Counter that interfaces with opencensus."),
		jen.Line(),
		jen.Type().ID(typeName).Struct(
			jen.ID("name").String(),
			jen.ID("actualCount").Uint64(),
			jen.ID("measure").PointerTo().Qual("go.opencensus.io/stats", "Int64Measure"),
			jen.ID("v").PointerTo().Qual("go.opencensus.io/stats/view", "View"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID(pointerVarName).PointerTo().ID(typeName)).ID("subtractFromCount").Params(
			constants.CtxParam(),
			jen.ID("value").Uint64(),
		).Block(
			jen.Qual("sync/atomic", "AddUint64").Call(
				jen.AddressOf().ID(pointerVarName).Dot("actualCount"),
				jen.BitwiseXOR().ID("value").Plus().One(),
			),
			jen.Qual("go.opencensus.io/stats", "Record").Call(
				constants.CtxVar(),
				jen.ID(pointerVarName).Dot("measure").Dot("M").Call(
					jen.Int64().Call(jen.Minus().ID("value")),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID(pointerVarName).PointerTo().ID(typeName)).ID("addToCount").Params(
			constants.CtxParam(),
			jen.ID("value").Uint64(),
		).Block(
			jen.Qual("sync/atomic", "AddUint64").Call(
				jen.AddressOf().ID(pointerVarName).Dot("actualCount"),
				jen.ID("value"),
			),
			jen.Qual("go.opencensus.io/stats", "Record").Call(
				constants.CtxVar(),
				jen.ID(pointerVarName).Dot("measure").Dot("M").Call(jen.Int64().Call(jen.ID("value"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Decrement satisfies our Counter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID(pointerVarName).PointerTo().ID(typeName)).ID("Decrement").Params(constants.CtxParam()).Block(
			jen.ID("c").Dot("subtractFromCount").Call(constants.CtxVar(), jen.One()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Increment satisfies our Counter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID(pointerVarName).PointerTo().ID(typeName)).ID("Increment").Params(constants.CtxParam()).Block(
			jen.ID("c").Dot("addToCount").Call(constants.CtxVar(), jen.One()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IncrementBy satisfies our Counter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID(pointerVarName).PointerTo().ID(typeName)).ID("IncrementBy").Params(constants.CtxParam(), jen.ID("value").Uint64()).Block(
			jen.ID("c").Dot("addToCount").Call(constants.CtxVar(), jen.ID("value")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUnitCounter provides a new counter."),
		jen.Line(),
		jen.Func().ID("ProvideUnitCounter").Params(jen.ID("counterName").ID("CounterName"), jen.ID("description").String()).Params(jen.ID("UnitCounter"), jen.Error()).Block(
			jen.ID("name").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%s_count"), jen.String().Call(jen.ID("counterName"))),
			jen.Comment("Counts/groups the lengths of lines read in."),
			jen.ID("count").Assign().Qual("go.opencensus.io/stats", "Int64").Call(jen.ID("name"), jen.ID("description"), jen.Lit("By")),
			jen.Line(),
			jen.ID("countView").Assign().AddressOf().Qual("go.opencensus.io/stats/view", "View").Valuesln(
				jen.ID("Name").MapAssign().ID("name"),
				jen.ID("Description").MapAssign().ID("description"),
				jen.ID("Measure").MapAssign().ID("count"),
				jen.ID("Aggregation").MapAssign().Qual("go.opencensus.io/stats/view", "Count").Call(),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().Qual("go.opencensus.io/stats/view", "Register").Call(jen.ID("countView")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("failed to register views: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID(pointerVarName).Assign().AddressOf().ID(typeName).Valuesln(
				jen.ID("name").MapAssign().ID("name"),
				jen.ID("measure").MapAssign().ID("count"),
				jen.ID("v").MapAssign().ID("countView"),
			),
			jen.Line(),
			jen.Return(jen.ID(pointerVarName), jen.Nil()),
		),
		jen.Line(),
	)

	return code
}
