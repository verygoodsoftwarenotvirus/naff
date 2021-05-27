package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestUsers").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.ID("usersTestSuite")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("usersBaseSuite").Struct(
				jen.ID("suite").Dot("Suite"),
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
				jen.ID("exampleUserList").Op("*").ID("types").Dot("UserList"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").ID("usersBaseSuite")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("usersBaseSuite")).ID("SetupTest").Params().Body(
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("s").Dot("exampleUser").Dot("HashedPassword").Op("=").Lit(""),
			jen.ID("s").Dot("exampleUser").Dot("TwoFactorSecret").Op("=").Lit(""),
			jen.ID("s").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
			jen.ID("s").Dot("exampleUserList").Op("=").ID("fakes").Dot("BuildFakeUserList").Call(),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("len").Call(jen.ID("s").Dot("exampleUserList").Dot("Users")), jen.ID("i").Op("++")).Body(
				jen.ID("s").Dot("exampleUserList").Dot("Users").Index(jen.ID("i")).Dot("HashedPassword").Op("=").Lit(""),
				jen.ID("s").Dot("exampleUserList").Dot("Users").Index(jen.ID("i")).Dot("TwoFactorSecret").Op("=").Lit(""),
				jen.ID("s").Dot("exampleUserList").Dot("Users").Index(jen.ID("i")).Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("usersTestSuite").Struct(
				jen.ID("suite").Dot("Suite"),
				jen.ID("usersBaseSuite"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("usersTestSuite")).ID("TestClient_GetUser").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/users/%d"),
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
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleUser"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleUser"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetUser").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
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
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("usersTestSuite")).ID("TestClient_GetUsers").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleUserList"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetUsers").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleUserList"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetUsers").Call(
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
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetUsers").Call(
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
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("usersTestSuite")).ID("TestClient_SearchForUsersByUsername").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users/search"),
			),
			jen.ID("exampleUsername").Op(":=").ID("s").Dot("exampleUser").Dot("Username"),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("q=%s"),
							jen.ID("exampleUsername"),
						),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleUserList").Dot("Users"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchForUsersByUsername").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleUsername"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleUserList").Dot("Users"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with empty query"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchForUsersByUsername").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(""),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchForUsersByUsername").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleUsername"),
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
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchForUsersByUsername").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleUsername"),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("usersTestSuite")).ID("TestClient_CreateUser").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildUserCreationResponseFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithRequestBodyValidation").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Op("&").ID("types").Dot("UserRegistrationInput").Values(),
						jen.ID("exampleInput"),
						jen.ID("expected"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateUser").Call(
						jen.ID("s").Dot("ctx"),
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
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateUser").Call(
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
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
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

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("usersTestSuite")).ID("TestClient_ArchiveUser").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/users/%d"),
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
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
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
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
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
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("usersTestSuite")).ID("TestClient_GetAuditLogForUser").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users/%d/audit"),
				jen.ID("expectedMethod").Op("=").Qual("net/http", "MethodGet"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("exampleAuditLogEntryList").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleAuditLogEntryList"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
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
				jen.Lit("with invalid user ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForUser").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
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
						jen.ID("s").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForUser").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("usersTestSuite")).ID("TestClient_UploadNewAvatar").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users/avatar/upload"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
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
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("exampleAvatar").Op(":=").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleExtension").Op(":=").ID("png"),
					jen.ID("err").Op(":=").ID("c").Dot("UploadNewAvatar").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleAvatar"),
						jen.ID("exampleExtension"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid avatar"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleExtension").Op(":=").ID("png"),
					jen.ID("err").Op(":=").ID("c").Dot("UploadNewAvatar").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
						jen.ID("exampleExtension"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid extension"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleAvatar").Op(":=").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("err").Op(":=").ID("c").Dot("UploadNewAvatar").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleAvatar"),
						jen.Lit(""),
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
					jen.ID("exampleAvatar").Op(":=").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleExtension").Op(":=").ID("png"),
					jen.ID("err").Op(":=").ID("c").Dot("UploadNewAvatar").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleAvatar"),
						jen.ID("exampleExtension"),
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
					jen.ID("exampleAvatar").Op(":=").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleExtension").Op(":=").ID("png"),
					jen.ID("err").Op(":=").ID("c").Dot("UploadNewAvatar").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleAvatar"),
						jen.ID("exampleExtension"),
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
