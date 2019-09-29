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
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
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
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "OAuth2Client").Values(jen.Dict{
					jen.ID("ID"):           jen.Lit(1),
					jen.ID("ClientID"):     jen.Lit("example"),
					jen.ID("ClientSecret"): jen.Lit("blah"),
				}),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
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
						jen.ID("require").Dot("NoError").Call(
							jen.ID(t),
							jen.Qual("encoding/json", "NewEncoder").Call(
								jen.ID("res"),
							).Dot("Encode").Call(
								jen.ID("expected"),
							),
						),
					),
				),
				),
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
				createCtx(),
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
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientList").Values(jen.Dict{
					jen.ID("Clients"): jen.Index().Qual(modelsPkg, "OAuth2Client").Values(jen.Dict{
						jen.ID("ID"):           jen.Lit(1),
						jen.ID("ClientID"):     jen.Lit("example"),
						jen.ID("ClientSecret"): jen.Lit("blah"),
					}),
				}),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
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
							jen.ID("require").Dot("NoError").Call(
								jen.ID(t),
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(
									jen.ID("expected"),
								),
							),
						),
					),
				),
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
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").Values(jen.Dict{
					jen.ID("UserLoginInput"): jen.Qual(modelsPkg, "UserLoginInput").Values(jen.Dict{
						jen.ID("Username"):  jen.Lit("username"),
						jen.ID("Password"):  jen.Lit("password"),
						jen.ID("TOTPToken"): jen.Lit("123456"),
					}),
				}),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildCreateOAuth2ClientRequest").Call(
					jen.ID("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.ID("exampleInput"),
				),
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
				createCtx(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").Values(jen.Dict{
					jen.ID("UserLoginInput"): jen.Qual(modelsPkg, "UserLoginInput").Values(jen.Dict{
						jen.ID("Username"):  jen.Lit("username"),
						jen.ID("Password"):  jen.Lit("password"),
						jen.ID("TOTPToken"): jen.Lit("123456"),
					}),
				}),
				jen.ID("exampleOutput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2Client").Values(jen.Dict{
					jen.ID("ClientID"):     jen.Lit("EXAMPLECLIENTID"),
					jen.ID("ClientSecret"): jen.Lit("EXAMPLECLIENTSECRET"),
				}),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
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
							jen.ID("require").Dot("NoError").Call(
								jen.ID(t),
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(
									jen.ID("exampleOutput"),
								),
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
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
			buildSubTest(
				"with invalid body",
				createCtx(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").Values(jen.Dict{
					jen.ID("UserLoginInput"): jen.Qual(modelsPkg, "UserLoginInput").Values(jen.Dict{
						jen.ID("Username"):  jen.Lit("username"),
						jen.ID("Password"):  jen.Lit("password"),
						jen.ID("TOTPToken"): jen.Lit("123456"),
					}),
				}),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
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
									jen.Lit(`

					BLAH

									`),
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
			buildSubTest(
				"with timeout",
				createCtx(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").Values(jen.Dict{
					jen.ID("UserLoginInput"): jen.Qual(modelsPkg, "UserLoginInput").Values(jen.Dict{
						jen.ID("Username"):  jen.Lit("username"),
						jen.ID("Password"):  jen.Lit("password"),
						jen.ID("TOTPToken"): jen.Lit("123456"),
					}),
				}),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
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
			buildSubTest(
				"with 404",
				createCtx(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "OAuth2ClientCreationInput").Values(jen.Dict{
					jen.ID("UserLoginInput"): jen.Qual(modelsPkg, "UserLoginInput").Values(jen.Dict{
						jen.ID("Username"):  jen.Lit("username"),
						jen.ID("Password"):  jen.Lit("password"),
						jen.ID("TOTPToken"): jen.Lit("123456"),
					}),
				}),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
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
			buildSubTest(
				"with no cookie",
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
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
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
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
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
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
							writeHeader("StatusOK"),
						),
					),
				),
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
