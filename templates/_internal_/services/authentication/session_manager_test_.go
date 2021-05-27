package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func sessionManagerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("mockSessionManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockSessionManager")).ID("Load").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("token").ID("string")).Params(jen.Qual("context", "Context"), jen.ID("error")).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("token"),
			),
			jen.Return().List(jen.ID("returnArgs").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("context", "Context")), jen.ID("returnArgs").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockSessionManager")).ID("RenewToken").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx")).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockSessionManager")).ID("Get").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("key").ID("string")).Params(jen.Interface()).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("key"),
			).Dot("Get").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockSessionManager")).ID("Put").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("key").ID("string"), jen.ID("val").Interface()).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("key"),
				jen.ID("val"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockSessionManager")).ID("Commit").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string"), jen.Qual("time", "Time"), jen.ID("error")).Body(
			jen.ID("returnArgs").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("returnArgs").Dot("String").Call(jen.Lit(0)), jen.ID("returnArgs").Dot("Get").Call(jen.Lit(1)).Assert(jen.Qual("time", "Time")), jen.ID("returnArgs").Dot("Error").Call(jen.Lit(2))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockSessionManager")).ID("Destroy").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx")).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
