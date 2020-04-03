package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metaDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("MetricAggregationMeasurement keeps track of how much time we spend collecting metrics"),
			jen.ID("MetricAggregationMeasurement").Equals().Qual("go.opencensus.io/stats", "Int64").Callln(
				jen.Lit("metrics_aggregation_time"),
				jen.Lit("cumulative time in nanoseconds spent aggregating metrics"),
				jen.Qual("go.opencensus.io/stats", "UnitDimensionless"),
			),
			jen.Line(),
			jen.Comment("MetricAggregationMeasurementView is the corresponding view for the above metric"),
			jen.ID("MetricAggregationMeasurementView").Equals().VarPointer().Qual("go.opencensus.io/stats/view", "View").Valuesln(
				jen.ID("Name").MapAssign().Lit("metrics_aggregation_time"),
				jen.ID("Measure").MapAssign().ID("MetricAggregationMeasurement"),
				jen.ID("Description").MapAssign().Lit("cumulative time in nanoseconds spent aggregating metrics"),
				jen.ID("Aggregation").MapAssign().Qual("go.opencensus.io/stats/view", "LastValue").Call(),
			),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Comment("RegisterDefaultViews registers default runtime views"),
		jen.Line(),
		jen.Func().ID("RegisterDefaultViews").Params().Params(jen.Error()).Block(
			jen.Return().Qual("go.opencensus.io/stats/view",
				"Register",
			).Call(jen.ID("DefaultRuntimeViews").Op("...")),
		),
		jen.Line(),
	)
	return ret
}
