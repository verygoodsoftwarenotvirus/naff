package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func counterDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

	ret.Add(
		jen.Type().ID("Counter").Interface(jen.ID("Increment").Params(), jen.ID("IncrementBy").Params(jen.ID("val").ID("uint64")), jen.ID("Decrement").Params()),
	jen.Line(),
	)

	ret.Add(
		jen.Type().ID("opencensusCounter").Struct(jen.ID("name").ID("string"), jen.ID("actualCount").ID("uint64"), jen.ID("count").Op("*").Qual("go.opencensus.io/stats", "Int64Measure"), jen.ID("counter").Op("*").ID("view").Dot(
		"View",
	)),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("Increment satisfies our Counter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("opencensusCounter")).ID("Increment").Params(jen.ID("ctx").Qual("context", "Context")).Block(
		jen.Qual("sync/atomic", "AddUint64").Call(jen.Op("&").ID("c").Dot(
			"actualCount",
		), jen.Lit(1)),
		jen.Qual("go.opencensus.io/stats", "Record").Call(jen.ID("ctx"), jen.ID("c").Dot(
			"count",
		).Dot(
			"M",
		).Call(jen.Lit(1))),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("IncrementBy satisfies our Counter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("opencensusCounter")).ID("IncrementBy").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("val").ID("uint64")).Block(
		jen.Qual("sync/atomic", "AddUint64").Call(jen.Op("&").ID("c").Dot(
			"actualCount",
		), jen.ID("val")),
		jen.Qual("go.opencensus.io/stats", "Record").Call(jen.ID("ctx"), jen.ID("c").Dot(
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
		jen.Func().Params(jen.ID("c").Op("*").ID("opencensusCounter")).ID("Decrement").Params(jen.ID("ctx").Qual("context", "Context")).Block(
		jen.Qual("sync/atomic", "AddUint64").Call(jen.Op("&").ID("c").Dot(
			"actualCount",
		), jen.Op("^").ID("uint64").Call(jen.Lit(0))),
		jen.Qual("go.opencensus.io/stats", "Record").Call(jen.ID("ctx"), jen.ID("c").Dot(
			"count",
		).Dot(
			"M",
		).Call(jen.Op("-").Lit(1))),
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
		jen.ID("name").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s_count"), jen.ID("string").Call(jen.ID("counterName"))),
		jen.ID("count").Op(":=").Qual("go.opencensus.io/stats", "Int64").Call(jen.ID("name"), jen.Lit(""), jen.Lit("By")),
		jen.ID("countView").Op(":=").Op("&").ID("view").Dot(
			"View",
		).Valuesln(jen.ID("Name").Op(":").ID("name"), jen.ID("Description").Op(":").ID("description"), jen.ID("Measure").Op(":").ID("count"), jen.ID("Aggregation").Op(":").ID("view").Dot(
			"Count",
		).Call()),
		jen.If(jen.ID("err").Op(":=").ID("view").Dot(
			"Register",
		).Call(jen.ID("countView")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("failed to register views: %w"), jen.ID("err"))),
		),
		jen.Return().List(jen.Op("&").ID("opencensusCounter").Valuesln(jen.ID("name").Op(":").ID("name"), jen.ID("count").Op(":").ID("count"), jen.ID("counter").Op(":").ID("countView")), jen.ID("nil")),
	),
	jen.Line(),
	)
	return ret
}
