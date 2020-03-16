package models

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

// DataType represents a data model
type DataType struct {
	Name                 wordsmith.SuperPalabra
	BelongsToUser        bool
	ReadRestrictedToUser bool
	BelongsToNobody      bool
	BelongsToStruct      wordsmith.SuperPalabra
	Fields               []DataField
}

// DataField represents a data model's field
type DataField struct {
	Name                  wordsmith.SuperPalabra
	Type                  string
	Pointer               bool
	DefaultAllowed        bool
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

func (typ DataType) BuildGetSomethingParams(pkg *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToUser {
		lp = append(lp, jen.ID("userID"))
	}

	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}

func (typ DataType) BuildGetSomethingArgs(pkg *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToUser {
		params = append(params, jen.ID("userID"))
	}

	return params
}

func (typ DataType) BuildGetListOfSomethingParams(pkg *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.BelongsToUser {
		lp = append(lp, jen.ID("userID"))
	}
	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	if !isModelsPackage {
		params = append(params, jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"))
	} else {
		params = append(params, jen.ID("filter").Op("*").ID("QueryFilter"))
	}

	return params
}

const creationObjectVarName = "input"

func (typ DataType) BuildCreateSomethingParams(pkg *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(creationObjectVarName).Op("*").IDf("%sCreationInput", sn))
	} else {
		params = append(params, jen.ID(creationObjectVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models", "v1"), fmt.Sprintf("%sCreationInput", sn)))
	}

	return params
}

func (typ DataType) BuildCreateSomethingArgs(pkg *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}

	params = append(params, lp...)
	params = append(params, jen.ID(creationObjectVarName))

	return params
}

func (typ DataType) BuildUpdateSomethingParams(pkg *Project, updatedVarName string, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	if updatedVarName == "" {
		panic("BuildUpdateSomethingParams called with empty updatedVarName")
	}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if len(lp) > 1 {
		params = append(params, jen.List(lp[:len(lp)-1]...).ID("uint64"))
	}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(updatedVarName).Op("*").ID(sn))
	} else {
		params = append(params, jen.ID(updatedVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models", "v1"), sn))
	}

	return params
}

func (typ DataType) BuildUpdateSomethingArgs(pkg *Project, updatedVarName string) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if len(lp) >= 1 {
		params = append(params, jen.List(lp[:len(lp)-1]...))
	}
	params = append(params, jen.ID(updatedVarName))

	return params
}

func (typ DataType) BuildGetSomethingForSomethingElseParams(pkg *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.List(lp...).ID("uint64"))
	params = append(params, jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models", "v1"), "QueryFilter"))

	return params
}

func (typ DataType) BuildGetSomethingForSomethingElseParamsForModelsPackage(pkg *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.List(lp...).ID("uint64"))
	params = append(params, jen.ID("filter").Op("*").ID("QueryFilter"))

	return params
}

func (typ DataType) BuildGetSomethingForSomethingElseArgs(pkg *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func (typ DataType) BuildGetListOfSomethingArgs(pkg *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.BelongsToUser {
		params = append(params, jen.ID("userID"))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func (typ DataType) BuildGetSomethingForUserParams(pkg *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	for _, pt := range pkg.FindOwnerTypeChain(typ) {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	lp = append(lp, jen.ID("userID"))
	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}
