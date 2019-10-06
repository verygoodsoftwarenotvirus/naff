package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mainDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Type().ID("ServiceAttacker").Struct(
		jen.ID("todoClient").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "V1Client"),
	),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func().ID("main").Params().Block(
		jen.ID("todoClient").Op(":=").ID("initializeClient").Call(jen.ID("oa2Client")),
		jen.Null().Var().ID("runTime").Op("=").Lit(10).Op("*").Qual("time", "Minute"),
		jen.If(
			jen.ID("rt").Op(":=").Qual("os", "Getenv").Call(jen.Lit("LOADTEST_RUN_TIME")),
			jen.ID("rt").Op("!=").Lit(""),
		).Block(
			jen.List(jen.ID("_rt"), jen.ID("err")).Op(":=").Qual("time", "ParseDuration").Call(jen.ID("rt")),
			jen.If(
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.ID("runTime").Op("=").ID("_rt"),
		),
		jen.ID("attacker").Op(":=").Op("&").ID("ServiceAttacker").Valuesln(jen.ID("todoClient").Op(":").ID("todoClient")),
		jen.ID("cfg").Op(":=").ID("hazana").Dot(
			"Config",
		).Valuesln(jen.ID("RPS").Op(":").Lit(50), jen.ID("AttackTimeSec").Op(":").ID("int").Call(jen.ID("runTime").Dot(
			"Seconds",
		).Call()), jen.ID("RampupTimeSec").Op(":").Lit(5), jen.ID("MaxAttackers").Op(":").Lit(50), jen.ID("Verbose").Op(":").ID("true"), jen.ID("DoTimeoutSec").Op(":").Lit(10)),
		jen.ID("r").Op(":=").ID("hazana").Dot(
			"Run",
		).Call(jen.ID("attacker"), jen.ID("cfg")),
		jen.ID("r").Dot(
			"Failed",
		).Op("=").ID("false"),
		jen.ID("hazana").Dot(
			"PrintReport",
		).Call(jen.ID("r")),
	),
	)
	return ret
}
