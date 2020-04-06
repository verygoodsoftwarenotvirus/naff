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
				jen.Line(),
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("v"), jen.Lit("description")),
				jen.List(jen.ID("c"), jen.ID("typOK")).Assign().ID("ct").Assert(jen.PointerTo().ID("opencensusCounter")),
				utils.RequireNotNil(jen.ID("c"), nil),
				utils.RequireTrue(jen.ID("typOK"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.Uint64().Call(jen.Zero()), nil),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(utils.CtxVar()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.Uint64().Call(jen.One()), nil),
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
				jen.Line(),
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("v"), jen.Lit("description")),
				jen.List(jen.ID("c"), jen.ID("typOK")).Assign().ID("ct").Assert(jen.PointerTo().ID("opencensusCounter")),
				utils.RequireNotNil(jen.ID("c"), nil),
				utils.RequireTrue(jen.ID("typOK"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.Uint64().Call(jen.Zero()), nil),
				jen.Line(),
				jen.ID("c").Dot("IncrementBy").Call(utils.CtxVar(), jen.Lit(666)),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.Uint64().Call(jen.Lit(666)), nil),
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
				jen.Line(),
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("v"), jen.Lit("description")),
				jen.List(jen.ID("c"), jen.ID("typOK")).Assign().ID("ct").Assert(jen.PointerTo().ID("opencensusCounter")),
				utils.RequireNotNil(jen.ID("c"), nil),
				utils.RequireTrue(jen.ID("typOK"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.Uint64().Call(jen.Zero()), nil),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(utils.CtxVar()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.Uint64().Call(jen.One()), nil),
				jen.Line(),
				jen.ID("c").Dot("Decrement").Call(utils.CtxVar()),
				utils.AssertEqual(jen.ID("c").Dot("actualCount"), jen.Uint64().Call(jen.Zero()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideUnitCounterProvider").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("obligatory"),
			utils.AssertNotNil(jen.ID("ProvideUnitCounterProvider").Call(), nil),
		),
		jen.Line(),
	)

	return ret
}
