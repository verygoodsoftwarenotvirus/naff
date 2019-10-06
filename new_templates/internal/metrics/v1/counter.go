package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func counterDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Type().ID("Counter").Interface(jen.ID("Increment").Params(), jen.ID("IncrementBy").Params(jen.ID("val").ID("uint64")), jen.ID("Decrement").Params()))
	ret.Add(jen.Null().Type().ID("opencensusCounter").Struct(
		jen.ID("name").ID("string"),
		jen.ID("actualCount").ID("uint64"),
		jen.ID("count").Op("*").Qual("go.opencensus.io/stats", "Int64Measure"),
		jen.ID("counter").Op("*").ID("view").Dot(
			"View",
		),
	),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
