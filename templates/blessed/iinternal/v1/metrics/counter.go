package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Comment("Counter counts things"),
		jen.Line(),
		jen.Type().ID("Counter").Interface(jen.ID("Increment").Params(), jen.ID("IncrementBy").Params(jen.ID("val").ID("uint64")), jen.ID("Decrement").Params()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("opencensusCounter is a Counter that interfaces with opencensus"),
		jen.Line(),
		jen.Type().ID("opencensusCounter").Struct(
			jen.ID("name").ID("string"),
			jen.ID("actualCount").ID("uint64"),
			jen.ID("count").ParamPointer().Qual("go.opencensus.io/stats", "Int64Measure"),
			jen.ID("counter").ParamPointer().Qual("go.opencensus.io/stats/view", "View"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Increment satisfies our Counter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("opencensusCounter")).ID("Increment").Params(utils.CtxVar().Qual("context", "Context")).Block(
			jen.Qual("sync/atomic", "AddUint64").Call(jen.VarPointer().ID("c").Dot("actualCount"), jen.Lit(1)),
			jen.Qual("go.opencensus.io/stats", "Record").Call(utils.CtxVar(), jen.ID("c").Dot("count").Dot("M").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IncrementBy satisfies our Counter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("opencensusCounter")).ID("IncrementBy").Params(utils.CtxParam(), jen.ID("val").ID("uint64")).Block(
			jen.Qual("sync/atomic", "AddUint64").Call(jen.VarPointer().ID("c").Dot(
				"actualCount",
			), jen.ID("val")),
			jen.Qual("go.opencensus.io/stats", "Record").Call(utils.CtxVar(), jen.ID("c").Dot(
				"count",
			).Dot(
				"M",
			).Call(jen.ID("int64").Call(jen.ID("val")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Decrement satisfies our Counter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("opencensusCounter")).ID("Decrement").Params(utils.CtxVar().Qual("context", "Context")).Block(
			jen.Qual("sync/atomic", "AddUint64").Call(jen.VarPointer().ID("c").Dot(
				"actualCount",
			), jen.Op("^").ID("uint64").Call(jen.Lit(0))),
			jen.Qual("go.opencensus.io/stats", "Record").Call(utils.CtxVar(), jen.ID("c").Dot(
				"count",
			).Dot(
				"M",
			).Call(jen.Lit(-1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUnitCounterProvider provides UnitCounter providers"),
		jen.Line(),
		jen.Func().ID("ProvideUnitCounterProvider").Params().Params(jen.ID("UnitCounterProvider")).Block(
			jen.Return().ID("ProvideUnitCounter"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUnitCounter provides a new counter"),
		jen.Line(),
		jen.Func().ID("ProvideUnitCounter").Params(jen.ID("counterName").ID("CounterName"), jen.ID("description").ID("string")).Params(jen.ID("UnitCounter"), jen.ID("error")).Block(
			jen.ID("name").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%s_count"), jen.ID("string").Call(jen.ID("counterName"))),
			jen.Comment("Counts/groups the lengths of lines read in."),
			jen.ID("count").Assign().Qual("go.opencensus.io/stats", "Int64").Call(jen.ID("name"), jen.Lit(""), jen.Lit("By")),
			jen.Line(),
			jen.ID("countView").Assign().VarPointer().Qual("go.opencensus.io/stats/view", "View").Valuesln(
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
			jen.Return().List(jen.VarPointer().ID("opencensusCounter").Valuesln(jen.ID("name").MapAssign().ID("name"), jen.ID("count").MapAssign().ID("count"), jen.ID("counter").MapAssign().ID("countView")), jen.Nil()),
		),
		jen.Line(),
	)
	return ret
}
