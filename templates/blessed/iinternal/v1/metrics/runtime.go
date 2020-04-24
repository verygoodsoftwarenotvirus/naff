package metrics

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func float64Metric(varName, measurementName, measurementDescription string, count bool) jen.Code {
	sn := fmt.Sprintf("Runtime%sMeasurement", varName)
	vn := fmt.Sprintf("Runtime%sView", varName)
	g := jen.Group{}

	return g.Add(
		jen.Commentf("%s captures the runtime memstats %s field", sn, varName),
		jen.Line(),
		statsDotFloat64(fmt.Sprintf("Runtime%sMeasurement", varName), measurementName, measurementDescription),
		jen.Line(),
		jen.Commentf("%s is the corresponding view for the above field", vn),
		jen.Line(),
		viewDotView(vn, measurementName, sn, measurementDescription, count),
		jen.Line(),
		jen.Line(),
	)
}

func int64Metric(varName, measurementName, measurementDescription string, count bool) jen.Code {
	sn := fmt.Sprintf("Runtime%sMeasurement", varName)
	vn := fmt.Sprintf("Runtime%sView", varName)
	g := jen.Group{}

	return g.Add(
		jen.Commentf("%s captures the runtime memstats %s field", sn, varName),
		jen.Line(),
		statsDotInt64(fmt.Sprintf("Runtime%sMeasurement", varName), measurementName, measurementDescription),
		jen.Line(),
		jen.Commentf("%s is the corresponding view for the above field", vn),
		jen.Line(),
		viewDotView(vn, measurementName, sn, measurementDescription, count),
		jen.Line(),
		jen.Line(),
	)
}

func viewDotView(viewName, measurementName, measurementVarName, description string, count bool) jen.Code {
	agg := jen.ID("Aggregation").MapAssign()
	if count {
		agg = agg.Qual("go.opencensus.io/stats/view", "Count").Call()
	} else {
		agg = agg.Qual("go.opencensus.io/stats/view", "LastValue").Call()
	}

	return jen.ID(viewName).Equals().AddressOf().Qual("go.opencensus.io/stats/view", "View").Valuesln(
		jen.ID("Name").MapAssign().Lit(measurementName),
		jen.ID("Measure").MapAssign().ID(measurementVarName),
		jen.ID("Description").MapAssign().Lit(description),
		agg,
	)
}

func statsDotFloat64(varName, name, description string) jen.Code {
	return jen.ID(varName).Equals().Qual("go.opencensus.io/stats", "Float64").Callln(
		jen.Lit(name),
		jen.Lit(description),
		jen.Qual("go.opencensus.io/stats", "UnitDimensionless"),
	)
}

func statsDotInt64(varName, name, description string) jen.Code {
	return jen.ID(varName).Equals().Qual("go.opencensus.io/stats", "Int64").Callln(
		jen.Lit(name),
		jen.Lit(description),
		jen.Qual("go.opencensus.io/stats", "UnitDimensionless"),
	)
}

func runtimeDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(proj, ret)

	ret.Comment("inspired by:")
	ret.Comment("https://github.com/opencensus-integrations/caddy/blob/c8498719b7c1c2a3c707355be2395a35f03e434e/caddy/caddymain/exporters.go#L54-L110")
	ret.Line()

	var (
		defs []jen.Code
		vals []jen.Code
	)

	for _, metric := range allMetrics {
		if metric.isFloat {
			defs = append(defs, float64Metric(metric.varName, metric.measurementName, metric.measurementDescription, metric.count))
		} else {
			defs = append(defs, int64Metric(metric.varName, metric.measurementName, metric.measurementDescription, metric.count))
		}
		vals = append(vals, jen.ID(fmt.Sprintf("Runtime%sView", metric.varName)))
	}

	defs = append(defs,
		jen.Line(),
		jen.Comment("MetricAggregationMeasurement keeps track of how much time we spend collecting metrics"),
		jen.ID("MetricAggregationMeasurement").Equals().Qual("go.opencensus.io/stats", "Int64").Callln(
			jen.Lit("metrics_aggregation_time"),
			jen.Lit("cumulative time in nanoseconds spent aggregating metrics"),
			jen.Qual("go.opencensus.io/stats", "UnitDimensionless"),
		),
		jen.Comment("MetricAggregationMeasurementView is the corresponding view for the above metric"),
		jen.ID("MetricAggregationMeasurementView").Equals().AddressOf().Qual("go.opencensus.io/stats/view", "View").Valuesln(
			jen.ID("Name").MapAssign().Lit("metrics_aggregation_time"),
			jen.ID("Measure").MapAssign().ID("MetricAggregationMeasurement"),
			jen.ID("Description").MapAssign().Lit("cumulative time in nanoseconds spent aggregating metrics"),
			jen.ID("Aggregation").MapAssign().Qual("go.opencensus.io/stats/view", "LastValue").Call(),
		),
		jen.Line(),
	)

	vals = append(vals,
		jen.ID("MetricAggregationMeasurementView"),
		jen.Comment("provided by ochttp"),
		jen.Qual("go.opencensus.io/plugin/ochttp", "ServerRequestCountView"),
		jen.Qual("go.opencensus.io/plugin/ochttp", "ServerRequestBytesView"),
		jen.Qual("go.opencensus.io/plugin/ochttp", "ServerResponseBytesView"),
		jen.Qual("go.opencensus.io/plugin/ochttp", "ServerLatencyView"),
		jen.Qual("go.opencensus.io/plugin/ochttp", "ServerRequestCountByMethod"),
		jen.Qual("go.opencensus.io/plugin/ochttp", "ServerResponseCountByStatusCode"),
	)

	defs = append(defs,
		jen.Comment("DefaultRuntimeViews represents the pre-configured views"),
		jen.ID("DefaultRuntimeViews").Equals().Index().PointerTo().Qual("go.opencensus.io/stats/view", "View").Valuesln(vals...),
	)

	ret.Add(
		jen.Var().Defs(defs...),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("RegisterDefaultViews registers default runtime views"),
		jen.Line(),
		jen.Func().ID("RegisterDefaultViews").Params().Params(jen.Error()).Block(
			jen.Return().Qual("go.opencensus.io/stats/view", "Register").Call(jen.ID("DefaultRuntimeViews").Spread()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("RecordRuntimeStats records runtime statistics at the provided interval."),
		jen.Line(),
		jen.Comment("Returns a stop function and an error"),
		jen.Line(),
		jen.Func().ID("RecordRuntimeStats").Params(jen.ID("interval").Qual("time", "Duration")).Params(jen.ID("stopFn").Func().Params()).Block(
			jen.Var().Defs(
				jen.ID("closeOnce").Qual("sync", "Once"),
				jen.ID("ticker").Equals().Qual("time", "NewTicker").Call(jen.ID("interval")),
				jen.ID("done").Equals().ID("make").Call(jen.Chan().Struct()),
			),
			jen.Line(),
			jen.Go().Func().Params().Block(
				jen.For().Block(
					jen.Select().Block(
						jen.Case(jen.Op("<-").ID("ticker").Dot("C")).Block(
							constants.CreateCtx(),
							jen.ID("startTime").Assign().Qual("time", "Now").Call(),
							jen.ID("ms").Assign().AddressOf().Qual("runtime", "MemStats").Values(),
							jen.Line(),
							jen.Qual("runtime", "ReadMemStats").Call(jen.ID("ms")),
							jen.Qual("go.opencensus.io/stats", "Record").Callln(
								constants.CtxVar(),
								jen.ID("RuntimeTotalAllocMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("TotalAlloc"))),
								jen.ID("RuntimeSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("Sys"))),
								jen.ID("RuntimeLookupsMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("Lookups"))),
								jen.ID("RuntimeMallocsMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("Mallocs"))),
								jen.ID("RuntimeFreesMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("Frees"))),
								jen.ID("RuntimeHeapAllocMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("HeapAlloc"))),
								jen.ID("RuntimeHeapSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("HeapSys"))),
								jen.ID("RuntimeHeapIdleMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("HeapIdle"))),
								jen.ID("RuntimeHeapInuseMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("HeapInuse"))),
								jen.ID("RuntimeHeapReleasedMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("HeapReleased"))),
								jen.ID("RuntimeHeapObjectsMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("HeapObjects"))),
								jen.ID("RuntimeStackInuseMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("StackInuse"))),
								jen.ID("RuntimeStackSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("StackSys"))),
								jen.ID("RuntimeMSpanInuseMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("MSpanInuse"))),
								jen.ID("RuntimeMSpanSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("MSpanSys"))),
								jen.ID("RuntimeMCacheInuseMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("MCacheInuse"))),
								jen.ID("RuntimeMCacheSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("MCacheSys"))),
								jen.ID("RuntimeBuckHashSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("BuckHashSys"))),
								jen.ID("RuntimeGCSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("GCSys"))),
								jen.ID("RuntimeOtherSysMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("OtherSys"))),
								jen.ID("RuntimeNextGCMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("NextGC"))),
								jen.ID("RuntimePauseTotalNsMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("PauseTotalNs"))),
								jen.ID("RuntimePauseNsMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("PauseNs").Index(jen.Parens(jen.ID("ms").Dot("NumGC").Plus().Lit(255)).Op("%").Lit(256)))),
								jen.ID("RuntimePauseEndMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("PauseEnd").Index(jen.Parens(jen.ID("ms").Dot("NumGC").Plus().Lit(255)).Op("%").Lit(256)))),
								jen.ID("RuntimeNumGCMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("NumGC"))),
								jen.ID("RuntimeNumForcedGCMeasurement").Dot("M").Call(jen.ID("int64").Call(jen.ID("ms").Dot("NumForcedGC"))),
								jen.ID("RuntimeGCCPUFractionMeasurement").Dot("M").Call(jen.ID("ms").Dot("GCCPUFraction")),
								jen.ID("MetricAggregationMeasurement").Dot("M").Call(jen.Qual("time", "Since").Call(jen.ID("startTime")).Dot("Nanoseconds").Call()),
							),
						),
						jen.Case(jen.Op("<-").ID("done")).Block(jen.ID("ticker").Dot("Stop").Call(), jen.Return()),
					),
				),
			).Call(),
			jen.Line(),
			jen.Return().Func().Params().Block(
				jen.ID("closeOnce").Dot(
					"Do",
				).Call(
					jen.Func().Params().Block(
						jen.ID("close").Call(
							jen.ID("done"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return ret
}

var (
	allMetrics = []struct {
		varName, measurementName, measurementDescription string
		count, isFloat                                   bool
	}{
		{
			varName:                "TotalAlloc",
			measurementName:        "total_alloc",
			measurementDescription: "cumulative bytes allocated for heap objects",
			count:                  true,
		},
		{
			varName:                "Sys",
			measurementName:        "sys",
			measurementDescription: "total bytes of memory obtained from the OS",
			count:                  false,
		},
		{
			varName:                "Lookups",
			measurementName:        "lookups",
			measurementDescription: "the number of pointer lookups performed by the runtime",
			count:                  false,
		},
		{
			varName:                "Mallocs",
			measurementName:        "mallocs",
			measurementDescription: "the cumulative count of heap objects allocated (the number of live objects is mallocs - frees)",
			count:                  true,
		},
		{
			varName:                "Frees",
			measurementName:        "frees",
			measurementDescription: "cumulative count of heap objects freed (the number of live objects is mallocs - frees)",
			count:                  true,
		},
		{
			varName:                "HeapAlloc",
			measurementName:        "heap_alloc",
			measurementDescription: "bytes of allocated heap objects",
			count:                  false,
		},
		{
			varName:                "HeapSys",
			measurementName:        "heap_sys",
			measurementDescription: "bytes of heap memory obtained from the OS",
			count:                  false,
		},
		{
			varName:                "HeapIdle",
			measurementName:        "heap_idle",
			measurementDescription: "bytes in idle (unused) spans",
			count:                  false,
		},
		{
			varName:                "HeapInuse",
			measurementName:        "heap_inuse",
			measurementDescription: "bytes in in-use spans",
			count:                  false,
		},
		{
			varName:                "HeapReleased",
			measurementName:        "heap_released",
			measurementDescription: "bytes of physical memory returned to the OS",
			count:                  false,
		},
		{
			varName:                "HeapObjects",
			measurementName:        "heap_objects",
			measurementDescription: "the number of allocated heap objects.",
			count:                  false,
		},
		{
			varName:                "StackInuse",
			measurementName:        "stack_inuse",
			measurementDescription: "bytes in stack spans.",
			count:                  false,
		},
		{
			varName:                "StackSys",
			measurementName:        "stack_sys",
			measurementDescription: "bytes of stack memory obtained from the OS.",
			count:                  false,
		},
		{
			varName:                "MSpanInuse",
			measurementName:        "mspan_inuse",
			measurementDescription: "bytes of allocated mspan structures.",
			count:                  false,
		},
		{
			varName:                "MSpanSys",
			measurementName:        "mspan_sys",
			measurementDescription: "bytes of memory obtained from the OS for mspan structures.",
			count:                  false,
		},
		{
			varName:                "MCacheInuse",
			measurementName:        "mcache_inuse",
			measurementDescription: "bytes of allocated mcache structures.",
			count:                  false,
		},
		{
			varName:                "MCacheSys",
			measurementName:        "mcache_sys",
			measurementDescription: "bytes of memory obtained from the OS for mcache structures.",
			count:                  false,
		},
		{
			varName:                "BuckHashSys",
			measurementName:        "buck_hash_sys",
			measurementDescription: "bytes of memory in profiling bucket hash tables.",
			count:                  false,
		},
		{
			varName:                "GCSys",
			measurementName:        "gc_sys",
			measurementDescription: "bytes of memory in garbage collection metadata.",
			count:                  false,
		},
		{
			varName:                "OtherSys",
			measurementName:        "other_sys",
			measurementDescription: "bytes of memory in miscellaneous off-heap runtime allocations.",
			count:                  false,
		},
		{
			varName:                "NextGC",
			measurementName:        "next_gc",
			measurementDescription: "the target heap size of the next GC cycle.",
			count:                  false,
		},
		{
			varName:                "PauseTotalNs",
			measurementName:        "pause_total_ns",
			measurementDescription: "the cumulative nanoseconds in GC",
			count:                  true,
		},
		{
			varName:                "PauseNs",
			measurementName:        "pause_ns",
			measurementDescription: "a circular buffer of recent GC stop-the-world pause times in nanoseconds",
			count:                  false,
		},
		{
			varName:                "PauseEnd",
			measurementName:        "pause_end",
			measurementDescription: "a circular buffer of recent GC pause end times, as nanoseconds since 1970 (the UNIX epoch).",
			count:                  false,
		},
		{
			varName:                "NumGC",
			measurementName:        "num_gc",
			measurementDescription: "the number of completed GC cycles.",
			count:                  true,
		},
		{
			varName:                "NumForcedGC",
			measurementName:        "num_forced_gc",
			measurementDescription: "the number of GC cycles that were forced by the application calling the GC function.",
			count:                  true,
		},
		{
			varName:                "GCCPUFraction",
			measurementName:        "gc_cpu_fraction",
			measurementDescription: "the fraction of this program's available CPU time used by the GC since the program started.",
			count:                  false,
			isFloat:                true,
		},
	}
)
