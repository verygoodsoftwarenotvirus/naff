package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildTestV1ClientBuildGetUserRequest(proj)...)
	code.Add(buildTestV1ClientGetUser(proj)...)
	code.Add(buildTestV1ClientBuildGetUsersRequest()...)
	code.Add(buildTestV1ClientGetUsers(proj)...)
	code.Add(buildTestV1ClientBuildCreateUserRequest(proj)...)
	code.Add(buildTestV1ClientCreateUser(proj)...)
	code.Add(buildTestV1ClientBuildArchiveUserRequest(proj)...)
	code.Add(buildTestV1ClientArchiveUser(proj)...)
	code.Add(buildTestV1ClientBuildLoginRequest(proj)...)
	code.Add(buildTestV1ClientLogin(proj)...)
	code.Add(buildTestV1ClientBuildValidateTOTPSecretRequest(proj)...)
	code.Add(buildTestV1ClientValidateTOTPSecret(proj)...)

	return code
}

func buildTestV1ClientBuildGetUserRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildGetUserRequest").Body(
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
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetUserRequest").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
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
						utils.FormatString("%d",
							jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
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

func buildTestV1ClientGetUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_GetUser").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Comment("the hashed password is never transmitted over the wire."),
				jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword").Equals().EmptyString(),
				jen.Comment("the two factor secret is transmitted over the wire only on creation."),
				jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret").Equals().EmptyString(),
				jen.Comment("the two factor secret validation is never transmitted over the wire."),
				jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertTrue(
						jen.Qual("strings", "HasSuffix").Call(
							jen.ID(constants.RequestVarName).Dot("URL").Dot("String").Call(),
							jen.Qual("strconv", "Itoa").Call(
								jen.Int().Call(
									jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
								),
							),
						),
						nil,
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						utils.FormatString("/users/%d", jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodGet"), nil),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("User"))),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("Salt").Equals().Nil(),
				jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword").Equals().EmptyString(),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1ClientBuildGetUsersRequest() []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildGetUsersRequest").Body(
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
					constants.CtxVar(),
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
	}

	return lines
}

func buildTestV1ClientGetUsers(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_GetUsers").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "UserList"),
				jen.Comment("the hashed password is never transmitted over the wire."),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Zero()).Dot("HashedPassword").Equals().EmptyString(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.One()).Dot("HashedPassword").Equals().EmptyString(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Lit(2)).Dot("HashedPassword").Equals().EmptyString(),
				jen.Comment("the two factor secret is transmitted over the wire only on creation."),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Zero()).Dot("TwoFactorSecret").Equals().EmptyString(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.One()).Dot("TwoFactorSecret").Equals().EmptyString(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Lit(2)).Dot("TwoFactorSecret").Equals().EmptyString(),
				jen.Comment("the two factor secret validation is never transmitted over the wire."),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Zero()).Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.One()).Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
				jen.ID(utils.BuildFakeVarName("UserList")).Dot("Users").Index(jen.Lit(2)).Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.Lit("/users"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("UserList"))),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					constants.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("UserList")), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					constants.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1ClientBuildCreateUserRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildCreateUserRequest").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserCreationInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildCreateUserRequest").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
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

