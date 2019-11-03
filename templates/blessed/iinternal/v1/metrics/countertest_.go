package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterTestDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_Increment").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("ct"), jen.ID("err")).Op(":=").ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Op(":=").ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
				jen.Line(),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(1))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_IncrementBy").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("ct"), jen.ID("err")).Op(":=").ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Op(":=").ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
				jen.Line(),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
				jen.Line(),
				jen.ID("c").Dot("IncrementBy").Call(jen.Qual("context", "Background").Call(), jen.Lit(666)),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(666))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_Decrement").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("ct"), jen.ID("err")).Op(":=").ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Op(":=").ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
				jen.Line(),
				jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(1))),
				jen.Line(),
				jen.ID("c").Dot("Decrement").Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideUnitCounterProvider").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("obligatory"),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("T"), jen.ID("ProvideUnitCounterProvider").Call()),
		),
		jen.Line(),
	)
	return ret
}
