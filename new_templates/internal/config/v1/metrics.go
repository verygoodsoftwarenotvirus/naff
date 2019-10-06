package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func metricsDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("MetricsNamespace").Op("=").Lit("todo_server").Var().ID("MinimumRuntimeCollectionInterval").Op("=").Qual("time", "Second"))
	ret.Add(jen.Null().Type().ID("metricsProvider").ID("string").Type().ID("tracingProvider").ID("string"))
	ret.Add(jen.Null().Var().ID("ErrInvalidMetricsProvider").Op("=").ID("errors").Dot(
		"New",
	).Call(jen.Lit("invalid metrics provider")).Var().ID("Prometheus").ID("metricsProvider").Op("=").Lit("prometheus").Var().ID("DefaultMetricsProvider").Op("=").ID("Prometheus").Var().ID("ErrInvalidTracingProvider").Op("=").ID("errors").Dot(
		"New",
	).Call(jen.Lit("invalid tracing provider")).Var().ID("Jaeger").ID("tracingProvider").Op("=").Lit("jaeger").Var().ID("DefaultTracingProvider").Op("=").ID("Jaeger"),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
