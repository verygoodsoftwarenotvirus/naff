package load

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mainDotGo(rootPkg string) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(ret)

	ret.Add(
		jen.Comment("ServiceAttacker implements hazana's Attacker interface"),
		jen.Line(),
		jen.Type().ID("ServiceAttacker").Struct(
			jen.ID("todoClient").Op("*").Qual(filepath.Join(rootPkg, "client/v1/http"), "V1Client"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Setup implement's hazana's Attacker interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("a").Op("*").ID("ServiceAttacker")).ID("Setup").Params(jen.ID("c").ID("hazana").Dot("Config")).Params(jen.ID("error")).Block(
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Do implement's hazana's Attacker interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("a").Op("*").ID("ServiceAttacker")).ID("Do").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("hazana").Dot("DoResult")).Block(
			jen.Comment("Do performs one request and is executed in a separate goroutine."),
			jen.Comment("The context is used to cancel the request on timeout."),
			jen.ID("act").Op(":=").ID("RandomAction").Call(jen.ID("a").Dot("todoClient")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("act").Dot("Action").Call(),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("req").Op("==").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").ID("ErrUnavailableYet")).Block(
					jen.Return().ID("hazana").Dot("DoResult").Valuesln(
						jen.ID("RequestLabel").Op(":").ID("act").Dot("Name"),
						jen.ID("Error").Op(":").ID("nil"), jen.ID("StatusCode").Op(":").Lit(200),
					),
				),
				jen.Qual("log", "Printf").Call(jen.Lit("something has gone awry: %v\n"), jen.ID("err")),
				jen.Return().ID("hazana").Dot("DoResult").Values(jen.ID("Error").Op(":").ID("err")),
			),
			jen.Line(),
			jen.Var().Defs(
				jen.ID("sc").ID("int"),
				jen.ID("bo").ID("int64"),
				jen.ID("bi").Index().ID("byte"),
			),
			jen.If(jen.ID("req").Dot("Body").Op("!=").ID("nil")).Block(
				jen.List(jen.ID("bi"), jen.ID("err")).Op("=").Qual("io/ioutil", "ReadAll").Call(jen.ID("req").Dot("Body")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().ID("hazana").Dot("DoResult").Values(jen.ID("Error").Op(":").ID("err")),
				),
				jen.ID("rdr").Op(":=").Qual("io/ioutil", "NopCloser").Call(jen.Qual("bytes", "NewBuffer").Call(jen.ID("bi"))),
				jen.ID("req").Dot("Body").Op("=").ID("rdr"),
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("a").Dot("todoClient").Dot("AuthenticatedClient").Call().Dot("Do").Call(jen.ID("req")),
			jen.If(jen.ID("res").Op("!=").ID("nil")).Block(
				jen.ID("sc").Op("=").ID("res").Dot("StatusCode"),
				jen.ID("bo").Op("=").ID("res").Dot("ContentLength"),
			),
			jen.Line(),
			jen.ID("dr").Op(":=").ID("hazana").Dot("DoResult").Valuesln(
				jen.ID("RequestLabel").Op(":").ID("act").Dot("Name"),
				jen.ID("Error").Op(":").ID("err"),
				jen.ID("StatusCode").Op(":").ID("sc"),
				jen.ID("BytesIn").Op(":").ID("int64").Call(jen.ID("len").Call(jen.ID("bi"))),
				jen.ID("BytesOut").Op(":").ID("bo")),
			jen.Return().ID("dr"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Teardown implements hazana's Attacker interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("a").Op("*").ID("ServiceAttacker")).ID("Teardown").Params().Params(jen.ID("error")).Block(
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Clone implements hazana's Attacker interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("a").Op("*").ID("ServiceAttacker")).ID("Clone").Params().Params(jen.ID("hazana").Dot("Attack")).Block(
			jen.Return().ID("a"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("main").Params().Block(
			jen.ID("todoClient").Op(":=").ID("initializeClient").Call(jen.ID("oa2Client")),
			jen.Line(),
			jen.Var().ID("runTime").Op("=").Lit(10).Op("*").Qual("time", "Minute"),
			jen.If(jen.ID("rt").Op(":=").Qual("os", "Getenv").Call(jen.Lit("LOADTEST_RUN_TIME")), jen.ID("rt").Op("!=").Lit("")).Block(
				jen.List(jen.ID("_rt"), jen.ID("err")).Op(":=").Qual("time", "ParseDuration").Call(jen.ID("rt")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("panic").Call(jen.ID("err")),
				),
				jen.ID("runTime").Op("=").ID("_rt"),
			),
			jen.Line(),
			jen.ID("attacker").Op(":=").Op("&").ID("ServiceAttacker").Values(jen.ID("todoClient").Op(":").ID("todoClient")),
			jen.ID("cfg").Op(":=").ID("hazana").Dot("Config").Valuesln(
				jen.ID("RPS").Op(":").Lit(50),
				jen.ID("AttackTimeSec").Op(":").ID("int").Call(jen.ID("runTime").Dot("Seconds").Call()),
				jen.ID("RampupTimeSec").Op(":").Lit(5),
				jen.ID("MaxAttackers").Op(":").Lit(50),
				jen.ID("Verbose").Op(":").ID("true"),
				jen.ID("DoTimeoutSec").Op(":").Lit(10),
			),
			jen.Line(),
			jen.ID("r").Op(":=").ID("hazana").Dot("Run").Call(jen.ID("attacker"), jen.ID("cfg")),
			jen.Line(),
			jen.Comment("inspect the report and compute whether the test has failed"),
			jen.Comment("e.g by looking at the success percentage and mean response time of each metric."),
			jen.ID("r").Dot("Failed").Op("=").ID("false"),
			jen.Line(),
			jen.ID("hazana").Dot("PrintReport").Call(jen.ID("r")),
		),
		jen.Line(),
	)
	return ret
}
