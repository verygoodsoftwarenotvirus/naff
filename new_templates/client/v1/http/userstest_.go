package client

import jen "github.com/dave/jennifer/jen"

func usersTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		testFunc("V1Client_BuildGetUserRequest").Block(
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
				).Op(":=").Id("c").Dot("BuildGetUserRequest").Call(
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
		testFunc("V1Client_GetUser").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("User").Values(
					jen.Id("ID").Op(":").Lit(1),
				),
				jen.Line(),
				createCtx(),
				jen.Line(),
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
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/users/%d"),
								jen.Id("expected").Dot("ID"),
							),
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
				jen.Line(),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("GetUser").Call(
					jen.Id("ctx"),
					jen.Id("expected").Dot("ID"),
				),
				jen.Id("require").Dot("NotNil").Call(
					jen.Id("t"),
					jen.Id("actual"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("expected"),
					jen.Id("actual"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildGetUsersRequest").Block(
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
				).Op(":=").Id("c").Dot("BuildGetUsersRequest").Call(
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
		testFunc("V1Client_GetUsers").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("UserList").Values(
					jen.Id("Users").Op(":").Index().Id("models").Dot("User").Values(
						jen.Values(
							jen.Id("ID").Op(":").Lit(1),
						),
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
							jen.Lit("/users"),
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
				).Op(":=").Id("c").Dot("GetUsers").Call(
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
		testFunc("V1Client_BuildCreateUserRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("UserInput").Values(),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildCreateUserRequest").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
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
		testFunc("V1Client_CreateUser").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("UserCreationResponse").Values(
					jen.Id("ID").Op(":").Lit(1),
				),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("UserInput").Values(),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/users"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
						),
						jen.Var().Id("x").Op("*").Id("models").Dot("UserInput"),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"),
							jen.Qual("encoding/json", "NewDecoder").Call(
								jen.Id("req").Dot("Body"),
							).Dot("Decode").Call(
								jen.Op("&").Id("x"),
							),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("exampleInput"),
							jen.Id("x"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
						writeHeader("StatusOK"),
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
				).Op(":=").Id("c").Dot("CreateUser").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
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
		testFunc("V1Client_BuildArchiveUserRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				jen.Id("expectedID").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildArchiveUserRequest").Call(
					jen.Id("ctx"),
					jen.Id("expectedID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				jen.Id("require").Dot("NotNil").Call(
					jen.Id("t"),
					jen.Id("actual").Dot("URL"),
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
		testFunc("V1Client_ArchiveUser").Block(
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
								jen.Lit("/users/%d"),
								jen.Id("expected"),
							),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodDelete"),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("err").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				).Dot("ArchiveUser").Call(
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

	ret.Add(
		testFunc("V1Client_BuildLoginRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("req"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildLoginRequest").Call(
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				requireNotNil(jen.Id("req"), nil),
				assertEqual(
					jen.Id("t"),
					jen.Id("req").Dot("Method"),
					jen.Qual("net/http", "MethodPost"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_Login").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/users/login"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
						),
						jen.Qual("net/http", "SetCookie").Call(
							jen.Id("res"),
							jen.Op("&").Qual("net/http", "Cookie").Values(
								jen.Id("Name").Op(":").Lit("hi"),
							),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("cookie"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("Login").Call(
					jen.Id("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				requireNotNil(jen.Id("cookie"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
			buildSubTest(
				"with timeout",
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/users/login"),
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
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Microsecond"),
				jen.List(
					jen.Id("cookie"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("Login").Call(
					jen.Id("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				jen.Id("require").Dot("Nil").Call(
					jen.Id("t"),
					jen.Id("cookie"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
			buildSubTest(
				"with missing cookie",
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/users/login"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("cookie"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("Login").Call(
					jen.Id("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				jen.Id("require").Dot("Nil").Call(
					jen.Id("t"),
					jen.Id("cookie"),
				),
				assertError(
					jen.Id("t"),
					jen.Id("err"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
