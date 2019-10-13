package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func runtimeDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("RuntimeTotalAllocMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("total_alloc"), jen.Lit("cumulative bytes allocated for heap objects"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeTotalAllocView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("total_alloc"), jen.ID("Measure").Op(":").ID("RuntimeTotalAllocMeasurement"), jen.ID("Description").Op(":").Lit("cumulative bytes allocated for heap objects"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"Count",
	).Call()).Var().ID("RuntimeSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("sys"), jen.Lit("total bytes of memory obtained from the OS"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("sys"), jen.ID("Measure").Op(":").ID("RuntimeSysMeasurement"), jen.ID("Description").Op(":").Lit("total bytes of memory obtained from the OS"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeLookupsMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("lookups"), jen.Lit("the number of pointer lookups performed by the runtime"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeLookupsView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("lookups"), jen.ID("Measure").Op(":").ID("RuntimeLookupsMeasurement"), jen.ID("Description").Op(":").Lit("the number of pointer lookups performed by the runtime"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeMallocsMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("mallocs"), jen.Lit("the cumulative count of heap objects allocated (the number of live objects is mallocs - frees)"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeNallocsView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("mallocs"), jen.ID("Measure").Op(":").ID("RuntimeMallocsMeasurement"), jen.ID("Description").Op(":").Lit("the cumulative count of heap objects allocated (the number of live objects is mallocs - frees)"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"Count",
	).Call()).Var().ID("RuntimeFreesMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("frees"), jen.Lit("cumulative count of heap objects freed (the number of live objects is mallocs - frees)"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeFreesView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("frees"), jen.ID("Measure").Op(":").ID("RuntimeFreesMeasurement"), jen.ID("Description").Op(":").Lit("cumulative count of heap objects freed (the number of live objects is mallocs - frees)"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"Count",
	).Call()).Var().ID("RuntimeHeapAllocMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("heap_alloc"), jen.Lit("bytes of allocated heap objects"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeHeapAllocView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("heap_alloc"), jen.ID("Measure").Op(":").ID("RuntimeHeapAllocMeasurement"), jen.ID("Description").Op(":").Lit("bytes of allocated heap objects"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeHeapSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("heap_sys"), jen.Lit("bytes of heap memory obtained from the OS"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeHeapSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("heap_sys"), jen.ID("Measure").Op(":").ID("RuntimeHeapSysMeasurement"), jen.ID("Description").Op(":").Lit("bytes of heap memory obtained from the OS"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeHeapIdleMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("heap_idle"), jen.Lit("bytes in idle (unused) spans"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeHeapIdleView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("heap_idle"), jen.ID("Measure").Op(":").ID("RuntimeHeapIdleMeasurement"), jen.ID("Description").Op(":").Lit("bytes in idle (unused) spans"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeHeapInuseMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("heap_inuse"), jen.Lit("bytes in in-use spans"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeHeapInuseView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("heap_inuse"), jen.ID("Measure").Op(":").ID("RuntimeHeapInuseMeasurement"), jen.ID("Description").Op(":").Lit("bytes in in-use spans"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeHeapReleasedMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("heap_released"), jen.Lit("bytes of physical memory returned to the OS"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeHeapReleasedView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("heap_released"), jen.ID("Measure").Op(":").ID("RuntimeHeapReleasedMeasurement"), jen.ID("Description").Op(":").Lit("bytes of physical memory returned to the OS"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeHeapObjectsMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("heap_objects"), jen.Lit("the number of allocated heap objects."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeHeapObjectsView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("heap_objects"), jen.ID("Measure").Op(":").ID("RuntimeHeapObjectsMeasurement"), jen.ID("Description").Op(":").Lit("the number of allocated heap objects."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeStackInuseMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("stack_inuse"), jen.Lit("bytes in stack spans."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeStackInuseView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("stack_inuse"), jen.ID("Measure").Op(":").ID("RuntimeStackInuseMeasurement"), jen.ID("Description").Op(":").Lit("bytes in stack spans."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeStackSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("stack_sys"), jen.Lit("bytes of stack memory obtained from the OS."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeStackSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("stack_sys"), jen.ID("Measure").Op(":").ID("RuntimeStackSysMeasurement"), jen.ID("Description").Op(":").Lit("bytes of stack memory obtained from the OS."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeMSpanInuseMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("mspan_inuse"), jen.Lit("bytes of allocated mspan structures."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimemSpanInuseView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("mspan_inuse"), jen.ID("Measure").Op(":").ID("RuntimeMSpanInuseMeasurement"), jen.ID("Description").Op(":").Lit("bytes of allocated mspan structures."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeMSpanSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("mspan_sys"), jen.Lit("bytes of memory obtained from the OS for mspan structures."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimemSpanSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("mspan_sys"), jen.ID("Measure").Op(":").ID("RuntimeMSpanSysMeasurement"), jen.ID("Description").Op(":").Lit("bytes of memory obtained from the OS for mspan structures."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeMCacheInuseMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("mcache_inuse"), jen.Lit("bytes of allocated mcache structures."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeMCacheInuseView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("mcache_inuse"), jen.ID("Measure").Op(":").ID("RuntimeMCacheInuseMeasurement"), jen.ID("Description").Op(":").Lit("bytes of allocated mcache structures."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeMCacheSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("mcache_sys"), jen.Lit("bytes of memory obtained from the OS for mcache structures."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeMCacheSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("mcache_sys"), jen.ID("Measure").Op(":").ID("RuntimeMCacheSysMeasurement"), jen.ID("Description").Op(":").Lit("bytes of memory obtained from the OS for mcache structures."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeBuckHashSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("buck_hash_sys"), jen.Lit("bytes of memory in profiling bucket hash tables."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeBuckHashSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("buck_hash_sys"), jen.ID("Measure").Op(":").ID("RuntimeBuckHashSysMeasurement"), jen.ID("Description").Op(":").Lit("bytes of memory in profiling bucket hash tables."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeGCSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("gc_sys"), jen.Lit("bytes of memory in garbage collection metadata."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeGCSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("gc_sys"), jen.ID("Measure").Op(":").ID("RuntimeGCSysMeasurement"), jen.ID("Description").Op(":").Lit("bytes of memory in garbage collection metadata."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeOtherSysMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("other_sys"), jen.Lit("bytes of memory in miscellaneous off-heap runtime allocations."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeOtherSysView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("other_sys"), jen.ID("Measure").Op(":").ID("RuntimeOtherSysMeasurement"), jen.ID("Description").Op(":").Lit("bytes of memory in miscellaneous off-heap runtime allocations."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeNextGCMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("next_gc"), jen.Lit("the target heap size of the next GC cycle."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeNextGCView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("next_gc"), jen.ID("Measure").Op(":").ID("RuntimeNextGCMeasurement"), jen.ID("Description").Op(":").Lit("the target heap size of the next GC cycle."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimePauseTotalNsMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("pause_total_ns"), jen.Lit("the cumulative nanoseconds in GC"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimePauseTotalNsView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("pause_total_ns"), jen.ID("Measure").Op(":").ID("RuntimePauseTotalNsMeasurement"), jen.ID("Description").Op(":").Lit("the cumulative nanoseconds in GC"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"Count",
	).Call()).Var().ID("RuntimePauseNsMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("pause_ns"), jen.Lit("a circular buffer of recent GC stop-the-world pause times in nanoseconds"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimePauseNsView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("pause_ns"), jen.ID("Measure").Op(":").ID("RuntimePauseNsMeasurement"), jen.ID("Description").Op(":").Lit("a circular buffer of recent GC stop-the-world pause times in nanoseconds"), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimePauseEndMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("pause_end"), jen.Lit("a circular buffer of recent GC pause end times, as nanoseconds since 1970 (the UNIX epoch)."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimePauseEndView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("pause_end"), jen.ID("Measure").Op(":").ID("RuntimePauseEndMeasurement"), jen.ID("Description").Op(":").Lit("a circular buffer of recent GC pause end times, as nanoseconds since 1970 (the UNIX epoch)."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("RuntimeNumGCMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("num_gc"), jen.Lit("the number of completed GC cycles."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeNumGCView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("num_gc"), jen.ID("Measure").Op(":").ID("RuntimeNumGCMeasurement"), jen.ID("Description").Op(":").Lit("the number of completed GC cycles."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"Count",
	).Call()).Var().ID("RuntimeNumForcedGCMeasurement").Op("=").Qual("go.opencensus.io/stats", "Int64").Call(jen.Lit("num_forced_gc"), jen.Lit("the number of GC cycles that were forced by the application calling the GC function."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeNumForcedGCView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("num_forced_gc"), jen.ID("Measure").Op(":").ID("RuntimeNumForcedGCMeasurement"), jen.ID("Description").Op(":").Lit("the number of GC cycles that were forced by the application calling the GC function."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"Count",
	).Call()).Var().ID("RuntimeGCCPUFractionMeasurement").Op("=").Qual("go.opencensus.io/stats", "Float64").Call(jen.Lit("gc_cpu_fraction"), jen.Lit("the fraction of this program's available CPU time used by the GC since the program started."), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("RuntimeGCCPUFractionView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("gc_cpu_fraction"), jen.ID("Measure").Op(":").ID("RuntimeGCCPUFractionMeasurement"), jen.ID("Description").Op(":").Lit("the fraction of this program's available CPU time used by the GC since the program started."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("CPUUsageMeasurement").Op("=").Qual("go.opencensus.io/stats", "Float64").Call(jen.Lit("cpu_usage"), jen.Lit("percent of CPU used"), jen.Qual("go.opencensus.io/stats", "UnitDimensionless")).Var().ID("CPUUsageView").Op("=").Op("&").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("Name").Op(":").Lit("cpu_usage"), jen.ID("Measure").Op(":").ID("CPUUsageMeasurement"), jen.ID("Description").Op(":").Lit("percent of CPU used."), jen.ID("Aggregation").Op(":").ID("view").Dot(
		"LastValue",
	).Call()).Var().ID("DefaultRuntimeViews").Op("=").Index().Op("*").ID("view").Dot(
		"View",
	).Valuesln(jen.ID("RuntimeTotalAllocView"), jen.ID("RuntimeSysView"), jen.ID("RuntimeLookupsView"), jen.ID("RuntimeNallocsView"), jen.ID("RuntimeFreesView"), jen.ID("RuntimeHeapAllocView"), jen.ID("RuntimeHeapSysView"), jen.ID("RuntimeHeapIdleView"), jen.ID("RuntimeHeapInuseView"), jen.ID("RuntimeHeapReleasedView"), jen.ID("RuntimeHeapObjectsView"), jen.ID("RuntimeStackInuseView"), jen.ID("RuntimeStackSysView"), jen.ID("RuntimemSpanInuseView"), jen.ID("RuntimemSpanSysView"), jen.ID("RuntimeMCacheInuseView"), jen.ID("RuntimeMCacheSysView"), jen.ID("RuntimeBuckHashSysView"), jen.ID("RuntimeGCSysView"), jen.ID("RuntimeOtherSysView"), jen.ID("RuntimeNextGCView"), jen.ID("RuntimePauseTotalNsView"), jen.ID("RuntimePauseNsView"), jen.ID("RuntimePauseEndView"), jen.ID("RuntimeNumGCView"), jen.ID("RuntimeNumForcedGCView"), jen.ID("RuntimeGCCPUFractionView"), jen.ID("CPUUsageView"), jen.ID("MetricAggregationMeasurementView"), jen.ID("ochttp").Dot(
		"ServerRequestCountView",
	), jen.ID("ochttp").Dot(
		"ServerRequestBytesView",
	), jen.ID("ochttp").Dot(
		"ServerResponseBytesView",
	), jen.ID("ochttp").Dot(
		"ServerLatencyView",
	), jen.ID("ochttp").Dot(
		"ServerRequestCountByMethod",
	), jen.ID("ochttp").Dot(
		"ServerResponseCountByStatusCode",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// RecordRuntimeStats records runtime statistics at the provided interval.").Comment("// Returns a stop function and an error").ID("RecordRuntimeStats").Params(jen.ID("interval").Qual("time", "Duration")).Params(jen.ID("stopFn").Params()).Block(
		jen.Null().Var().ID("closeOnce").Qual("sync", "Once").Var().ID("ticker").Op("=").Qual("time", "NewTicker").Call(jen.ID("interval")).Var().ID("done").Op("=").ID("make").Call(jen.Chan().Struct()),
		jen.ID("ms").Op(":=").Op("&").Qual("runtime", "MemStats").Valuesln(),
		jen.Go().Func().Params().Block(
			jen.For().Block(
				jen.Select().Block(
					jen.Case(jen.Op("<-").ID("ticker").Dot(
						"C",
					)).Block(jen.ID("startTime").Op(":=").Qual("time", "Now").Call(), jen.ID("ctx").Op(":=").Qual("context", "Background").Call(), jen.Qual("runtime", "ReadMemStats").Call(jen.ID("ms")), jen.Qual("go.opencensus.io/stats", "Record").Call(jen.ID("ctx"), jen.ID("RuntimeTotalAllocMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"TotalAlloc",
					))), jen.ID("RuntimeSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"Sys",
					))), jen.ID("RuntimeLookupsMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"Lookups",
					))), jen.ID("RuntimeMallocsMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"Mallocs",
					))), jen.ID("RuntimeFreesMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"Frees",
					))), jen.ID("RuntimeHeapAllocMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"HeapAlloc",
					))), jen.ID("RuntimeHeapSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"HeapSys",
					))), jen.ID("RuntimeHeapIdleMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"HeapIdle",
					))), jen.ID("RuntimeHeapInuseMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"HeapInuse",
					))), jen.ID("RuntimeHeapReleasedMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"HeapReleased",
					))), jen.ID("RuntimeHeapObjectsMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"HeapObjects",
					))), jen.ID("RuntimeStackInuseMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"StackInuse",
					))), jen.ID("RuntimeStackSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"StackSys",
					))), jen.ID("RuntimeMSpanInuseMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"MSpanInuse",
					))), jen.ID("RuntimeMSpanSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"MSpanSys",
					))), jen.ID("RuntimeMCacheInuseMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"MCacheInuse",
					))), jen.ID("RuntimeMCacheSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"MCacheSys",
					))), jen.ID("RuntimeBuckHashSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"BuckHashSys",
					))), jen.ID("RuntimeGCSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"GCSys",
					))), jen.ID("RuntimeOtherSysMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"OtherSys",
					))), jen.ID("RuntimeNextGCMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"NextGC",
					))), jen.ID("RuntimePauseTotalNsMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"PauseTotalNs",
					))), jen.ID("RuntimePauseNsMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"PauseNs",
					).Index(jen.Parens(jen.ID("ms").Dot(
						"NumGC",
					).Op("+").Lit(255)).Op("%").Lit(256)))), jen.ID("RuntimePauseEndMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"PauseEnd",
					).Index(jen.Parens(jen.ID("ms").Dot(
						"NumGC",
					).Op("+").Lit(255)).Op("%").Lit(256)))), jen.ID("RuntimeNumGCMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"NumGC",
					))), jen.ID("RuntimeNumForcedGCMeasurement").Dot(
						"M",
					).Call(jen.ID("int64").Call(jen.ID("ms").Dot(
						"NumForcedGC",
					))), jen.ID("RuntimeGCCPUFractionMeasurement").Dot(
						"M",
					).Call(jen.ID("ms").Dot(
						"GCCPUFraction",
					))), jen.Qual("go.opencensus.io/stats", "Record").Call(jen.ID("ctx"), jen.ID("MetricAggregationMeasurement").Dot(
						"M",
					).Call(jen.Qual("time", "Since").Call(jen.ID("startTime")).Dot(
						"Nanoseconds",
					).Call()))),
					jen.Case(jen.Op("<-").ID("done")).Block(jen.ID("ticker").Dot(
						"Stop",
					).Call(), jen.Return()),
				),
			),
		).Call(),
		jen.Return().Func().Params().Block(
			jen.ID("closeOnce").Dot(
				"Do",
			).Call(jen.Func().Params().Block(
				jen.ID("close").Call(jen.ID("done")),
			)),
		),
	),

		jen.Line(),
	)
	return ret
}
