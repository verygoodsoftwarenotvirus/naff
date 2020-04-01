package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_Increment").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Assign().ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(1))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_IncrementBy").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Assign().ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
				jen.Line(),
				jen.ID("c").Dot("IncrementBy").Call(utils.CtxVar(), jen.Lit(666)),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(666))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_opencensusCounter_Decrement").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.List(jen.ID("ct"), jen.Err()).Assign().ID("ProvideUnitCounter").Call(jen.Lit("counter"), jen.Lit("description")),
				jen.ID("c").Assign().ID("ct").Assert(jen.Op("*").ID("opencensusCounter")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
				jen.Line(),
				jen.ID("c").Dot("Increment").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(1))),
				jen.Line(),
				jen.ID("c").Dot("Decrement").Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("c").Dot("actualCount"), jen.ID("uint64").Call(jen.Lit(0))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideUnitCounterProvider").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("obligatory"),
			jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("T"), jen.ID("ProvideUnitCounterProvider").Call()),
		),
		jen.Line(),
	)
	return ret
}
