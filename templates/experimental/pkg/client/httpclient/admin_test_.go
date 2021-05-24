package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAdmin").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.ID("adminTestSuite")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("adminTestSuite").Struct(
			jen.ID("suite").Dot("Suite"),
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("exampleAccount").Op("*").ID("types").Dot("Account"),
			jen.ID("exampleAccountList").Op("*").ID("types").Dot("AccountList"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").ID("adminTestSuite")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("adminTestSuite")).ID("SetupTest").Params().Body(
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dot("exampleAccount").Op("=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.ID("s").Dot("exampleAccountList").Op("=").ID("fakes").Dot("BuildFakeAccountList").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("adminTestSuite")).ID("TestClient_UpdateUserReputation").Params().Body(
			jen.Var().ID("expectedPath").Op("=").Lit("/api/v1/admin/users/status"),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusAccepted"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateUserReputation").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusAccepted"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateUserReputation").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("UserReputationUpdateInput").Valuesln(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusAccepted"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateUserReputation").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with bad request response"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusBadRequest"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateUserReputation").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with otherwise invalid status code response"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusInternalServerError"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateUserReputation").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateUserReputation").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateUserReputation").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
