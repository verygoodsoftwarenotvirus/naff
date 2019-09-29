package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func usersDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(jen.Var().ID("usersBasePath").Op("=").Lit("users"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildGetUserRequest").Params(
			ctxParam(),
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
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("GetUser").Params(
			ctxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(
			jen.ID("user").Op("*").Qual(modelsPkg, "User"),
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
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildGetUsersRequest").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("usersBasePath"),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("GetUsers").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual(modelsPkg, "UserList"),
			jen.ID("error"),
		).Block(
			jen.ID("users").Op(":=").Op("&").Qual(modelsPkg, "UserList").Values(),
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetUsersRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
				),
			),
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
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildCreateUserRequest").Params(
			ctxParam(),
			jen.ID("body").Op("*").Qual(modelsPkg, "UserInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("nil"),
				jen.ID("usersBasePath"),
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
		jen.Comment(""),
		jen.Line(),
		newClientMethod("CreateUser").Params(
			ctxParam(),
			jen.ID("input").Op("*").Qual(modelsPkg, "UserInput"),
		).Params(
			jen.Op("*").Qual(modelsPkg, "UserCreationResponse"),
			jen.ID("error"),
		).Block(
			jen.ID("user").Op(":=").Op("&").Qual(modelsPkg, "UserCreationResponse").Values(),
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
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildArchiveUserRequest").Params(
			ctxParam(),
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
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("ArchiveUser").Params(
			ctxParam(),
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
					jen.Lit("building request"),
					jen.ID("err"),
				),
			),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
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
				jen.Op("&").Qual(modelsPkg, "UserLoginInput").Values(jen.Dict{
					jen.ID("Username"):  jen.ID("username"),
					jen.ID("Password"):  jen.ID("password"),
					jen.ID("TOTPToken"): jen.ID("totpToken"),
				}),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("creating body from struct"),
					jen.ID("err"),
				),
				),
			),
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
		jen.Comment(""),
		jen.Line(),
		newClientMethod("Login").Params(
			ctxParam(),
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
				jen.Return().List(jen.ID("nil"),
					jen.ID("err")),
			),
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("plainClient").Dot("Do").Call(
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("encountered error executing login request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.If(jen.ID("c").Dot("Debug")).Block(
				jen.List(
					jen.ID("b"),
					jen.ID("err"),
				).Op(":=").ID("httputil").Dot("DumpResponse").Call(
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
			jen.Return().List(
				jen.ID("nil"),
				jen.ID("errors").Dot("New").Call(
					jen.Lit("no cookies returned from request"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
