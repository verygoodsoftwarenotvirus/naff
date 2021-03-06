package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("main")

	utils.AddImports(proj, code)

	code.Add(buildMainServiceAttacker(proj)...)
	code.Add(buildMainSetup()...)
	code.Add(buildMainDo(proj)...)
	code.Add(buildMainTeardown()...)
	code.Add(buildMainClone()...)
	code.Add(buildMainMain(proj)...)

	return code
}

func buildMainServiceAttacker(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ServiceAttacker implements hazana's Attacker interface."),
		jen.Line(),
		jen.Type().ID("ServiceAttacker").Struct(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client"),
		),
		jen.Line(),
	}

	return lines
}

func buildMainSetup() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Setup implements hazana's Attacker interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("a").PointerTo().ID("ServiceAttacker")).ID("Setup").Params(jen.Underscore().Qual("github.com/emicklei/hazana", "Config")).Params(jen.Error()).Body(
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}

	return lines
}

func buildMainDo(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("Do implements hazana's Attacker interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("a").PointerTo().ID("ServiceAttacker")).ID("Do").Params(
			jen.Underscore().Qual("context", "Context"),
		).Params(jen.Qual("github.com/emicklei/hazana", "DoResult")).Body(
			jen.Comment("Do performs one request and is executed in a separate goroutine."),
			jen.Comment("The context is used to cancel the request on timeout."),
			jen.ID("act").Assign().ID("RandomAction").Call(jen.ID("a").Dotf("%sClient", proj.Name.UnexportedVarName())),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("act").Dot("Action").Call(),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID(constants.RequestVarName).IsEqualTo().ID("nil")).Body(
				jen.If(jen.Err().IsEqualTo().ID("ErrUnavailableYet")).Body(
					jen.Return().Qual("github.com/emicklei/hazana", "DoResult").Valuesln(
						jen.ID("RequestLabel").MapAssign().ID("act").Dot("Name"),
						jen.ID("Error").MapAssign().ID("nil"), jen.ID("StatusCode").MapAssign().Lit(200),
					),
				),
				jen.Qual("log", "Printf").Call(jen.Lit("something has gone awry: %v\n"), jen.Err()),
				jen.Return().Qual("github.com/emicklei/hazana", "DoResult").Values(jen.ID("Error").MapAssign().Err()),
			),
			jen.Line(),
			jen.Var().Defs(
				jen.ID("sc").ID("int"),
				jen.ID("bo").ID("int64"),
				jen.ID("bi").Index().Byte(),
			),
			jen.If(jen.ID(constants.RequestVarName).Dot("Body").DoesNotEqual().ID("nil")).Body(
				jen.List(jen.ID("bi"), jen.Err()).Equals().Qual("io/ioutil", "ReadAll").Call(jen.ID(constants.RequestVarName).Dot("Body")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().Qual("github.com/emicklei/hazana", "DoResult").Values(jen.ID("Error").MapAssign().Err()),
				),
				jen.ID("rdr").Assign().Qual("io/ioutil", "NopCloser").Call(jen.Qual("bytes", "NewBuffer").Call(jen.ID("bi"))),
				jen.ID(constants.RequestVarName).Dot("Body").Equals().ID("rdr"),
			),
			jen.Line(),
			jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().ID("a").Dotf("%sClient", proj.Name.UnexportedVarName()).Dot("AuthenticatedClient").Call().Dot("Do").Call(jen.ID(constants.RequestVarName)),
			jen.If(jen.ID(constants.ResponseVarName).DoesNotEqual().ID("nil")).Body(
				jen.ID("sc").Equals().ID(constants.ResponseVarName).Dot("StatusCode"),
				jen.ID("bo").Equals().ID(constants.ResponseVarName).Dot("ContentLength"),
			),
			jen.Line(),
			jen.ID("dr").Assign().Qual("github.com/emicklei/hazana", "DoResult").Valuesln(
				jen.ID("RequestLabel").MapAssign().ID("act").Dot("Name"),
				jen.ID("Error").MapAssign().Err(),
				jen.ID("StatusCode").MapAssign().ID("sc"),
				jen.ID("BytesIn").MapAssign().ID("int64").Call(jen.Len(jen.ID("bi"))),
				jen.ID("BytesOut").MapAssign().ID("bo")),
			jen.Return().ID("dr"),
		),
		jen.Line(),
	}

	return lines
}

func buildMainTeardown() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Teardown implements hazana's Attacker interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("a").PointerTo().ID("ServiceAttacker")).ID("Teardown").Params().Params(jen.Error()).Body(
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}

	return lines
}

func buildMainClone() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Clone implements hazana's Attacker interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("a").PointerTo().ID("ServiceAttacker")).ID("Clone").Params().Params(jen.Qual("github.com/emicklei/hazana", "Attack")).Body(
			jen.Return().ID("a"),
		),
		jen.Line(),
	}

	return lines
}

func buildMainMain(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("main").Params().Body(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Assign().ID("initializeClient").Call(jen.ID("oa2Client")),
			jen.Line(),
			jen.Var().ID("runTime").Equals().Lit(10).Times().Qual("time", "Minute"),
			jen.If(jen.ID("rt").Assign().Qual("os", "Getenv").Call(jen.Lit("LOADTEST_RUN_TIME")), jen.ID("rt").DoesNotEqual().EmptyString()).Body(
				jen.List(jen.ID("_rt"), jen.Err()).Assign().Qual("time", "ParseDuration").Call(jen.ID("rt")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("panic").Call(jen.Err()),
				),
				jen.ID("runTime").Equals().ID("_rt"),
			),
			jen.Line(),
			jen.ID("attacker").Assign().AddressOf().ID("ServiceAttacker").Values(jen.IDf("%sClient", proj.Name.UnexportedVarName()).MapAssign().IDf("%sClient", proj.Name.UnexportedVarName())),
			jen.ID("cfg").Assign().Qual("github.com/emicklei/hazana", "Config").Valuesln(
				jen.ID("RPS").MapAssign().Lit(50),
				jen.ID("AttackTimeSec").MapAssign().ID("int").Call(jen.ID("runTime").Dot("Seconds").Call()),
				jen.ID("RampupTimeSec").MapAssign().Lit(5),
				jen.ID("MaxAttackers").MapAssign().Lit(50),
				jen.ID("Verbose").MapAssign().True(),
				jen.ID("DoTimeoutSec").MapAssign().Lit(10),
			),
			jen.Line(),
			jen.ID("r").Assign().Qual("github.com/emicklei/hazana", "Run").Call(jen.ID("attacker"), jen.ID("cfg")),
			jen.Line(),
			jen.Comment("inspect the report and compute whether the test has failed"),
			jen.Comment("e.g by looking at the success percentage and mean response time of each metric."),
			jen.ID("r").Dot("Failed").Equals().False(),
			jen.Line(),
			jen.Qual("github.com/emicklei/hazana", "PrintReport").Call(jen.ID("r")),
		),
		jen.Line(),
	}

	return lines
}
