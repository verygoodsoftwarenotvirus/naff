package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuilder_BuildGetUserRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/users/%d"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildGetUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildGetUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildGetUsersRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildGetUsersRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildSearchForUsersByUsernameRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users/search"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleUsername").Op(":=").ID("fakes").Dot("BuildFakeUser").Call().Dot("Username"),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("q=%s"),
							jen.ID("exampleUsername"),
						),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildSearchForUsersByUsernameRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleUsername"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with empty username"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildSearchForUsersByUsernameRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(""),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildCreateUserRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildCreateUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildCreateUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildArchiveUserRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/users/%d"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodDelete"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildArchiveUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildArchiveUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildArbitraryImage builds an image with a bunch of colors in it."),
		jen.Line(),
		jen.Func().ID("buildArbitraryImage").Params(jen.ID("widthAndHeight").ID("int")).Params(jen.Qual("image", "Image")).Body(
			jen.ID("img").Op(":=").Qual("image", "NewRGBA").Call(jen.Qual("image", "Rectangle").Valuesln(jen.ID("Min").Op(":").Qual("image", "Point").Values(), jen.ID("Max").Op(":").Qual("image", "Point").Valuesln(jen.ID("X").Op(":").ID("widthAndHeight"), jen.ID("Y").Op(":").ID("widthAndHeight")))),
			jen.For(jen.ID("x").Op(":=").Lit(0), jen.ID("x").Op("<").ID("widthAndHeight"), jen.ID("x").Op("++")).Body(
				jen.For(jen.ID("y").Op(":=").Lit(0), jen.ID("y").Op("<").ID("widthAndHeight"), jen.ID("y").Op("++")).Body(
					jen.ID("img").Dot("Set").Call(
						jen.ID("x"),
						jen.ID("y"),
						jen.Qual("image/color", "RGBA").Valuesln(jen.ID("R").Op(":").ID("uint8").Call(jen.ID("x").Op("%").Qual("math", "MaxUint8")), jen.ID("G").Op(":").ID("uint8").Call(jen.ID("y").Op("%").Qual("math", "MaxUint8")), jen.ID("B").Op(":").ID("uint8").Call(jen.ID("x").Op("+").ID("y").Op("%").Qual("math", "MaxUint8")), jen.ID("A").Op(":").Qual("math", "MaxUint8")),
					))),
			jen.Return().ID("img"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildPNGBytes").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("i").Qual("image", "Image")).Params(jen.Index().ID("byte")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("b").Op(":=").ID("new").Call(jen.Qual("bytes", "Buffer")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.Qual("image/png", "Encode").Call(
					jen.ID("b"),
					jen.ID("i"),
				),
			),
			jen.Return().ID("b").Dot("Bytes").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildAvatarUploadRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users/avatar/upload"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard jpeg"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("avatar").Op(":=").ID("buildArbitraryImage").Call(jen.Lit(123)),
					jen.ID("avatarBytes").Op(":=").ID("buildPNGBytes").Call(
						jen.ID("t"),
						jen.ID("avatar"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAvatarUploadRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("avatarBytes"),
						jen.Lit("jpeg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard png"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("avatar").Op(":=").ID("buildArbitraryImage").Call(jen.Lit(123)),
					jen.ID("avatarBytes").Op(":=").ID("buildPNGBytes").Call(
						jen.ID("t"),
						jen.ID("avatar"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAvatarUploadRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("avatarBytes"),
						jen.Lit("png"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard gif"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("avatar").Op(":=").ID("buildArbitraryImage").Call(jen.Lit(123)),
					jen.ID("avatarBytes").Op(":=").ID("buildPNGBytes").Call(
						jen.ID("t"),
						jen.ID("avatar"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAvatarUploadRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("avatarBytes"),
						jen.Lit("gif"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with empty avatar"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAvatarUploadRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
						jen.Lit("jpeg"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid extension"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("avatar").Op(":=").ID("buildArbitraryImage").Call(jen.Lit(123)),
					jen.ID("avatarBytes").Op(":=").ID("buildPNGBytes").Call(
						jen.ID("t"),
						jen.ID("avatar"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAvatarUploadRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("avatarBytes"),
						jen.Lit(""),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("helper").Dot("builder").Op("=").ID("buildTestRequestBuilderWithInvalidURL").Call(),
					jen.ID("avatar").Op(":=").ID("buildArbitraryImage").Call(jen.Lit(123)),
					jen.ID("avatarBytes").Op(":=").ID("buildPNGBytes").Call(
						jen.ID("t"),
						jen.ID("avatar"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAvatarUploadRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("avatarBytes"),
						jen.Lit("png"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildGetAuditLogForUserRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/users/%d/audit"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildGetAuditLogForUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildGetAuditLogForUserRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
