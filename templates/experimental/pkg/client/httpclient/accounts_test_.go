package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAccounts").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.ID("accountsTestSuite")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("accountsTestSuite").Struct(
				jen.ID("suite").Dot("Suite"),
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("exampleAccount").Op("*").ID("types").Dot("Account"),
				jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
				jen.ID("exampleAccountList").Op("*").ID("types").Dot("AccountList"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").ID("accountsTestSuite")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("SetupTest").Params().Body(
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("s").Dot("exampleAccount").Op("=").ID("fakes").Dot("BuildFakeAccount").Call(),
			jen.ID("s").Dot("exampleAccount").Dot("BelongsToUser").Op("=").ID("s").Dot("exampleUser").Dot("ID"),
			jen.ID("s").Dot("exampleAccountList").Op("=").ID("fakes").Dot("BuildFakeAccountList").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_SwitchActiveAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/account/select"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("s").Dot("exampleAccount").Dot("BelongsToUser").Op("=").Lit(0),
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
					jen.ID("c").Dot("authMethod").Op("=").ID("cookieAuthMethod"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("SwitchActiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("authMethod").Op("=").ID("cookieAuthMethod"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("SwitchActiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.Lit(0),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("c").Dot("authMethod").Op("=").ID("cookieAuthMethod"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("SwitchActiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("c").Dot("authMethod").Op("=").ID("cookieAuthMethod"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("SwitchActiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_GetAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleAccount"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleAccount"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_GetAccounts").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/accounts"),
			),
			jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
				jen.ID("true"),
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
				jen.ID("expectedPath"),
			),
			jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleAccountList"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("filter"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleAccountList"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_CreateAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/accounts"),
			),
			jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
				jen.ID("false"),
				jen.Qual("net/http", "MethodPost"),
				jen.Lit(""),
				jen.ID("expectedPath"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("s").Dot("exampleAccount").Dot("BelongsToUser").Op("=").Lit(0),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("s").Dot("exampleAccount")),
					jen.ID("c").Op(":=").ID("buildTestClientWithRequestBodyValidation").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleInput"),
						jen.ID("exampleInput"),
						jen.ID("s").Dot("exampleAccount"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleAccount"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AccountCreationInput").Valuesln(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("s").Dot("exampleAccount")),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("s").Dot("exampleAccount").Dot("BelongsToUser").Op("=").Lit(0),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("s").Dot("exampleAccount")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_UpdateAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPut"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleAccount"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount"),
						),
						jen.Lit("no error should be returned"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("nil"),
						),
						jen.Lit("error should be returned"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("UpdateAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleAccount"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount"),
						),
						jen.Lit("error should be returned"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_ArchiveAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodDelete"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
						jen.Lit("no error should be returned"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.Lit(0),
						),
						jen.Lit("no error should be returned"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
						jen.Lit("error should be returned"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
						jen.Lit("no error should be returned"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_AddUserToAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d/member"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAddUserToAccountInput").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("exampleInput").Dot("AccountID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
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
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
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
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAddUserToAccountInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAddUserToAccountInput").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_MarkAsDefault").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d/default"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAsDefault").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAsDefault").Call(
							jen.ID("s").Dot("ctx"),
							jen.Lit(0),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAsDefault").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAsDefault").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_RemoveUserFromAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d/members/%d"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("query").Op(":=").Qual("net/url", "Values").Valuesln(jen.ID("keys").Dot("ReasonKey").Op(":").Index().ID("string").Valuesln(jen.ID("t").Dot("Name").Call())).Dot("Encode").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodDelete"),
						jen.ID("query"),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.Lit(0),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.Lit(0),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid reason"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.Lit(""),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_ModifyMemberPermissions").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d/members/%d/permissions"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPatch"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyMemberPermissions").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyMemberPermissions").Call(
							jen.ID("s").Dot("ctx"),
							jen.Lit(0),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyMemberPermissions").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.Lit(0),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyMemberPermissions").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("ModifyUserPermissionsInput").Valuesln(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyMemberPermissions").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyMemberPermissions").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyMemberPermissions").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("s").Dot("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_TransferAccountOwnership").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/accounts/%d/transfer"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("s").Dot("ctx"),
							jen.Lit(0),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AccountOwnershipTransferInput").Valuesln(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("s").Dot("exampleAccount").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("accountsTestSuite")).ID("TestClient_GetAuditLogForAccount").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/accounts/%d/audit"),
				jen.ID("expectedMethod").Op("=").Qual("net/http", "MethodGet"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleAuditLogEntryList").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleAuditLogEntryList"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAuditLogEntryList"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
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
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForAccount").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
