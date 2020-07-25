package mock

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataServerDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	sn := typ.Name.Singular()

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataServer", sn)).Equals().Parens(jen.PointerTo().IDf("%sDataServer", sn)).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Commentf("%sDataServer is a mocked models.%sDataServer for testing.", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataServer", sn).Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreationInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	if proj.SearchEnabled() {
		code.Add(
			jen.Comment("SearchHandler implements our interface requirements."),
			jen.Line(),
			jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("SearchHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			).Params().Block(
				jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			),
			jen.Line(),
		)
	}

	code.Add(
		jen.Comment("ListHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ListHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("CreateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ExistenceHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ExistenceHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ReadHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("UpdateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ArchiveHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Block(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	)

	return code
}