func buildTestV1ClientCreateUser(proj *models.Project) []jen.Code {
	fmp := proj.FakeModelsPackage()

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_CreateUser").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(fmp, "BuildFakeUserCreationInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.ID("expected").Assign().Qual(fmp, "BuildDatabaseCreationResponse").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.Lit("/users"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Var().ID("x").PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput"),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewDecoder").Call(jen.ID(constants.RequestVarName).Dot("Body")).Dot("Decode").Call(jen.AddressOf().ID("x")),
						nil,
					),
					utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Input")), jen.ID("x"), nil),
					jen.Line(),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID("expected")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(fmp, "BuildFakeUserCreationInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1ClientBuildArchiveUserRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildArchiveUserRequest").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodDelete"),
				utils.BuildFakeVar(proj, "User"),
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.RequireNotNil(jen.ID("actual").Dot("URL"), nil),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						utils.FormatString("%d",
							jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
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

func buildTestV1ClientArchiveUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_ArchiveUser").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						utils.FormatString("/users/%d", jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodDelete"),
						nil,
					),
				),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot("ArchiveUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("ArchiveUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1ClientBuildLoginRequest(proj *models.Project) []jen.Code {
	fmp := proj.FakeModelsPackage()

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildLoginRequest").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(fmp, "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot("BuildLoginRequest").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
				),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				utils.AssertEqual(
					jen.ID(constants.RequestVarName).Dot("Method"),
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
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot("BuildLoginRequest").Call(constants.CtxVar(), jen.Nil()),
				utils.AssertNil(jen.ID(constants.RequestVarName), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1ClientLogin(proj *models.Project) []jen.Code {
	fmp := proj.FakeModelsPackage()

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_Login").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			jen.Const().ID("expectedPath").Equals().Lit("/users/login"),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(fmp, "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.ID("expectedPath"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Qual("net/http", "SetCookie").Call(
						jen.ID(constants.ResponseVarName),
						jen.AddressOf().Qual("net/http", "Cookie").Values(
							jen.ID("Name").MapAssign().ID("exampleUser").Dot("Username"),
						),
					),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("c").Dot("Login").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
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
					constants.CtxVar(),
					jen.Nil(),
				),
				utils.AssertNil(jen.ID("cookie"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(fmp, "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("c").Dot("Login").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
				),
				utils.AssertNil(jen.ID("cookie"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(fmp, "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.ID("expectedPath"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
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
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(fmp, "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.ID("expectedPath"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
				),
				utils.RequireNil(jen.ID("cookie"), nil),
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

func buildTestV1ClientBuildValidateTOTPSecretRequest(proj *models.Project) []jen.Code {
	fmp := proj.FakeModelsPackage()

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_BuildValidateTOTPSecretRequest").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(fmp, "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(fmp, "BuildFakeTOTPSecretValidationInputForUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildVerifyTOTPSecretRequest").Call(
					constants.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID("exampleInput").Dot("TOTPToken"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodPost"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestV1ClientValidateTOTPSecret(proj *models.Project) []jen.Code {
	fmp := proj.FakeModelsPackage()

	lines := []jen.Code{
		utils.OuterTestFunc("V1Client_ValidateTOTPSecret").Body(
			utils.ParallelTest(nil),
			jen.Line(),
			jen.Const().ID("expectedPath").Equals().Lit("/users/totp_secret/verify"),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleUser").Assign().Qual(fmp, "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(fmp, "BuildFakeTOTPSecretValidationInputForUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.ID("expectedPath"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("VerifyTOTPSecret").Call(
					constants.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID("exampleInput").Dot("TOTPToken"),
				),
				utils.AssertNoError(jen.Err(), nil),
			), jen.Line(),
			utils.BuildSubTest(
				"with bad request response",
				jen.ID("exampleUser").Assign().Qual(fmp, "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(fmp, "BuildFakeTOTPSecretValidationInputForUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.ID("expectedPath"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("VerifyTOTPSecret").Call(
					constants.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID("exampleInput").Dot("TOTPToken"),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("ErrInvalidTOTPToken"), jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with otherwise invalid status code response",
				jen.ID("exampleUser").Assign().Qual(fmp, "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(fmp, "BuildFakeTOTPSecretValidationInputForUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.ID("expectedPath"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("VerifyTOTPSecret").Call(
					constants.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID("exampleInput").Dot("TOTPToken"),
				),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.ID("exampleUser").Assign().Qual(fmp, "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(fmp, "BuildFakeTOTPSecretValidationInputForUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("VerifyTOTPSecret").Call(
					constants.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID("exampleInput").Dot("TOTPToken"),
				),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with timeout",
				jen.ID("exampleUser").Assign().Qual(fmp, "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(fmp, "BuildFakeTOTPSecretValidationInputForUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.ID("expectedPath"),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Qual("time", "Sleep").Call(jen.Lit(10).Times().Qual("time", "Minute")),
					jen.Line(),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusAccepted")),
				),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.ID("c").Dot("plainClient").Dot("Timeout").Equals().Qual("time", "Millisecond"),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("VerifyTOTPSecret").Call(
					constants.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID("exampleInput").Dot("TOTPToken"),
				),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
