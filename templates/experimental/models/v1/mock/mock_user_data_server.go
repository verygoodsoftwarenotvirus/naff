package mock

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockUserDataServerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"UserDataServer",
	).Op("=").Parens(jen.Op("*").ID("UserDataServer")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("UserDataServer").Struct(jen.ID("mock").Dot(
		"Mock",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UserLoginInputMiddleware is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("UserLoginInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UserInputMiddleware is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// PasswordUpdateInputMiddleware is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// TOTPSecretRefreshInputMiddleware is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("next")),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ListHandler is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateHandler is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ReadHandler is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// NewTOTPSecretHandler is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("NewTOTPSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdatePasswordHandler is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("UpdatePasswordHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveHandler is a mock method to satisfy our interface requirements").Params(jen.ID("m").Op("*").ID("UserDataServer")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(),
		jen.Return().ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
	),

		jen.Line(),
	)
	return ret
}
