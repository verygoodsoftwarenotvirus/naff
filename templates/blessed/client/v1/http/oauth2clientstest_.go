package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(buildV1Client_BuildGetOAuth2ClientRequest()...)
	ret.Add(buildV1Client_GetOAuth2Client(proj)...)
	ret.Add(buildV1Client_BuildGetOAuth2ClientsRequest()...)
	ret.Add(buildV1Client_GetOAuth2Clients(proj)...)
	ret.Add(buildV1Client_BuildCreateOAuth2ClientRequest(proj)...)
	ret.Add(buildV1Client_CreateOAuth2Client(proj)...)
	ret.Add(buildV1Client_BuildArchiveOAuth2ClientRequest()...)
	ret.Add(buildV1Client_ArchiveOAuth2Client()...)

	return ret
}

func buildV1Client_BuildGetOAuth2ClientRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildGetOAuth2ClientRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("expectedID").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
					utils.CtxVar(),
					jen.ID("expectedID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				utils.AssertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_GetOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_GetOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(proj.OutputPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("ClientID").MapAssign().Add(utils.FakeStringFunc()),
					jen.ID("ClientSecret").MapAssign().Lit("blah"),
				),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertTrue(
								jen.Qual("strings", "HasSuffix").Call(
									jen.ID("req").Dot("URL").Dot("String").Call(),
									jen.Qual("strconv", "Itoa").Call(
										jen.ID("int").Call(
											jen.ID("expected").Dot("ID"),
										),
									),
								),
								nil,
							),
							utils.AssertEqual(
								jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("/api/v1/oauth2/clients/%d"),
									jen.ID("expected").Dot("ID"),
								),
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("expected and actual path don't match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
								nil,
							),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("GetOAuth2Client").Call(
					utils.CtxVar(),
					jen.ID("expected").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_BuildGetOAuth2ClientsRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildGetOAuth2ClientsRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_GetOAuth2Clients(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_GetOAuth2Clients").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().VarPointer().Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientList").Valuesln(
					jen.ID("Clients").MapAssign().Index().Qual(filepath.Join(outPath, "models/v1"), "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
							jen.ID("ClientID").MapAssign().Add(utils.FakeStringFunc()),
							jen.ID("ClientSecret").MapAssign().Add(utils.FakeUUIDFunc()),
						),
					),
				),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/api/v1/oauth2/clients"),
								jen.Lit("expected and actual path don't match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
								nil,
							),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("GetOAuth2Clients").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_BuildCreateOAuth2ClientRequest(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildCreateOAuth2ClientRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").MapAssign().Qual(filepath.Join(outPath, "models/v1"), "UserLoginInput").Valuesln(
						jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
						jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
						jen.ID("TOTPToken").MapAssign().Lit("123456"),
					),
				),
				jen.List(
					jen.ID("req"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildCreateOAuth2ClientRequest").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertNoError(
					jen.Err(),
					nil,
				),
				utils.AssertEqual(
					jen.Qual("net/http", "MethodPost"),
					jen.ID("req").Dot("Method"),
					nil,
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_CreateOAuth2Client(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_CreateOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").MapAssign().Qual(filepath.Join(outPath, "models/v1"), "UserLoginInput").Valuesln(
						jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
						jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
						jen.ID("TOTPToken").MapAssign().Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("exampleOutput").Assign().VarPointer().Qual(filepath.Join(outPath, "models/v1"), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("EXAMPLECLIENTID"),
					jen.ID("ClientSecret").MapAssign().Lit("EXAMPLECLIENTSECRET"),
				),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.Lit("/oauth2/client"),
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("expected and actual path don't match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleOutput")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.Err(),
				).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertNoError(
					jen.Err(),
					nil,
				),
				utils.AssertNotNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid body",
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").MapAssign().Qual(filepath.Join(outPath, "models/v1"), "UserLoginInput").Valuesln(
						jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
						jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
						jen.ID("TOTPToken").MapAssign().Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/oauth2/client"),
								jen.Lit("expected and actual path don't match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							jen.List(
								jen.ID("_"),
								jen.Err(),
							).Assign().ID("res").Dot("Write").Call(
								jen.Index().ID("byte").Call(
									jen.Lit("BLAH"),
								),
							),
							utils.AssertNoError(
								jen.Err(),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.Err(),
				).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
				utils.AssertNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").MapAssign().Qual(filepath.Join(outPath, "models/v1"), "UserLoginInput").Valuesln(
						jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
						jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
						jen.ID("TOTPToken").MapAssign().Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/oauth2/client"),
								jen.Lit("expected and actual path don't match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							jen.Qual("time", "Sleep").Call(
								jen.Lit(10).Times().Qual("time", "Hour"),
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Millisecond"),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.Err(),
				).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
				utils.AssertNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with 404",
				jen.ID("exampleInput").Assign().VarPointer().Qual(filepath.Join(outPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").MapAssign().Qual(filepath.Join(outPath, "models/v1"), "UserLoginInput").Valuesln(
						jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
						jen.ID("Password").MapAssign().Add(utils.FakePasswordFunc()),
						jen.ID("TOTPToken").MapAssign().Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/oauth2/client"),
								jen.Lit("expected and actual path don't match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							utils.WriteHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.Err(),
				).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
				utils.AssertEqual(jen.Err(),
					jen.ID("ErrNotFound"), nil),
				utils.AssertNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with no cookie",
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("_"),
					jen.Err(),
				).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.Nil(),
					jen.Nil(),
				),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_BuildArchiveOAuth2ClientRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildArchiveOAuth2ClientRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodDelete"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				jen.ID("expectedID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
					utils.CtxVar(),
					jen.ID("expectedID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.RequireNotNil(jen.ID("actual").Dot("URL"), nil),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_ArchiveOAuth2Client() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_ArchiveOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").ParamPointer().Qual("net/http", "Request"),
						).Block(
							utils.AssertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("/api/v1/oauth2/clients/%d"),
									jen.ID("expected"),
								),
								jen.Lit("expected and actual path don't match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodDelete"),
								nil,
							),
							jen.Line(),
							utils.WriteHeader("StatusOK"),
						),
					),
				),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				).Dot("ArchiveOAuth2Client").Call(
					utils.CtxVar(),
					jen.ID("expected"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}
