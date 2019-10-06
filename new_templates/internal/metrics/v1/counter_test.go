package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func counterTestDotGo() *jen.File {
	ret := jen.NewFile("metrics")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("Test_opencensusCounter_Increment").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("ct"), jen.ID("err")).Op(":=").ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
			jen.ID("c").Op(":=").ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"actualCount",
			), jen.ID("uint64").Call(jen.Lit(0))),
			jen.ID("c").Dot(
				"Increment",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"actualCount",
			), jen.ID("uint64").Call(jen.Lit(1))),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_opencensusCounter_IncrementBy").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("ct"), jen.ID("err")).Op(":=").ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
			jen.ID("c").Op(":=").ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"actualCount",
			), jen.ID("uint64").Call(jen.Lit(0))),
			jen.ID("c").Dot(
				"IncrementBy",
			).Call(jen.Qual("context", "Background").Call(), jen.Lit(666)),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"actualCount",
			), jen.ID("uint64").Call(jen.Lit(666))),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_opencensusCounter_Decrement").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("ct"), jen.ID("err")).Op(":=").ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
			jen.ID("c").Op(":=").ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"actualCount",
			), jen.ID("uint64").Call(jen.Lit(0))),
			jen.ID("c").Dot(
				"Increment",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"actualCount",
			), jen.ID("uint64").Call(jen.Lit(1))),
			jen.ID("c").Dot(
				"Decrement",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"actualCount",
			), jen.ID("uint64").Call(jen.Lit(0))),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestProvideUnitCounterProvider").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("assert").Dot(
			"NotNil",
		).Call(jen.ID("T"), jen.ID("ProvideUnitCounterProvider").Call()),
	),
	)
	return ret
}
