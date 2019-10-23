package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(ret)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildGetOAuth2ClientRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
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
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_GetOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2Client").Valuesln(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("ClientID").Op(":").Lit("example"),
					jen.ID("ClientSecret").Op(":").Lit("blah"),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
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
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetOAuth2Client").Call(
					jen.ID("ctx"),
					jen.ID("expected").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildGetOAuth2ClientsRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
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
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_GetOAuth2Clients").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2ClientList").Valuesln(
					jen.ID("Clients").Op(":").Index().Qual(utils.ModelsPkg, "OAuth2Client").Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(1),
							jen.ID("ClientID").Op(":").Lit("example"),
							jen.ID("ClientSecret").Op(":").Lit("blah"),
						),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
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
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetOAuth2Clients").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildCreateOAuth2ClientRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").Op(":").Qual(utils.ModelsPkg, "UserLoginInput").Valuesln(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildCreateOAuth2ClientRequest").Call(
					jen.ID("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertNoError(
					jen.ID("err"),
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
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_CreateOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").Op(":").Qual(utils.ModelsPkg, "UserLoginInput").Valuesln(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("exampleOutput").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2Client").Valuesln(
					jen.ID("ClientID").Op(":").Lit("EXAMPLECLIENTID"),
					jen.ID("ClientSecret").Op(":").Lit("EXAMPLECLIENTSECRET"),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
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
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("CreateOAuth2Client").Call(
					jen.ID("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertNoError(
					jen.ID("err"),
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
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").Op(":").Qual(utils.ModelsPkg, "UserLoginInput").Valuesln(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
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
								jen.ID("err"),
							).Op(":=").ID("res").Dot("Write").Call(
								jen.Index().ID("byte").Call(
									jen.Lit("BLAH"),
								),
							),
							utils.AssertNoError(
								jen.ID("err"),
								nil,
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("CreateOAuth2Client").Call(
					jen.ID("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertError(
					jen.ID("err"),
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
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").Op(":").Qual(utils.ModelsPkg, "UserLoginInput").Valuesln(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
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
								jen.Lit(10).Op("*").Qual("time", "Hour"),
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("CreateOAuth2Client").Call(
					jen.ID("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertError(
					jen.ID("err"),
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
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, "OAuth2ClientCreationInput").Valuesln(
					jen.ID("UserLoginInput").Op(":").Qual(utils.ModelsPkg, "UserLoginInput").Valuesln(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
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
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("oac"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("CreateOAuth2Client").Call(
					jen.ID("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
				utils.AssertEqual(jen.ID("err"),
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
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("_"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("CreateOAuth2Client").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
					jen.ID("nil"),
				),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildArchiveOAuth2ClientRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodDelete"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
					jen.ID("ctx"),
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
					jen.ID("err"),
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
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_ArchiveOAuth2Client").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
					jen.Qual("net/http", "HandlerFunc").Callln(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
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
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				).Dot("ArchiveOAuth2Client").Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
