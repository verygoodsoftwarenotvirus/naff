package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func utilsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildTestEnsureQueryIsRestrictedToUser(proj)...)

	return code
}

func buildTestEnsureQueryIsRestrictedToUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestEnsureQueryIsRestrictedToUser").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("leaves good queries alone",
				jen.ID("exampleUserID").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call().Dot("ID"),
				jen.Line(),
				jen.ID("exampleQuery").Assign().Qual("fmt", "Sprintf").Call(
					jen.Lit("things +belongsToUser:%d"),
					jen.ID("exampleUserID"),
				),
				jen.ID("expectation").Assign().Qual("fmt", "Sprintf").Call(
					jen.Lit("things +belongsToUser:%d"),
					jen.ID("exampleUserID"),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("ensureQueryIsRestrictedToUser").Call(
					jen.ID("exampleQuery"),
					jen.ID("exampleUserID"),
				),
				utils.AssertEqual(
					jen.ID("expectation"),
					jen.ID("actual"),
					jen.Lit("expected %q to equal %q"),
					jen.ID("expectation"),
					jen.ID("actual"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext("basic replacement",
				jen.ID("exampleUserID").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call().Dot("ID"),
				jen.Line(),
				jen.ID("exampleQuery").Assign().Lit("things"),
				jen.ID("expectation").Assign().Qual("fmt", "Sprintf").Call(
					jen.Lit("things +belongsToUser:%d"),
					jen.ID("exampleUserID"),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("ensureQueryIsRestrictedToUser").Call(
					jen.ID("exampleQuery"),
					jen.ID("exampleUserID"),
				),
				utils.AssertEqual(
					jen.ID("expectation"),
					jen.ID("actual"),
					jen.Lit("expected %q to equal %q"),
					jen.ID("expectation"),
					jen.ID("actual"),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext("with invalid user restriction",
				jen.ID("exampleUserID").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call().Dot("ID"),
				jen.Line(),
				jen.ID("exampleQuery").Assign().Qual("fmt", "Sprintf").Call(
					jen.Lit("stuff belongsToUser:%d"),
					jen.ID("exampleUserID"),
				),
				jen.ID("expectation").Assign().Qual("fmt", "Sprintf").Call(
					jen.Lit("stuff +belongsToUser:%d"),
					jen.ID("exampleUserID"),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("ensureQueryIsRestrictedToUser").Call(
					jen.ID("exampleQuery"),
					jen.ID("exampleUserID"),
				),
				utils.AssertEqual(
					jen.ID("expectation"),
					jen.ID("actual"),
					jen.Lit("expected %q to equal %q"),
					jen.ID("expectation"),
					jen.ID("actual"),
				),
			),
		),
	}

	return lines
}
