package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_Increment").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Assign().ID("ct").Assert(jen.PointerTo().ID("opencensusCounter")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0)), nil),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(utils.CtxVar()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(1)), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_IncrementBy").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Assign().ID("ct").Assert(jen.PointerTo().ID("opencensusCounter")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0)), nil),
				jen.Line(),
				jen.ID("c").Dot("IncrementBy").Call(utils.CtxVar(), jen.Lit(666)),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(666)), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_Decrement").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Assign().ID("ct").Assert(jen.PointerTo().ID("opencensusCounter")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0)), nil),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(utils.CtxVar()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(1)), nil),
				jen.Line(),
				jen.ID("c").Dot("Decrement").Call(utils.CtxVar()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0)), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideUnitCounterProvider").Params(jen.ID("t").ParamPointer().Qual("testing", "t")).Block(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("obligatory"),
			utils.AssertNotNil(jen.ID("ProvideUnitCounterProvider").Call(), nil),
		),
		jen.Line(),
	)
	return ret
}
