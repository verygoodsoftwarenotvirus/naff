package mock

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockIterableDataManagerDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	n := typ.Name
	sn := n.Singular()

	code.Add(
		jen.Var().Underscore().Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)).Equals().Parens(jen.PointerTo().IDf("%sDataManager", sn)).Call(jen.Nil()),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("%sDataManager is a mocked types.%sDataManager for testing.", sn, sn),
		jen.Newline(),
		jen.Type().IDf("%sDataManager", sn).Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Newline(),
	)

	code.Add(buildSomethingExists(proj, typ)...)
	code.Add(buildGetSomething(proj, typ)...)
	code.Add(buildGetAllSomethingsCount(proj, typ)...)
	code.Add(buildGetAllSomethings(proj, typ)...)
	code.Add(buildGetListOfSomething(proj, typ)...)
	code.Add(buildGetSomethingsWithIDs(proj, typ)...)
	code.Add(buildCreateSomething(proj, typ)...)
	code.Add(buildUpdateSomething(proj, typ)...)
	code.Add(buildArchiveSomething(typ)...)
	code.Add(buildGetAuditLogEntriesForSomething(proj, typ)...)

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
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).ID(funcName).Params(params...).Params(jen.Bool(), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Zero()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Newline(),
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
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%s", sn).Params(params...).Params(jen.PointerTo().Qual(proj.TypesPackage(), sn),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), sn)), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Newline(),
	}

	return lines
}

func buildGetAllSomethingsCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	lines := []jen.Code{
		jen.Commentf("GetAll%sCount is a mock function.", pn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("GetAll%sCount", pn).Params(
			constants.CtxParam(),
		).Params(jen.Uint64(), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Newline(),
	}

	return lines
}

func buildGetAllSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	lines := []jen.Code{
		jen.Commentf("GetAll%s is a mock function.", pn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("GetAll%s", pn).Params(
			constants.CtxParam(),
			jen.ID("results").Chan().Index().PointerTo().Qual(proj.TypesPackage(), sn),
			jen.ID("bucketSize").Uint16(),
		).Params(jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("results"), jen.ID("bucketSize")),
			jen.Return().List(
				jen.ID("args").Dot("Error").Call(jen.Zero()),
			),
		),
		jen.Newline(),
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
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%s", pn).Params(params...).Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(
					jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)),
				),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildGetListOfSomethingFromIDsParams(typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxParam()}

	lp := []jen.Code{}
	if typ.BelongsToStruct != nil {
		lp = append(lp, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	if typ.BelongsToAccount {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	params = append(params,
		jen.ID("limit").Uint8(),
		jen.ID("ids").Index().Uint64(),
	)

	return params
}

func buildGetListOfSomethingFromIDsArgs(typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	lp := []jen.Code{}
	if typ.BelongsToStruct != nil {
		lp = append(lp, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	if typ.BelongsToAccount {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...))
	}

	params = append(params,
		jen.ID("limit"),
		jen.ID("ids"),
	)

	return params
}

func buildGetSomethingsWithIDs(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	params := buildGetListOfSomethingFromIDsParams(typ)
	callArgs := buildGetListOfSomethingFromIDsArgs(typ)

	lines := []jen.Code{
		jen.Commentf("Get%sWithIDs is a mock function.", pn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Get%sWithIDs", pn).
			Params(params...).
			Params(jen.Index().PointerTo().Qual(proj.TypesPackage(), sn), jen.Error()).
			Body(
				jen.ID("args").Assign().ID("m").Dot("Called").Call(callArgs...),
				jen.Return().List(
					jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(
						jen.Index().PointerTo().Qual(proj.TypesPackage(), sn),
					),
					jen.ID("args").Dot("Error").Call(jen.One()),
				),
			),
		jen.Newline(),
	}

	return lines
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildMockInterfaceDefinitionCreationMethodParams(proj)
	args := typ.BuildMockInterfaceDefinitionCreationMethodCallArgs()

	lines := []jen.Code{
		jen.Commentf("Create%s is a mock function.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.PointerTo().Qual(proj.TypesPackage(), sn),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(
				args...,
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), sn)), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Newline(),
	}

	return lines
}

func buildUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildMockDataManagerUpdateMethodParams(proj, "updated")
	args := typ.BuildMockDataManagerUpdateMethodCallArgs("updated")

	lines := []jen.Code{
		jen.Commentf("Update%s is a mock function.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(
				args...,
			).Dot("Error").Call(jen.Zero()),
		),
		jen.Newline(),
	}

	return lines
}

func buildArchiveSomething(typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	params := typ.BuildInterfaceDefinitionArchiveMethodParams()
	callArgs := typ.BuildInterfaceDefinitionArchiveMethodCallArgs()

	lines := []jen.Code{
		jen.Commentf("Archive%s is a mock function.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("Archive%s", sn).Params(params...).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(callArgs...).Dot("Error").Call(jen.Zero()),
		),
		jen.Newline(),
	}

	return lines
}

func buildGetAuditLogEntriesForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	lines := []jen.Code{
		jen.Commentf("GetAuditLogEntriesFor%s is a mock function.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("m").PointerTo().IDf("%sDataManager", sn)).IDf("GetAuditLogEntriesFor%s", sn).Params(constants.CtxParam(), jen.IDf("%sID", n.UnexportedVarName()).Uint64()).Params(jen.Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry"), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.IDf("%sID", n.UnexportedVarName())),
			jen.Return().List(
				jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry")),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Newline(),
	}

	return lines
}
