package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataManagerDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	n := typ.Name
	sn := n.Singular()

	ret.Add(
		jen.Var().ID("_").ID("models").Dotf("%sDataManager", sn).Op("=").Parens(jen.Op("*").IDf("%sDataManager", sn)).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("%sDataManager is a mocked models.%sDataManager for testing", sn, sn),
		jen.Line(),
		jen.Type().IDf("%sDataManager", sn).Struct(
			jen.Qual("github.com/stretchr/testify/mock", "Mock"),
		),
		jen.Line(),
	)

	ret.Add(buildGetSomething(pkg, typ)...)
	ret.Add(buildGetSomethingCount(pkg, typ)...)
	ret.Add(buildGetAllSomethingsCount(pkg, typ)...)
	ret.Add(buildGetListOfSomething(pkg, typ)...)

	if typ.BelongsToUser {
		ret.Add(buildGetAllSomethingsForUser(pkg, typ)...)
	} else if typ.BelongsToStruct != nil {
		ret.Add(buildGetAllSomethingsForSomethingElse(pkg, typ)...)
	}

	ret.Add(buildCreateSomething(pkg, typ)...)
	ret.Add(buildUpdateSomething(pkg, typ)...)
	ret.Add(buildArchiveSomething(pkg, typ)...)

	return ret
}

func buildGetSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	uvn := n.UnexportedVarName()

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
	}
	callArgs := []jen.Code{
		jen.ID("ctx"), jen.IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		params = append(params, jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64"))
		callArgs = append(callArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		params = append(params, jen.List(jen.IDf("%sID", uvn), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).ID("uint64"))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	} else {
		params = append(params, jen.IDf("%sID", uvn).ID("uint64"))
	}

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(params...).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomethingCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
		jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
	}
	callArgs := []jen.Code{
		jen.ID("ctx"), jen.ID("filter"),
	}

	if typ.BelongsToUser {
		params = append(params, jen.ID("userID").ID("uint64"))
		callArgs = append(callArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64"))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	lines := []jen.Code{
		jen.Commentf("Get%sCount is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%sCount", sn).Params(params...).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllSomethingsCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	lines := []jen.Code{
		jen.Commentf("GetAll%sCount is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
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

	params := []jen.Code{
		jen.ID("ctx").Qual("context", "Context"),
		jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
	}
	callArgs := []jen.Code{
		jen.ID("ctx"), jen.ID("filter"),
	}

	if typ.BelongsToUser {
		params = append(params, jen.ID("userID").ID("uint64"))
		callArgs = append(callArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64"))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	lines := []jen.Code{
		jen.Commentf("Get%s is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Get%s", pn).Params(params...).Params(jen.Op("*").ID("models").Dotf("%sList", sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(callArgs...),
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

	lines := []jen.Code{
		jen.Commentf("GetAll%sForUser is a mock function", pn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sForUser", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("userID")),
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

	lines := []jen.Code{
		jen.Commentf("GetAll%sFor%s is a mock function", pn, typ.BelongsToStruct.Singular()),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular()).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	lines := []jen.Code{
		jen.Commentf("Create%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Create%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dotf("%sCreationInput", sn)).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	lines := []jen.Code{
		jen.Commentf("Update%s is a mock function", sn),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").IDf("%sDataManager", sn)).IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("updated")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := []jen.Code{jen.ID("ctx").Qual("context", "Context")}
	callArgs := []jen.Code{
		jen.ID("ctx"), jen.IDf("%sID", typ.Name.UnexportedVarName()),
	}

	if typ.BelongsToUser {
		params = append(params, jen.List(jen.IDf("%sID", typ.Name.UnexportedVarName()), jen.ID("userID")).ID("uint64"))
		callArgs = append(callArgs, jen.ID("userID"))
	} else if typ.BelongsToStruct != nil {
		params = append(params, jen.List(jen.IDf("%sID", typ.Name.UnexportedVarName()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).ID("uint64"))
		callArgs = append(callArgs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	} else if typ.BelongsToNobody {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).ID("uint64"))
	}

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
