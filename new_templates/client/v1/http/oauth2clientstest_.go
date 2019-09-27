package client

import jen "github.com/dave/jennifer/jen"

func oauth2ClientsTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		testFunc("V1Client_BuildGetOAuth2ClientRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("expectedID").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildGetOAuth2ClientRequest").Call(
					jen.Id("ctx"),
					jen.Id("expectedID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertTrue(
					jen.Id("t"),
					jen.Qual("strings", "HasSuffix").Call(
						jen.Id("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.Id("expectedID"),
						),
					),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_GetOAuth2Client").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("OAuth2Client").Values(
					jen.Id("ID").Op(":").Lit(1),
					jen.Id("ClientID").Op(":").Lit("example"),
					jen.Id("ClientSecret").Op(":").Lit("blah"),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertTrue(
							jen.Id("t"),
							jen.Qual("strings", "HasSuffix").Call(
								jen.Id("req").Dot("URL").Dot("String").Call(),
								jen.Qual("strconv", "Itoa").Call(
									jen.Id("int").Call(
										jen.Id("expected").Dot("ID"),
									),
								),
							),
						),
						assertEqual(
							jen.Id("t"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/api/v1/oauth2/clients/%d"),
								jen.Id("expected").Dot("ID"),
							),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodGet"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("GetOAuth2Client").Call(
					jen.Id("ctx"),
					jen.Id("expected").Dot("ID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildGetOAuth2ClientsRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildGetOAuth2ClientsRequest").Call(
					jen.Id("ctx"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_GetOAuth2Clients").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("OAuth2ClientList").Values(
					jen.Id("Clients").Op(":").Index().Id("models").Dot("OAuth2Client").Values(
						jen.Values(jen.Dict{
							jen.Id("ID"):           jen.Lit(1),
							jen.Id("ClientID"):     jen.Lit("example"),
							jen.Id("ClientSecret"): jen.Lit("blah"),
						}),
					),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/api/v1/oauth2/clients"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodGet"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("GetOAuth2Clients").Call(
					jen.Id("ctx"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildCreateOAuth2ClientRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(
					jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(
						jen.Id("Username").Op(":").Lit("username"),
						jen.Id("Password").Op(":").Lit("password"),
						jen.Id("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildCreateOAuth2ClientRequest").Call(
					jen.Id("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.Id("exampleInput"),
				),
				requireNotNil(jen.Id("req"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Qual("net/http", "MethodPost"),
					jen.Id("req").Dot("Method"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_CreateOAuth2Client").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				createCtx(),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(
					jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(
						jen.Id("Username").Op(":").Lit("username"),
						jen.Id("Password").Op(":").Lit("password"),
						jen.Id("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Id("exampleOutput").Op(":=").Op("&").Id("models").Dot("OAuth2Client").Values(
					jen.Id("ClientID").Op(":").Lit("EXAMPLECLIENTID"),
					jen.Id("ClientSecret").Op(":").Lit("EXAMPLECLIENTSECRET"),
				),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Lit("/oauth2/client"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("exampleOutput"),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("oac"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(
					jen.Id("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.Id("exampleInput"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
				assertNotNil(
					jen.Id("t"),
					jen.Id("oac"),
				),
			),
			buildSubTest(
				"with invalid body",
				createCtx(),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(
					jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(
						jen.Id("Username").Op(":").Lit("username"),
						jen.Id("Password").Op(":").Lit("password"),
						jen.Id("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/oauth2/client"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
						),
						jen.List(
							jen.Id("_"),
							jen.Id("err"),
						).Op(":=").Id("res").Dot("Write").Call(
							jen.Index().Id("byte").Call(
								jen.Lit(`

					BLAH

									`),
							),
						),
						assertNoError(
							jen.Id("t"),
							jen.Id("err"),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("oac"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(
					jen.Id("ctx"), jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.Id("exampleInput"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
				assertNil(
					jen.Id("t"),
					jen.Id("oac"),
				),
			),
			buildSubTest(
				"with timeout",
				createCtx(),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(
					jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(
						jen.Id("Username").Op(":").Lit("username"),
						jen.Id("Password").Op(":").Lit("password"),
						jen.Id("TOTPToken").Op(":").Lit("123456"),
					),
				),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/oauth2/client"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
						),
						jen.Qual("time", "Sleep").Call(
							jen.Lit(10).Op("*").Qual("time", "Hour"),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Millisecond"),
				jen.List(
					jen.Id("oac"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(
					jen.Id("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.Id("exampleInput"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
				assertNil(
					jen.Id("t"),
					jen.Id("oac"),
				),
			),
			buildSubTest(
				"with 404",
				createCtx(),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("OAuth2ClientCreationInput").Values(
					jen.Id("UserLoginInput").Op(":").Id("models").Dot("UserLoginInput").Values(jen.Dict{
						jen.Id("Username"):  jen.Lit("username"),
						jen.Id("Password"):  jen.Lit("password"),
						jen.Id("TOTPToken"): jen.Lit("123456"),
					}),
				),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/oauth2/client"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
						),
						writeHeader("StatusNotFound"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("oac"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(
					jen.Id("ctx"),
					jen.Op("&").Qual("net/http", "Cookie").Values(),
					jen.Id("exampleInput"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
				assertEqual(jen.Id("err"), jen.Id("ErrNotFound"), nil),
				assertNil(
					jen.Id("t"),
					jen.Id("oac"),
				),
			),
			buildSubTest(
				"with no cookie",
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("_"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("CreateOAuth2Client").Call(
					jen.Id("ctx"),
					jen.Id("nil"),
					jen.Id("nil"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildArchiveOAuth2ClientRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("expectedID").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
					jen.Id("ctx"),
					jen.Id("expectedID"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_ArchiveOAuth2Client").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/api/v1/oauth2/clients/%d"),
								jen.Id("expected"),
							),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodDelete")), writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("err").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				).Dot("ArchiveOAuth2Client").Call(
					jen.Id("ctx"),
					jen.Id("expected"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
