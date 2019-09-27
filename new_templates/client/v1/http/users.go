package client

import jen "github.com/dave/jennifer/jen"

func usersDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(jen.Var().Id("usersBasePath").Op("=").Lit("users"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildGetUserRequest").Params(
			ctxParam(),
			jen.Id("userID").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(
				jen.Id("nil"),
				jen.Id("usersBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("userID"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("GetUser").Params(
			ctxParam(),
			jen.Id("userID").Id("uint64"),
		).Params(
			jen.Id("user").Op("*").Id("models").Dot("User"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetUserRequest").Call(
				jen.Id("ctx"),
				jen.Id("userID"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("user"),
			),
			jen.Return().List(
				jen.Id("user"),
				jen.Id("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildGetUsersRequest").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(
				jen.Id("filter").Dot("ToValues").Call(),
				jen.Id("usersBasePath"),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("GetUsers").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Op("*").Id("models").Dot("UserList"),
			jen.Id("error"),
		).Block(
			jen.Id("users").Op(":=").Op("&").Id("models").Dot("UserList").Values(),
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetUsersRequest").Call(
				jen.Id("ctx"),
				jen.Id("filter"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Id("err"),
				),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("users"),
			),
			jen.Return().List(jen.Id("users"),
				jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildCreateUserRequest").Params(
			ctxParam(),
			jen.Id("body").Op("*").Id("models").Dot("UserInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(
				jen.Id("nil"),
				jen.Id("usersBasePath"),
			),
			jen.Return().Id("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.Id("uri"),
				jen.Id("body"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("CreateUser").Params(
			ctxParam(),
			jen.Id("input").Op("*").Id("models").Dot("UserInput"),
		).Params(
			jen.Op("*").Id("models").Dot("UserCreationResponse"),
			jen.Id("error"),
		).Block(
			jen.Id("user").Op(":=").Op("&").Id("models").Dot("UserCreationResponse").Values(),
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildCreateUserRequest").Call(
				jen.Id("ctx"),
				jen.Id("input"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("executeUnathenticatedDataRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("user"),
			),
			jen.Return().List(jen.Id("user"),
				jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildArchiveUserRequest").Params(
			ctxParam(),
			jen.Id("userID").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(
				jen.Id("nil"),
				jen.Id("usersBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("userID"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("ArchiveUser").Params(
			ctxParam(),
			jen.Id("userID").Id("uint64"),
		).Params(jen.Id("error")).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildArchiveUserRequest").Call(
				jen.Id("ctx"),
				jen.Id("userID"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request"),
					jen.Id("err"),
				),
			),
			jen.Return().Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildLoginRequest").Params(
			jen.List(
				jen.Id("username"),
				jen.Id("password"),
				jen.Id("totpToken"),
			).Id("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("body"),
				jen.Id("err"),
			).Op(":=").Id("createBodyFromStruct").Call(
				jen.Op("&").Id("models").Dot("UserLoginInput").Values(
					jen.Id("Username").Op(":").Id("username"),
					jen.Id("Password").Op(":").Id("password"),
					jen.Id("TOTPToken").Op(":").Id("totpToken"),
				),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("creating body from struct"),
					jen.Id("err"),
				),
				),
			),
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(
				jen.Id("nil"),
				jen.Id("usersBasePath"),
				jen.Lit("login"),
			),
			jen.Return().Id("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.Id("uri"),
				jen.Id("body"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("Login").Params(
			ctxParam(),
			jen.List(
				jen.Id("username"),
				jen.Id("password"),
				jen.Id("totpToken"),
			).Id("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Cookie"),
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildLoginRequest").Call(
				jen.Id("username"),
				jen.Id("password"),
				jen.Id("totpToken"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"),
					jen.Id("err")),
			),
			jen.List(
				jen.Id("res"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("plainClient").Dot("Do").Call(
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("encountered error executing login request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.If(jen.Id("c").Dot("Debug")).Block(
				jen.List(
					jen.Id("b"),
					jen.Id("err"),
				).Op(":=").Id("httputil").Dot("DumpResponse").Call(
					jen.Id("res"),
					jen.Id("true"),
				),
				jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
					jen.Id("c").Dot("logger").Dot("Error").Call(
						jen.Id("err"),
						jen.Lit("dumping response"),
					),
				),
				jen.Id("c").Dot("logger").Dot("WithValue").Call(
					jen.Lit("response"),
					jen.Id("string").Call(
						jen.Id("b"),
					),
				).Dot("Debug").Call(
					jen.Lit("login response received"),
				),
			),
			jen.Defer().Func().Params().Block(
				jen.If(jen.Id("err").Op(":=").Id("res").Dot("Body").Dot("Close").Call(),
					jen.Id("err").Op("!=").Id("nil"),
				).Block(
					jen.Id("c").Dot("logger").Dot("Error").Call(
						jen.Id("err"),
						jen.Lit("closing response body"),
					),
				),
			).Call(),
			jen.Id("cookies").Op(":=").Id("res").Dot("Cookies").Call(),
			jen.If(jen.Id("len").Call(
				jen.Id("cookies"),
			).Op(">").Lit(0),
			).Block(
				jen.Return().List(jen.Id("cookies").Index(
					jen.Lit(0),
				),
					jen.Id("nil"),
				),
			),
			jen.Return().List(
				jen.Id("nil"),
				jen.Id("errors").Dot("New").Call(
					jen.Lit("no cookies returned from request"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
