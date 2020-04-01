package mock

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataManagerDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg, ret)

	n := typ.Name
	sn := n.Singular()

	ret.Add(
		jen.Var().ID("_").ID("models").Dotf("%sDataManager", sn).Equals().Parens(jen.Op("*").IDf("%sDataManager", sn)).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("%sDataManager is a mocked models.%sDataManager for testing", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataManager", sn).Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(buildSomethingExists(pkg, typ)...)
	ret.Add(buildGetSomething(pkg, typ)...)
	ret.Add(buildGetSomethingCount(pkg, typ)...)
	ret.Add(buildGetAllSomethingsCount(typ)...)
	ret.Add(buildGetListOfSomething(pkg, typ)...)

	if typ.BelongsToUser {
		ret.Add(buildGetAllSomethingsForUser(pkg, typ)...)
	}
	if typ.BelongsToStruct != nil {
		ret.Add(buildGetAllSomethingsForSomethingElse(pkg, typ)...)
	}

	ret.Add(buildCreateSomething(pkg, typ)...)
	ret.Add(buildUpdateSomething(pkg, typ)...)
	ret.Add(buildArchiveSomething(pkg, typ)...)

	return ret
}

func buildSomethingExists(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	funcName := fmt.Sprintf("%sExists", sn)

	params := typ.BuildGetSomethingParams(pkg)
	callArgs := typ.BuildGetSomethingArgs(pkg)

	lines := []jen.Code{
		jen.Commentf("%s is a mock function", funcName),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).ID(funcName).Params(params...).Params(jen.Bool(), jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildGetSomethingParams(pkg)
	callArgs := typ.BuildGetSomethingArgs(pkg)

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(params...).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomethingCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildGetListOfSomethingParams(pkg, false)
	callArgs := typ.BuildGetListOfSomethingArgs(pkg)

	lines := []jen.Code{
		jen.Commentf("Get%sCount is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%sCount", sn).Params(params...).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
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
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sCount", pn).Params(
			utils.CtxVar().Qual("context", "Context"),
		).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetListOfSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := typ.BuildGetListOfSomethingParams(pkg, false)
	callArgs := typ.BuildGetListOfSomethingArgs(pkg)

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", pn).Params(params...).Params(jen.Op("*").ID("models").Dotf("%sList", sn),
			jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dotf("%sList", sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllSomethingsForUser(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := typ.BuildGetSomethingForUserParams(pkg)

	lines := []jen.Code{
		jen.Commentf("GetAll%sForUser is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sForUser", pn).Params(
			params...,
		).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllSomethingsForSomethingElse(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := typ.BuildGetSomethingForSomethingElseParams(pkg)

	lines := []jen.Code{
		jen.Commentf("GetAll%sFor%s is a mock function", pn, typ.BelongsToStruct.Singular()),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular()).Params(
			params...,
		).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildCreateSomethingParams(pkg, false)
	args := typ.BuildCreateSomethingArgs(pkg)

	lines := []jen.Code{
		jen.Commentf("Create%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(
				args...,
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildUpdateSomethingParams(pkg, "updated", false)
	args := typ.BuildUpdateSomethingArgs(pkg, "updated")

	lines := []jen.Code{
		jen.Commentf("Update%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(
				args...,
			).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildGetSomethingParams(pkg)
	callArgs := typ.BuildGetSomethingArgs(pkg)

	lines := []jen.Code{
		jen.Commentf("Archive%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Archive%s", sn).Params(params...).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(callArgs...).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	}

	return lines
}
