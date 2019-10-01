package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersDotGo() *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(ret)

	ret.Add(jen.Const().ID("usersBasePath").Op("=").Lit("users"))

	ret.Add(
		jen.Comment("BuildGetUserRequest builds an HTTP request for fetching a user"),
		jen.Line(),
		newClientMethod("BuildGetUserRequest").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("userID"),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUser retrieves a user"),
		jen.Line(),
		newClientMethod("GetUser").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(
			jen.ID("user").Op("*").Qual(utils.ModelsPkg, "User"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetUserRequest").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("user"),
			),
			jen.Return().List(
				jen.ID("user"),
				jen.ID("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildGetUsersRequest builds an HTTP request for fetching a user"),
		jen.Line(),
		newClientMethod("BuildGetUsersRequest").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("usersBasePath"),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUsers retrieves a list of users"),
		jen.Line(),
		newClientMethod("GetUsers").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual(utils.ModelsPkg, "UserList"),
			jen.ID("error"),
		).Block(
			jen.ID("users").Op(":=").Op("&").Qual(utils.ModelsPkg, "UserList").Values(),
			jen.Line(),
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetUsersRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("users"),
			),
			jen.Return().List(jen.ID("users"),
				jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildCreateUserRequest builds an HTTP request for creating a user"),
		jen.Line(),
		newClientMethod("BuildCreateUserRequest").Params(
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(utils.ModelsPkg, "UserInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("nil"),
				jen.ID("usersBasePath"),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateUser creates a new user"),
		jen.Line(),
		newClientMethod("CreateUser").Params(
			utils.CtxParam(),
			jen.ID("input").Op("*").Qual(utils.ModelsPkg, "UserInput"),
		).Params(
			jen.Op("*").Qual(utils.ModelsPkg, "UserCreationResponse"),
			jen.ID("error"),
		).Block(
			jen.ID("user").Op(":=").Op("&").Qual(utils.ModelsPkg, "UserCreationResponse").Values(),
			jen.Line(),
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildCreateUserRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.ID("err").Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("user"),
			),
			jen.Return().List(jen.ID("user"),
				jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildArchiveUserRequest builds an HTTP request for updating a user"),
		jen.Line(),
		newClientMethod("BuildArchiveUserRequest").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("userID"),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveUser archives a user"),
		jen.Line(),
		newClientMethod("ArchiveUser").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildArchiveUserRequest").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildLoginRequest builds an authenticating HTTP request"),
		jen.Line(),
		newClientMethod("BuildLoginRequest").Params(
			jen.List(
				jen.ID("username"),
				jen.ID("password"),
				jen.ID("totpToken"),
			).ID("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("body"),
				jen.ID("err"),
			).Op(":=").ID("createBodyFromStruct").Call(
				jen.Op("&").Qual(utils.ModelsPkg, "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("username"),
					jen.ID("Password").Op(":").ID("password"),
					jen.ID("TOTPToken").Op(":").ID("totpToken"),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("creating body from struct: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("login"),
			),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Login will, when provided the correct credentials, fetch a login cookie"),
		jen.Line(),
		newClientMethod("Login").Params(
			utils.CtxParam(),
			jen.List(
				jen.ID("username"),
				jen.ID("password"),
				jen.ID("totpToken"),
			).ID("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Cookie"),
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildLoginRequest").Call(
				jen.ID("username"),
				jen.ID("password"),
				jen.ID("totpToken"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("plainClient").Dot("Do").Call(
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("encountered error executing login request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.ID("c").Dot("Debug")).Block(
				jen.List(
					jen.ID("b"),
					jen.ID("err"),
				).Op(":=").Qual("net/http/httputil", "DumpResponse").Call(
					jen.ID("res"),
					jen.ID("true"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("c").Dot("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("dumping response"),
					),
				),
				jen.ID("c").Dot("logger").Dot("WithValue").Call(
					jen.Lit("response"),
					jen.ID("string").Call(
						jen.ID("b"),
					),
				).Dot("Debug").Call(
					jen.Lit("login response received"),
				),
			),
			jen.Line(),
			jen.Defer().Func().Params().Block(
				jen.If(jen.ID("err").Op(":=").ID("res").Dot("Body").Dot("Close").Call(),
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.ID("c").Dot("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("closing response body"),
					),
				),
			).Call(),
			jen.Line(),
			jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
			jen.If(jen.ID("len").Call(
				jen.ID("cookies"),
			).Op(">").Lit(0),
			).Block(
				jen.Return().List(jen.ID("cookies").Index(
					jen.Lit(0),
				),
					jen.ID("nil"),
				),
			),
			jen.Line(),
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual("errors", "New").Call(
					jen.Lit("no cookies returned from request"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
