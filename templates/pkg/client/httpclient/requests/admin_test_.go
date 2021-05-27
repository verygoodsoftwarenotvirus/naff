package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuilder_BuildUserReputationUpdateInputRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/admin/users/status"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildUserReputationUpdateInputRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildUserReputationUpdateInputRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildUserReputationUpdateInputRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").ID("types").Dot("UserReputationUpdateInput").Values(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
