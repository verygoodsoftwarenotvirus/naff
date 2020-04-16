package mock

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataManagerDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	n := typ.Name
	sn := n.Singular()

	ret.Add(
		jen.Var().Underscore().ID("models").Dotf("%sDataManager", sn).Equals().Parens(jen.PointerTo().IDf("%sDataManager", sn)).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("%sDataManager is a mocked models.%sDataManager for testing", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataManager", sn).Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(buildSomethingExists(proj, typ)...)
	ret.Add(buildGetSomething(proj, typ)...)
	ret.Add(buildGetAllSomethingsCount(typ)...)
	ret.Add(buildGetListOfSomething(proj, typ)...)
	ret.Add(buildCreateSomething(proj, typ)...)
	ret.Add(buildUpdateSomething(proj, typ)...)
	ret.Add(buildArchiveSomething(proj, typ)...)

	return ret
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	funcName := fmt.Sprintf("%sExists", sn)

	params := typ.BuildInterfaceDefinitionExistenceMethodParams(proj)
	callArgs := typ.BuildInterfaceDefinitionExistenceMethodCallArgs(proj)

	lines := []jen.Code{
		jen.Commentf("%s is a mock function", funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).ID(funcName).Params(params...).Params(jen.Bool(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Zero()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildInterfaceDefinitionRetrievalMethodParams(proj)
	callArgs := typ.BuildInterfaceDefinitionRetrievalMethodCallArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(params...).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllSomethingsCount(typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	lines := []jen.Code{
		jen.Commentf("GetAll%sCount is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("GetAll%sCount", pn).Params(
			constants.CtxParam(),
		).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := typ.BuildMockDataManagerListRetrievalMethodParams(proj)
	callArgs := typ.BuildMockDataManagerListRetrievalMethodCallArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%s", pn).Params(params...).Params(jen.PointerTo().ID("models").Dotf("%sList", sn),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().ID("models").Dotf("%sList", sn)), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildMockInterfaceDefinitionCreationMethodParams(proj)
	args := typ.BuildMockInterfaceDefinitionCreationMethodCallArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Create%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(
				args...,
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildMockDataManagerUpdateMethodParams(proj, "updated")
	args := typ.BuildMockDataManagerUpdateMethodCallArgs(proj, "updated")

	lines := []jen.Code{
		jen.Commentf("Update%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(
				args...,
			).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildInterfaceDefinitionArchiveMethodParams(proj)
	callArgs := typ.BuildInterfaceDefinitionArchiveMethodCallArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Archive%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Archive%s", sn).Params(params...).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(callArgs...).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
