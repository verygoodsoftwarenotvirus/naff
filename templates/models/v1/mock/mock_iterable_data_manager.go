package mock

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataManagerDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	n := typ.Name
	sn := n.Singular()

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", sn)).Equals().Parens(jen.PointerTo().IDf("%sDataManager", sn)).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Commentf("%sDataManager is a mocked models.%sDataManager for testing.", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataManager", sn).Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	)

	code.Add(buildSomethingExists(proj, typ)...)
	code.Add(buildGetSomething(proj, typ)...)
	code.Add(buildGetAllSomethingsCount(proj, typ)...)
	code.Add(buildGetAllSomethings(proj, typ)...)
	code.Add(buildGetListOfSomething(proj, typ)...)
	code.Add(buildGetSomethingsWithIDs(proj, typ)...)
	code.Add(buildCreateSomething(proj, typ)...)
	code.Add(buildUpdateSomething(proj, typ)...)
	code.Add(buildArchiveSomething(proj, typ)...)

	return code
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	funcName := fmt.Sprintf("%sExists", sn)

	params := typ.BuildInterfaceDefinitionExistenceMethodParams(proj)
	callArgs := typ.BuildInterfaceDefinitionExistenceMethodCallArgs(proj)

	lines := []jen.Code{
		jen.Commentf("%s is a mock function.", funcName),
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
		jen.Commentf("Get%s is a mock function.", sn),
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

func buildGetAllSomethingsCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	lines := []jen.Code{
		jen.Commentf("GetAll%sCount is a mock function.", pn),
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

func buildGetAllSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	lines := []jen.Code{
		jen.Commentf("GetAll%s is a mock function.", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("GetAll%s", pn).Params(
			constants.CtxParam(),
			jen.ID("results").Chan().Index().Qual(proj.ModelsV1Package(), sn),
		).Params(jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("results")),
			jen.Return().List(
				jen.ID("args").Dot("Error").Call(jen.Zero()),
			),
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
		jen.Commentf("Get%s is a mock function.", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%s", pn).Params(params...).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(
					jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)),
				),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomethingsWithIDs(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := []jen.Code{
		constants.CtxParam(),
		constants.UserIDParam(),
		jen.ID("limit").Uint8(),
		jen.ID("ids").Index().Uint64(),
	}
	callArgs := []jen.Code{
		constants.CtxVar(),
		constants.UserIDVar(),
		jen.ID("limit"),
		jen.ID("ids"),
	}

	lines := []jen.Code{
		jen.Commentf("Get%sWithIDs is a mock function.", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%sWithIDs", pn).
			Params(params...).
			Params(jen.Index().Qual(proj.ModelsV1Package(), sn), jen.Error()).
			Block(
				jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
				jen.Return().List(
					jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(
						jen.Index().Qual(proj.ModelsV1Package(), sn),
					),
					jen.ID("args").Dot("Error").Call(jen.One()),
				),
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
		jen.Commentf("Create%s is a mock function.", sn),
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
		jen.Commentf("Update%s is a mock function.", sn),
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
		jen.Commentf("Archive%s is a mock function.", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Archive%s", sn).Params(params...).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(callArgs...).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
