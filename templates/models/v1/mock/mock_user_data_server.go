package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockUserDataServerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "UserDataServer").Equals().Parens(jen.PointerTo().ID("UserDataServer")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserDataServer is a mocked models.UserDataServer for testing"),
		jen.Line(),
		jen.Type().ID("UserDataServer").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserLoginInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("UserLoginInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PasswordUpdateInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TOTPSecretVerificationInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("TOTPSecretVerificationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TOTPSecretRefreshInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ListHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TOTPSecretVerificationHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("TOTPSecretVerificationHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewTOTPSecretHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("NewTOTPSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdatePasswordHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("UpdatePasswordHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	return code
}
