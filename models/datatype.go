package models

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

// DataType represents a data model
type DataType struct {
	Name             wordsmith.SuperPalabra
	BelongsToUser    bool
	BelongsToNobody  bool
	RestrictedToUser bool
	BelongsToStruct  wordsmith.SuperPalabra
	Fields           []DataField
}

// DataField represents a data model's field
type DataField struct {
	Name                  wordsmith.SuperPalabra
	Type                  string
	Pointer               bool
	DefaultValue          string
	ValidForCreationInput bool
	ValidForUpdateInput   bool
}

// CtxParam is a shorthand for a context param
func ctxParam() jen.Code {
	return ctxVar().Qual("context", "Context")
}

// CtxParam is a shorthand for a context param
func ctxVar() *jen.Statement {
	return jen.ID("ctx")
}

func (typ DataType) BuildGetSomethingParams(proj *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToUser && typ.RestrictedToUser {
		lp = append(lp, jen.ID("userID"))
	}

	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}

func (typ DataType) BuildGetSomethingArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID("userID"))
	}

	return params
}

func (typ DataType) BuildGetSomethingArgsWithExampleVariables(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for i, pt := range owners {
		pts := pt.Name.Singular()
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			params = append(params, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", pts))
		} else {
			params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
		}
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	return params
}

func (typ DataType) BuildGetSomethingLogValues(proj *Project) jen.Code {
	params := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.Litf("%s_id", pt.Name.RouteName()).Op(":").IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.Litf("%s_id", typ.Name.RouteName()).Op(":").IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.Lit("user_id").Op(":").ID("userID"))
	}

	return jen.Map(jen.ID("string")).Interface().Valuesln(params...)
}

func (typ DataType) BuildGetListOfSomethingLogValues(proj *Project) jen.Code {
	params := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.Litf("%s_id", pt.Name.RouteName()).Op(":").IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.Lit("user_id").Op(":").ID("userID"))
	}

	return jen.Map(jen.ID("string")).Interface().Valuesln(params...)
}

func (typ DataType) BuildSomethingExistsArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID("userID"))
	}

	return params
}

func (typ DataType) BuildGetListOfSomethingParams(proj *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		lp = append(lp, jen.ID("userID"))
	}
	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	if !isModelsPackage {
		params = append(params, jen.ID("filter").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter"))
	} else {
		params = append(params, jen.ID("filter").Op("*").ID("QueryFilter"))
	}

	return params
}

const creationObjectVarName = "input"

func (typ DataType) BuildCreateSomethingParams(proj *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(creationObjectVarName).Op("*").IDf("%sCreationInput", sn))
	} else {
		params = append(params, jen.ID(creationObjectVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models", "v1"), fmt.Sprintf("%sCreationInput", sn)))
	}

	return params
}

func (typ DataType) BuildCreateSomethingQueryParams(proj *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(creationObjectVarName).Op("*").ID(sn))
	} else {
		params = append(params, jen.ID(creationObjectVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn))
	}

	return params
}

func (typ DataType) BuildCreateSomethingArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}

	params = append(params, lp...)
	params = append(params, jen.ID(creationObjectVarName))

	return params
}

func (typ DataType) BuildUpdateSomethingParams(proj *Project, updatedVarName string, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	if updatedVarName == "" {
		panic("BuildUpdateSomethingParams called with empty updatedVarName")
	}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if len(lp) > 1 {
		params = append(params, jen.List(lp[:len(lp)-1]...).ID("uint64"))
	}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(updatedVarName).Op("*").ID(sn))
	} else {
		params = append(params, jen.ID(updatedVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn))
	}

	return params
}

func (typ DataType) BuildUpdateSomethingArgs(proj *Project, updatedVarName string) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if len(lp) >= 1 {
		params = append(params, jen.List(lp[:len(lp)-1]...))
	}
	params = append(params, jen.ID(updatedVarName))

	return params
}

func (typ DataType) BuildGetSomethingForSomethingElseParams(proj *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.List(lp...).ID("uint64"))
	params = append(params, jen.ID("filter").Op("*").Qual(filepath.Join(proj.OutputPath, "models", "v1"), "QueryFilter"))

	return params
}

func (typ DataType) BuildGetSomethingForSomethingElseParamsForModelsPackage(proj *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.List(lp...).ID("uint64"))
	params = append(params, jen.ID("filter").Op("*").ID("QueryFilter"))

	return params
}

func (typ DataType) BuildGetSomethingForSomethingElseArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func (typ DataType) BuildGetListOfSomethingArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID("userID"))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func (typ DataType) BuildGetSomethingForUserParams(proj *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	lp = append(lp, jen.ID("userID"))
	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}

func buildFakeVarName(typName string) string {
	return fmt.Sprintf("example%s", typName)
}

func (typ DataType) BuildVarDeclarationsOfDependentStructs(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		pts := pt.Name.Singular()
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) BuildVarDeclarationsOfDependentStructsForListFunction(proj *Project) []jen.Code {
	lines := []jen.Code{}
	//sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}
	//lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) BuildVarDeclarationsOfDependentStructsForUpdateFunction(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()
	owners := proj.FindOwnerTypeChain(typ)

	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			pts := pt.Name.Singular()
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) BuildCreationVarDeclarationsOfDependentStructs(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines,
			jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call(),
		)
	}

	lines = append(lines,
		jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call(),
	)

	return lines
}

