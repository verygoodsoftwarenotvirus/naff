package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildGetUserRequest").Block(
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
				).Op(":=").ID("c").Dot("BuildGetUserRequest").Call(
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
		utils.OuterTestFunc("V1Client_GetUser").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Valuesln(
					jen.ID("ID").Op(":").Lit(1),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
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
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/users/%d"),
							jen.ID("expected").Dot("ID"),
						),
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
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetUser").Call(
					jen.ID("ctx"),
					jen.ID("expected").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(
					jen.ID("actual"),
					nil,
				),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
					jen.ID("expected"),
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildGetUsersRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildGetUsersRequest").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
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
		utils.OuterTestFunc("V1Client_GetUsers").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserList").Values(
					jen.ID("Users").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(
						jen.Values(jen.ID("ID").Op(":").Lit(1)),
					),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users"),
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
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetUsers").Call(
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
		utils.OuterTestFunc("V1Client_BuildCreateUserRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput").Values(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildCreateUserRequest").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
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
		utils.OuterTestFunc("V1Client_CreateUser").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserCreationResponse").Values(jen.ID("ID").Op(":").Lit(1)),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput").Values(),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users"),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Var().ID("x").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput"),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewDecoder").Call(
							jen.ID("req").Dot("Body"),
						).Dot("Decode").Call(
							jen.Op("&").ID("x"),
						),
						nil,
					),
					utils.AssertEqual(
						jen.ID("exampleInput"),
						jen.ID("x"),
						nil,
					),
					jen.Line(),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
						nil,
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
				).Op(":=").ID("c").Dot("CreateUser").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
					jen.ID("expected"),
					jen.ID("actual"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildArchiveUserRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodDelete"),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildArchiveUserRequest").Call(
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
		utils.OuterTestFunc("V1Client_ArchiveUser").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/users/%d"),
							jen.ID("expected"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodDelete"),
						nil,
					),
				),
				jen.Line(),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				).Dot("ArchiveUser").Call(
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

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildLoginRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("req"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildLoginRequest").Call(
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertEqual(
					jen.ID("req").Dot("Method"),
					jen.Qual("net/http", "MethodPost"),
					nil,
				),
				utils.AssertNoError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_Login").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users/login"),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Qual("net/http", "SetCookie").Call(
						jen.ID("res"),
						jen.Op("&").Qual("net/http", "Cookie").Values(
							jen.ID("Name").Op(":").Lit("hi"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("cookie"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("Login").Call(
					jen.ID("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				utils.RequireNotNil(jen.ID("cookie"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users/login"),
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
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Op("=").Lit(500).Op("*").Qual("time", "Microsecond"),
				jen.Line(),
				jen.List(
					jen.ID("cookie"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("Login").Call(
					jen.ID("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				utils.RequireNil(jen.ID("cookie"), nil),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with missing cookie",
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users/login"),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("cookie"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("Login").Call(
					jen.ID("ctx"),
					jen.Lit("username"),
					jen.Lit("password"),
					jen.Lit("123456"),
				),
				utils.RequireNil(jen.ID("cookie"), nil),
				utils.AssertError(
					jen.ID("err"),
					nil,
				),
			),
		),
		jen.Line(),
	)

	return ret
}
