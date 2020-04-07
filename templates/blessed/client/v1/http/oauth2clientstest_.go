package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(buildV1Client_BuildGetOAuth2ClientRequest(proj)...)
	ret.Add(buildV1Client_GetOAuth2Client(proj)...)
	ret.Add(buildV1Client_BuildGetOAuth2ClientsRequest()...)
	ret.Add(buildV1Client_GetOAuth2Clients(proj)...)
	ret.Add(buildV1Client_BuildCreateOAuth2ClientRequest(proj)...)
	ret.Add(buildV1Client_CreateOAuth2Client(proj)...)
	ret.Add(buildV1Client_BuildArchiveOAuth2ClientRequest(proj)...)
	ret.Add(buildV1Client_ArchiveOAuth2Client(proj)...)

	return ret
}

func buildV1Client_BuildGetOAuth2ClientRequest(proj *models.Project) []jen.Code {
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
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleOAuth2Client").Dot("ID"),
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
							jen.ID("exampleOAuth2Client").Dot("ID"),
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

	happyPathSubtestLines := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
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
								jen.Int().Call(
									jen.ID("exampleOAuth2Client").Dot("ID"),
								),
							),
						),
						nil,
					),
					utils.AssertEqual(
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/api/v1/oauth2/clients/%d"),
							jen.ID("exampleOAuth2Client").Dot("ID"),
						),
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleOAuth2Client")),
						nil,
					),
				),
			),
		),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(
			utils.CtxVar(),
			jen.ID("exampleOAuth2Client").Dot("ID"),
		),
		jen.Line(),
		utils.RequireNotNil(jen.ID("actual"), nil),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
		utils.AssertEqual(jen.ID("exampleOAuth2Client"), jen.ID("actual"), nil),
	}

	invalidClientURLSubtestLines := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
		jen.Line(),
		jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(
			utils.CtxVar(),
			jen.ID("exampleOAuth2Client").Dot("ID"),
		),
		jen.Line(),
		utils.AssertNil(jen.ID("actual"), nil),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	}

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_GetOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", invalidClientURLSubtestLines...),
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
	happyPathSubtestLines := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2ClientList"),
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
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleOAuth2ClientList")),
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
		utils.AssertEqual(jen.ID("exampleOAuth2ClientList"), jen.ID("actual"), nil),
	}

	invalidClientURLSubtestLines := []jen.Code{
		jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(
			utils.CtxVar(),
			jen.Nil(),
		),
		jen.Line(),
		utils.AssertNil(jen.ID("actual"), nil),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	}

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_GetOAuth2Clients").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", invalidClientURLSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_BuildCreateOAuth2ClientRequest(proj *models.Project) []jen.Code {
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
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2ClientCreationInputFromClient").Call(jen.ID("exampleOAuth2Client")),
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
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_CreateOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID("exampleOAuth2Client")),
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
								jen.Lit("expected and actual paths do not match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							utils.RequireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleOAuth2Client")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleOAuth2Client"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID("exampleOAuth2Client")),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid response from server",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID("exampleOAuth2Client")),
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
								jen.Lit("expected and actual paths do not match"),
							),
							utils.AssertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							jen.List(jen.Underscore(), jen.Err()).Assign().ID("res").Dot("Write").Call(jen.Index().Byte().Call(jen.Lit("BLAH"))),
							utils.AssertNoError(jen.Err(), nil),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.VarPointer().Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"without cookie",
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(
					utils.CtxVar(),
					jen.Nil(),
					jen.Nil(),
				),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildV1Client_BuildArchiveOAuth2ClientRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildArchiveOAuth2ClientRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodDelete"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleOAuth2Client").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.RequireNotNil(jen.ID("actual").Dot("URL"), nil),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("exampleOAuth2Client").Dot("ID"),
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

func buildV1Client_ArchiveOAuth2Client(proj *models.Project) []jen.Code {
	happyPathSubtestLines := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
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
							jen.ID("exampleOAuth2Client").Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
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
		).Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleOAuth2Client").Dot("ID")),
		utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
	}

	invalidClientURLSubtestLines := []jen.Code{
		utils.BuildFakeVar(proj, "OAuth2Client"),
		jen.Line(),

		jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(
			jen.ID("t"),
		).Dot("ArchiveOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleOAuth2Client").Dot("ID")),
		utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
	}

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_ArchiveOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", happyPathSubtestLines...),
			jen.Line(),
			utils.BuildSubTest("with invalid client URL", invalidClientURLSubtestLines...),
		),
		jen.Line(),
	}

	return lines
}