func (typ DataType) BuildCreationVarDeclarationsOfDependentStructsSkippingPossibleOwnerStruct(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lines = append(lines,
				jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call(),
			)
		}
	}

	lines = append(lines,
		jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call(),
	)

	return lines
}

func (typ DataType) BuildFormatStringForSingleInstanceRoute(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForListRoute(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s", typ.Name.PluralRouteName())

	return modelRoute
}

func (typ DataType) BuildFormatStringForSingleInstanceCreationRoute(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += typ.Name.PluralRouteName()

	return modelRoute
}

func (typ DataType) BuildFormatCallArgsForSingleInstanceRoute(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", pt.Name.Singular()))
		} else {
			callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}
	callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForListRoute(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		//if i == len(owners)-1 && typ.BelongsToStruct != nil {
		//	callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		//} else {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		//}
	}

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForSingleInstanceCreationRoute(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForSingleInstanceRouteThatIncludesItsOwnType(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)

	for i, pt := range owners {
		if typ.BelongsToStruct != nil && i == len(owners)-1 {
			callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}
	callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	return callArgs
}

func (typ DataType) BuildParamsForMethodThatHandlesAnInstanceWithIDs(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	} else {
		if !call {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64())
		} else {
			params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		}
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatHandlesAnInstanceWithStructs(proj *Project) []jen.Code {
	owners := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(owners) > 0 {
		for i, pt := range owners {
			if i == len(owners)-1 && typ.BelongsToStruct != nil {
				listParams = append(listParams, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", pt.Name.Singular()))
			} else {
				listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
			}
		}
		listParams = append(listParams, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatRetrievesAListOfADataType(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))
	} else {
		params = append(params, jen.ID(constants.FilterVarName))
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatRetrievesAListOfADataTypeFromStructs(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				listParams = append(listParams, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
			} else {
				listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
			}
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))
	} else {
		params = append(params, jen.ID(constants.FilterVarName))
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatCreatesADataType(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))
	} else {
		params = append(params, jen.ID("input"))
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatCreatesADataTypeFromStructs(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				listParams = append(listParams, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
			} else {
				listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
			}
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(buildFakeVarName("Input")).PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))
	} else {
		params = append(params, jen.ID(buildFakeVarName("Input")))
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatFetchesAListOfDataTypes(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))
	} else {
		params = append(params, jen.ID(constants.FilterVarName))
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatFetchesAListOfDataTypesFromStructs(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				listParams = append(listParams, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
			} else {
				listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
			}
		}
		if !call {
			params = append(params, jen.List(listParams...).Uint64())
		} else {
			params = append(params, listParams...)
		}
	}

	if !call {
		params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))
	} else {
		params = append(params, jen.ID(constants.FilterVarName))
	}

	return params
}

func (typ DataType) BuildCallForMethodThatFetchesAListOfDataTypesFromStructsForListFunction(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			//if i == len(parents)-1 && typ.BelongsToStruct != nil {
			//	listParams = append(listParams, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
			//} else {
			listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
			//}
		}
		params = append(params, listParams...)
	}
	params = append(params, jen.ID(constants.FilterVarName))

	return params
}

func (typ DataType) BuildParamsForMethodThatIncludesItsOwnTypeInItsParams(proj *Project, call bool) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		func() jen.Code {
			if call {
				return constants.CtxVar()
			} else {
				return constants.CtxParam()
			}
		}(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			if !call {
				params = append(params, jen.List(listParams...).Uint64())
			} else {
				params = append(params, listParams...)
			}
		}
	}

	if !call {
		params = append(params, jen.ID(typ.Name.UnexportedVarName()).PointerTo().Qual(proj.ModelsV1Package(), typ.Name.Singular()))
	} else {
		params = append(params, jen.ID(typ.Name.UnexportedVarName()))
	}

	return params
}

func (typ DataType) BuildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, listParams...)
		}
	}
	params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())))

	return params
}

func (typ DataType) BuildRequisiteFakeVarDecs(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	return lines
}

func (typ DataType) BuildRequisiteFakeVarDecsForListFunction(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
	}
	//if len(owners) == 0 {
	//	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	//}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgs(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser && typ.RestrictedToUser {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return lines
}

func (typ DataType) BuildRequisiteIDCallArgsForCreation(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}

	return lines
}

func (typ DataType) BuildRequisiteIDCallArgsForListFunction(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.BelongsToUser && typ.RestrictedToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteIDDeclarationsForUpdateFunction(proj *Project) []jen.Code {
	lines := []jen.Code{constants.CreateCtx()}

	owners := proj.FindOwnerTypeChain(typ)

	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	return lines
}
func (typ DataType) BuildRequisiteIDCallArgsForUpdateFunction(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}

	return lines
}

func (typ DataType) BuildRequisiteIDCallArgsWithPreCreatedUser(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}
