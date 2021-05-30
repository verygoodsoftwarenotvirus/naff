package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func utilsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestEnsureQueryIsRestrictedToUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("leaves good queries alone"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUserID").Op(":=").ID("fakes").Dot("BuildFakeUser").Call().Dot("ID"),
					jen.ID("exampleQuery").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("things +belongsToAccount:%d"),
						jen.ID("exampleUserID"),
					),
					jen.ID("expectation").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("things +belongsToAccount:%d"),
						jen.ID("exampleUserID"),
					),
					jen.ID("actual").Op(":=").ID("ensureQueryIsRestrictedToUser").Call(
						jen.ID("exampleQuery"),
						jen.ID("exampleUserID"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectation"),
						jen.ID("actual"),
						jen.Lit("expected %q to equal %q"),
						jen.ID("expectation"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("basic replacement"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUserID").Op(":=").ID("fakes").Dot("BuildFakeUser").Call().Dot("ID"),
					jen.ID("exampleQuery").Op(":=").Lit("things"),
					jen.ID("expectation").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("things +belongsToAccount:%d"),
						jen.ID("exampleUserID"),
					),
					jen.ID("actual").Op(":=").ID("ensureQueryIsRestrictedToUser").Call(
						jen.ID("exampleQuery"),
						jen.ID("exampleUserID"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectation"),
						jen.ID("actual"),
						jen.Lit("expected %q to equal %q"),
						jen.ID("expectation"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user restriction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUserID").Op(":=").ID("fakes").Dot("BuildFakeUser").Call().Dot("ID"),
					jen.ID("exampleQuery").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("stuff belongsToAccount:%d"),
						jen.ID("exampleUserID"),
					),
					jen.ID("expectation").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("stuff +belongsToAccount:%d"),
						jen.ID("exampleUserID"),
					),
					jen.ID("actual").Op(":=").ID("ensureQueryIsRestrictedToUser").Call(
						jen.ID("exampleQuery"),
						jen.ID("exampleUserID"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectation"),
						jen.ID("actual"),
						jen.Lit("expected %q to equal %q"),
						jen.ID("expectation"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
