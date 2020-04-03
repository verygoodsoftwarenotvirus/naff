package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildGetUserRequest").Block(
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
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetUserRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
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
							jen.ID("exampleUser").Dot("ID"),
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
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Comment("the hashed password is never transmitted over the wire"),
				jen.ID("exampleUser").Dot("HashedPassword").Equals().Lit(""),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertTrue(
						jen.Qual("strings", "HasSuffix").Call(
							jen.ID("req").Dot("URL").Dot("String").Call(),
							jen.Qual("strconv", "Itoa").Call(
								jen.ID("int").Call(
									jen.ID("exampleUser").Dot("ID"),
								),
							),
						),
						nil,
					),
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/users/%d"),
							jen.ID("exampleUser").Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodGet"), nil),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleUser")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID("exampleUser"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleUser").Dot("Salt").Equals().Nil(),
				jen.ID("exampleUser").Dot("HashedPassword").Equals().Lit(""),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
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
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetUsersRequest").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
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
				jen.ID("exampleUserList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserList").Call(),
				jen.Comment("the hashed password is never transmitted over the wire"),
				jen.ID("exampleUserList").Dot("Users").Index(jen.Lit(0)).Dot("HashedPassword").Equals().Lit(""),
				jen.ID("exampleUserList").Dot("Users").Index(jen.Lit(1)).Dot("HashedPassword").Equals().Lit(""),
				jen.ID("exampleUserList").Dot("Users").Index(jen.Lit(2)).Dot("HashedPassword").Equals().Lit(""),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleUserList")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID("exampleUserList"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
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
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserCreationInputFromUser").Call(jen.ID("exampleUser")),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildCreateUserRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
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
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_CreateUser").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserCreationInputFromUser").Call(jen.ID("exampleUser")),
				jen.ID("expected").Assign().Qual(proj.FakeModelsPackage(), "BuildDatabaseCreationResponse").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Var().ID("x").PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput"),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")).Dot("Decode").Call(jen.VarPointer().ID("x")),
						nil,
					),
					utils.AssertEqual(jen.ID("exampleInput"), jen.ID("x"), nil),
					jen.Line(),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserCreationInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
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
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildArchiveUserRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.RequireNotNil(jen.ID("actual").Dot("URL"), nil),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("exampleUser").Dot("ID"),
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
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_ArchiveUser").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/users/%d"),
							jen.ID("exampleUser").Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodDelete"),
						nil,
					),
				),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot("ArchiveUser").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
				),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("ArchiveUser").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
				),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc("V1Client_BuildLoginRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildLoginRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertEqual(
					jen.ID("req").Dot("Method"),
					jen.Qual("net/http", "MethodPost"),
					nil,
				),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil input",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildLoginRequest").Call(utils.CtxVar(), jen.Nil()),
				utils.AssertNil(jen.ID("req"), nil),
				utils.AssertError(jen.Err(), nil),
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
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users/login"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Qual("net/http", "SetCookie").Call(
						jen.ID("res"),
						jen.VarPointer().Qual("net/http", "Cookie").Values(
							jen.ID("Name").MapAssign().ID("exampleUser").Dot("Username"),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("c").Dot("Login").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				utils.RequireNotNil(jen.ID("cookie"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil input",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("c").Dot("Login").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				utils.AssertNil(jen.ID("cookie"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("c").Dot("Login").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				utils.AssertNil(jen.ID("cookie"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users/login"),
						jen.Lit("expected and actual paths do not match"),
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
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Equals().Lit(500).Times().Qual("time", "Microsecond"),
				jen.Line(),
				jen.List(
					jen.ID("cookie"),
					jen.Err(),
				).Assign().ID("c").Dot("Login").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				utils.RequireNil(jen.ID("cookie"), nil),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with missing cookie",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit("/users/login"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.Line(),
				jen.List(
					jen.ID("cookie"),
					jen.Err(),
				).Assign().ID("c").Dot("Login").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				utils.RequireNil(jen.ID("cookie"), nil),
				utils.AssertError(
					jen.Err(),
					nil,
				),
			),
		),
		jen.Line(),
	)

	return ret
}
