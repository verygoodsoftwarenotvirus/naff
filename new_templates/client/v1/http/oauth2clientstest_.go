package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func oauth2ClientsTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		testFunc("V1Client_BuildGetOAuth2ClientRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				assertEqual(
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
		testFunc("V1Client_GetOAuth2Client").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "OAuth2Client").ValuesLn(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("ClientID").Op(":").Lit("example"),
					jen.ID("ClientSecret").Op(":").Lit("blah"),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertTrue(
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
							assertEqual(
								jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("/api/v1/oauth2/clients/%d"),
									jen.ID("expected").Dot("ID"),
								),
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							requireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
								nil,
							),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildGetOAuth2ClientsRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
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
		testFunc("V1Client_GetOAuth2Clients").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientList").ValuesLn(
					jen.ID("Clients").Op(":").Index().Qual(modelsPkg, "OAuth2Client").ValuesLn(
						jen.ValuesLn(
							jen.ID("ID").Op(":").Lit(1),
							jen.ID("ClientID").Op(":").Lit("example"),
							jen.ID("ClientSecret").Op(":").Lit("blah"),
						),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/api/v1/oauth2/clients"),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							requireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
								nil,
							),
						),
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildCreateOAuth2ClientRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").ValuesLn(
					jen.ID("UserLoginInput").Op(":").Qual(modelsPkg, "UserLoginInput").ValuesLn(
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
				requireNotNil(jen.ID("req"), nil),
				assertNoError(
					jen.ID("err"),
					nil,
				),
				assertEqual(
					jen.Qual("net/http", "MethodPost"),
					jen.ID("req").Dot("Method"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_CreateOAuth2Client").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").ValuesLn(
					jen.ID("UserLoginInput").Op(":").Qual(modelsPkg, "UserLoginInput").ValuesLn(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("exampleOutput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2Client").ValuesLn(
					jen.ID("ClientID").Op(":").Lit("EXAMPLECLIENTID"),
					jen.ID("ClientSecret").Op(":").Lit("EXAMPLECLIENTSECRET"),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.Lit("/oauth2/client"),
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							requireNoError(
								jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleOutput")),
								nil,
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				assertNoError(
					jen.ID("err"),
					nil,
				),
				assertNotNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with invalid body",
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").ValuesLn(
					jen.ID("UserLoginInput").Op(":").Qual(modelsPkg, "UserLoginInput").ValuesLn(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/oauth2/client"),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
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
							assertNoError(
								jen.ID("err"),
								nil,
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				assertError(
					jen.ID("err"),
					nil,
				),
				assertNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with timeout",
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").ValuesLn(
					jen.ID("UserLoginInput").Op(":").Qual(modelsPkg, "UserLoginInput").ValuesLn(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/oauth2/client"),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
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
					jen.ID(t),
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
				assertError(
					jen.ID("err"),
					nil,
				),
				assertNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with 404",
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").ValuesLn(
					jen.ID("UserLoginInput").Op(":").Qual(modelsPkg, "UserLoginInput").ValuesLn(
						jen.ID("Username").Op(":").Lit("username"),
						jen.ID("Password").Op(":").Lit("password"),
						jen.ID("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/oauth2/client"),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							writeHeader("StatusNotFound"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				assertError(
					jen.ID("err"),
					nil,
				),
				assertEqual(jen.ID("err"),
					jen.ID("ErrNotFound"), nil),
				assertNil(
					jen.ID("oac"),
					nil,
				),
			),
			jen.Line(),
			buildSubTest(
				"with no cookie",
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				assertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildArchiveOAuth2ClientRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
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
				requireNotNil(jen.ID("actual"), nil),
				requireNotNil(jen.ID("actual").Dot("URL"), nil),
				assertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
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
		testFunc("V1Client_ArchiveOAuth2Client").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").CallLn(
					jen.Qual("net/http", "HandlerFunc").CallLn(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("/api/v1/oauth2/clients/%d"),
									jen.ID("expected"),
								),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodDelete"),
								nil,
							),
							jen.Line(),
							writeHeader("StatusOK"),
						),
					),
				),
				jen.Line(),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot("ArchiveOAuth2Client").Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
