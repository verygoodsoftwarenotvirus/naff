package mock

import (
	"fmt"

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
		jen.Var().ID("_").ID("models").Dotf("%sDataManager", sn).Equals().Parens(jen.PointerTo().IDf("%sDataManager", sn)).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("%sDataManager is a mocked models.%sDataManager for testing", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataManager", sn).Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(buildSomethingExists(proj, typ)...)
	ret.Add(buildGetSomething(proj, typ)...)
	//ret.Add(buildGetSomethingCount(proj, typ)...)
	ret.Add(buildGetAllSomethingsCount(typ)...)
	ret.Add(buildGetListOfSomething(proj, typ)...)

	//if typ.BelongsToUser {
	//	ret.Add(buildGetAllSomethingsForUser(proj, typ)...)
	//}
	if typ.BelongsToStruct != nil {
		ret.Add(buildGetAllSomethingsForSomethingElse(proj, typ)...)
	}

	ret.Add(buildCreateSomething(proj, typ)...)
	ret.Add(buildUpdateSomething(proj, typ)...)
	ret.Add(buildArchiveSomething(proj, typ)...)

	return ret
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	funcName := fmt.Sprintf("%sExists", sn)

	params := typ.BuildGetSomethingParams(proj)
	callArgs := typ.BuildGetSomethingArgs(proj)

	lines := []jen.Code{
		jen.Commentf("%s is a mock function", funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).ID(funcName).Params(params...).Params(jen.Bool(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildGetSomethingParams(proj)
	callArgs := typ.BuildGetSomethingArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(params...).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildGetListOfSomethingParams(proj, false)
	callArgs := typ.BuildGetListOfSomethingArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Get%sCount is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%sCount", sn).Params(params...).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
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
			utils.CtxParam(),
		).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := typ.BuildGetListOfSomethingParams(proj, false)
	callArgs := typ.BuildGetListOfSomethingArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%s", pn).Params(params...).Params(jen.PointerTo().ID("models").Dotf("%sList", sn),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().ID("models").Dotf("%sList", sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllSomethingsForUser(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := typ.BuildGetSomethingForUserParams(proj)

	lines := []jen.Code{
		jen.Commentf("GetAll%sForUser is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("GetAll%sForUser", pn).Params(
			params...,
		).Params(jen.Index().Qual(proj.ModelsV1Package(), sn),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Qual(proj.ModelsV1Package(), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllSomethingsForSomethingElse(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := typ.BuildGetSomethingForSomethingElseParams(proj)

	lines := []jen.Code{
		jen.Commentf("GetAll%sFor%s is a mock function", pn, typ.BelongsToStruct.Singular()),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular()).Params(
			params...,
		).Params(jen.Index().Qual(proj.ModelsV1Package(), sn),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Qual(proj.ModelsV1Package(), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildCreateSomethingParams(proj, false)
	args := typ.BuildCreateSomethingArgs(proj)

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
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildUpdateSomethingParams(proj, "updated", false)
	args := typ.BuildUpdateSomethingArgs(proj, "updated")

	lines := []jen.Code{
		jen.Commentf("Update%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(
				args...,
			).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildGetSomethingParams(proj)
	callArgs := typ.BuildGetSomethingArgs(proj)

	lines := []jen.Code{
		jen.Commentf("Archive%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Archive%s", sn).Params(params...).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(callArgs...).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	}

	return lines
}
