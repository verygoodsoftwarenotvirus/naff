package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func usersTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		testFunc("V1Client_BuildGetUserRequest").Block(
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
				).Op(":=").ID("c").Dot("BuildGetUserRequest").Call(
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
		testFunc("V1Client_GetUser").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "User").Values(jen.Dict{
					jen.ID("ID"): jen.Lit(1),
				}),
				jen.Line(),
				createCtx(),
				jen.Line(),
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
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/users/%d"),
								jen.ID("expected").Dot("ID"),
							),
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
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetUser").Call(
					jen.ID("ctx"),
					jen.ID("expected").Dot("ID"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID(t),
					jen.ID("actual"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("expected"),
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildGetUsersRequest").Block(
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
				).Op(":=").ID("c").Dot("BuildGetUsersRequest").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
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
		testFunc("V1Client_GetUsers").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "UserList").Values(jen.Dict{
					jen.ID("Users"): jen.Index().Qual(modelsPkg, "User").Values(jen.Dict{
						jen.ID("ID"): jen.Lit(1),
					}),
				}),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.ID("res").Qual("net/http", "ResponseWriter"),
						jen.ID("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.Lit("/users"),
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
				).Op(":=").ID("c").Dot("GetUsers").Call(
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
		testFunc("V1Client_BuildCreateUserRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "UserInput").Values(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildCreateUserRequest").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
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
		testFunc("V1Client_CreateUser").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "UserCreationResponse").Values(jen.Dict{
					jen.ID("ID"): jen.Lit(1),
				}),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "UserInput").Values(),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.ID("res").Qual("net/http", "ResponseWriter"),
						jen.ID("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.Lit("/users"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.ID("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
							nil,
						),
						jen.Var().ID("x").Op("*").Qual(modelsPkg, "UserInput"),
						jen.ID("require").Dot("NoError").Call(
							jen.ID(t),
							jen.Qual("encoding/json", "NewDecoder").Call(
								jen.ID("req").Dot("Body"),
							).Dot("Decode").Call(
								jen.Op("&").ID("x"),
							),
						),
						assertEqual(
							jen.ID("exampleInput"),
							jen.ID("x"),
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
						writeHeader("StatusOK"),
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
				).Op(":=").ID("c").Dot("CreateUser").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("expected"),
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildArchiveUserRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildArchiveUserRequest").Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				requireNotNil(jen.ID("actual"), nil),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID(t),
					jen.ID("actual").Dot("URL"),
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
		testFunc("V1Client_ArchiveUser").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.ID("res").Qual("net/http", "ResponseWriter"),
						jen.ID("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/users/%d"),
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
				).Dot("ArchiveUser").Call(
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

	ret.Add(
		testFunc("V1Client_BuildLoginRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildLoginRequest").Call(
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				requireNotNil(jen.ID("req"), nil),
				assertEqual(
					jen.ID("req").Dot("Method"),
					jen.Qual("net/http", "MethodPost"),
					nil,
				),
				assertNoError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_Login").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.ID("res").Qual("net/http", "ResponseWriter"),
						jen.ID("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.Lit("/users/login"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.ID("req").Dot("Method"),
							jen.Qual("net/http", "MethodPost"),
							nil,
						),
						jen.Qual("net/http", "SetCookie").Call(
							jen.ID("res"),
							jen.Op("&").Qual("net/http", "Cookie").Values(jen.Dict{
								jen.ID("Name"): jen.Lit("hi"),
							}),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("cookie"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("Login").Call(
					jen.ID("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				requireNotNil(jen.ID("cookie"), nil),
				assertNoError(
					jen.ID("err"),
					nil,
				),
			),
			buildSubTest(
				"with timeout",
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/users/login"),
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
							writeHeader("StatusOK"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Microsecond"),
				jen.List(
					jen.ID("cookie"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("Login").Call(
					jen.ID("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				jen.ID("require").Dot("Nil").Call(
					jen.ID(t),
					jen.ID("cookie"),
				),
				assertError(
					jen.ID("err"),
					nil,
				),
			),
			buildSubTest(
				"with missing cookie",
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit("/users/login"),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							writeHeader("StatusOK"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("cookie"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("Login").Call(
					jen.ID("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				jen.ID("require").Dot("Nil").Call(
					jen.ID(t),
					jen.ID("cookie"),
				),
				assertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	return ret
}
