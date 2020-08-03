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

	code.Add(buildIterableDataServer(typ)...)
	code.Add(buildIterableCreationInputMiddleware(typ)...)
	code.Add(buildIterableUpdateInputMiddleware(typ)...)
	if proj.SearchEnabled() {
		code.Add(buildIterableSearchHandler(typ)...)
	}
	code.Add(buildIterableListHandler(typ)...)
	code.Add(buildIterableCreateHandler(typ)...)
	code.Add(buildIterableExistenceHandler(typ)...)
	code.Add(buildIterableReadHandler(typ)...)
	code.Add(buildIterableUpdateHandler(typ)...)
	code.Add(buildIterableArchiveHandler(typ)...)

	return code
}

func buildIterableDataServer(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Commentf("%sDataServer is a mocked models.%sDataServer for testing.", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataServer", sn).Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildIterableCreationInputMiddleware(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("CreationInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableUpdateInputMiddleware(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("UpdateInputMiddleware implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableSearchHandler(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("SearchHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("SearchHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableListHandler(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("ListHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ListHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableCreateHandler(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("CreateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("CreateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableExistenceHandler(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("ExistenceHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ExistenceHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableReadHandler(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("ReadHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ReadHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableUpdateHandler(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("UpdateHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("UpdateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildIterableArchiveHandler(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Comment("ArchiveHandler implements our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataServer", sn)).ID("ArchiveHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}
