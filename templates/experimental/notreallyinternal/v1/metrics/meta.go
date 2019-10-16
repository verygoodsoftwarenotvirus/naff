package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func metaDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("MetricAggregationMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("metrics_aggregation_time"), jen.Lit("cumulative time in nanoseconds spent aggregating metrics"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("MetricAggregationMeasurementView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("metrics_aggregation_time"), jen.ID("Measure").Op(":").ID("MetricAggregationMeasurement"), jen.ID("Description").Op(":").Lit("cumulative time in nanoseconds spent aggregating metrics"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("RegisterDefaultViews registers default runtime views"),
		jen.Line(),
		jen.Func().ID("RegisterDefaultViews").Params().Params(jen.ID("error")).Block(
		jen.Return().ID("view").Dot(
			"Register",
		).Call(jen.ID("DefaultRuntimeViews").Op("...")),
	),
	jen.Line(),
	)
	return ret
}
