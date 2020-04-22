package models

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"path/filepath"
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

func buildFakeVarName(typName string) string {
	return fmt.Sprintf("example%s", typName)
}

// CtxParam is a shorthand for a context param
func ctxParam() jen.Code {
	return ctxVar().Qual("context", "Context")
}

// CtxParam is a shorthand for a context param
func ctxVar() *jen.Statement {
	return jen.ID("ctx")
}

func (typ DataType) RestrictedToUserAtSomeLevel(proj *Project) bool {
	for _, o := range proj.FindOwnerTypeChain(typ) {
		if o.BelongsToUser && o.RestrictedToUser {
			return true
		}
	}

	return typ.BelongsToUser && typ.RestrictedToUser
}

func (typ DataType) buildGetSomethingParams(proj *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		lp = append(lp, jen.ID("userID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	return params
}

func (typ DataType) buildArchiveSomethingParams(proj *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	if typ.BelongsToStruct != nil {
		lp = append(lp, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToUser && typ.RestrictedToUser {
		lp = append(lp, jen.ID("userID"))
	}

	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}

func (typ DataType) BuildInterfaceDefinitionExistenceMethodParams(proj *Project) []jen.Code {
	return typ.buildGetSomethingParams(proj)
}

func (typ DataType) BuildInterfaceDefinitionRetrievalMethodParams(proj *Project) []jen.Code {
	return typ.buildGetSomethingParams(proj)
}

func (typ DataType) BuildInterfaceDefinitionArchiveMethodParams(proj *Project) []jen.Code {
	return typ.buildArchiveSomethingParams(proj)
}

func (typ DataType) BuildDBClientArchiveMethodParams(proj *Project) []jen.Code {
	return typ.buildArchiveSomethingParams(proj)
}

func (typ DataType) BuildDBClientRetrievalMethodParams(proj *Project) []jen.Code {
	return typ.buildGetSomethingParams(proj)
}

func (typ DataType) BuildDBClientExistenceMethodParams(proj *Project) []jen.Code {
	return typ.buildGetSomethingParams(proj)
}

func (typ DataType) BuildDBQuerierArchiveMethodParams(proj *Project) []jen.Code {
	return typ.buildArchiveSomethingParams(proj)
}

func (typ DataType) BuildDBQuerierArchiveQueryMethodParams(proj *Project) []jen.Code {
	params := typ.buildArchiveSomethingParams(proj)

	return params[1:]
}

func (typ DataType) BuildDBQuerierRetrievalMethodParams(proj *Project) []jen.Code {
	params := typ.buildGetSomethingParams(proj)

	return params[1:]
}

func (typ DataType) BuildDBQuerierRetrievalQueryMethodParams(proj *Project) []jen.Code {
	return typ.buildGetSomethingParams(proj)
}

func (typ DataType) BuildDBQuerierExistenceQueryMethodParams(proj *Project) []jen.Code {
	params := typ.buildGetSomethingParams(proj)

	return params[1:]
}

func (typ DataType) ModifyQueryBuildingStatementWithJoinClauses(proj *Project, qbStmt *jen.Statement) *jen.Statement {
	if typ.BelongsToStruct != nil {
		qbStmt = qbStmt.Dotln("Join").Call(
			jen.IDf("%sOn%sJoinClause", typ.BelongsToStruct.PluralUnexportedVarName(), typ.Name.Plural()),
		)
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i := len(owners) - 1; i >= 0; i-- {
		pt := owners[i]

		if pt.BelongsToStruct != nil {
			qbStmt = qbStmt.Dotln("Join").Call(
				jen.IDf("%sOn%sJoinClause", pt.BelongsToStruct.PluralUnexportedVarName(), pt.Name.Plural()),
			)
		}
	}

	return qbStmt
}

func (typ DataType) ModifyQueryBuilderWithJoinClauses(proj *Project, qb squirrel.SelectBuilder) squirrel.SelectBuilder {
	if typ.BelongsToStruct != nil {
		qb = qb.Join(
			fmt.Sprintf("%s ON %s.%s=%s.id", typ.BelongsToStruct.PluralRouteName(), typ.Name.PluralRouteName(), fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()), typ.BelongsToStruct.PluralRouteName()),
		)
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i := len(owners) - 1; i >= 0; i-- {
		pt := owners[i]

		if pt.BelongsToStruct != nil {
			qb = qb.Join(
				fmt.Sprintf("%s ON %s.%s=%s.id", pt.BelongsToStruct.PluralUnexportedVarName(), pt.Name.PluralRouteName(), fmt.Sprintf("belongs_to_%s", pt.BelongsToStruct.RouteName()), pt.BelongsToStruct.PluralUnexportedVarName()),
			)
		}
	}

	return qb
}

func (typ DataType) buildDBQuerierSingleInstanceQueryMethodConditionalClauses(proj *Project) []jen.Code {
	n := typ.Name
	uvn := n.UnexportedVarName()
	puvn := n.PluralUnexportedVarName()

	whereValues := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.IDf("%sTableName", puvn)).MapAssign().IDf("%sID", uvn),
	}
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		whereValues = append(
			whereValues,
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.id"),
				jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
			).MapAssign().IDf("%sID", pt.Name.UnexportedVarName()),
		)

		if pt.BelongsToUser && pt.RestrictedToUser {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
					jen.IDf("%sUserOwnershipColumn", pt.Name.PluralUnexportedVarName()),
				).MapAssign().ID("userID"),
			)
		}

		if pt.BelongsToStruct != nil {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
					jen.IDf("%sTableOwnershipColumn", pt.Name.PluralUnexportedVarName()),
				).MapAssign().IDf("%sID", pt.BelongsToStruct.UnexportedVarName()),
			)
		}
	}

	if typ.BelongsToStruct != nil {
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sUserOwnershipColumn", puvn)).MapAssign().ID("userID"))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierExistenceQueryMethodConditionalClauses(proj *Project) []jen.Code {
	return typ.buildDBQuerierSingleInstanceQueryMethodConditionalClauses(proj)
}

func (typ DataType) BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(proj *Project) squirrel.Eq {
	return typ.buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(proj)
}

func (typ DataType) BuildDBQuerierRetrievalQueryMethodConditionalClauses(proj *Project) []jen.Code {
	return typ.buildDBQuerierSingleInstanceQueryMethodConditionalClauses(proj)
}

type Coder interface {
	Code() jen.Code
}

var _ Coder = codeWrapper{}

type codeWrapper struct {
	repr jen.Code
}

func NewCodeWrapper(c jen.Code) Coder {
	return codeWrapper{}
}

func (c codeWrapper) Code() jen.Code {
	return c.repr
}

func (typ DataType) buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(proj *Project) squirrel.Eq {
	n := typ.Name
	uvn := n.UnexportedVarName()
	puvn := n.PluralUnexportedVarName()
	tableName := typ.Name.PluralRouteName()

	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): fmt.Sprintf("%sID", uvn),
	}
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		whereValues[fmt.Sprintf("%s.id", tableName)] = fmt.Sprintf("%sID", pt.Name.UnexportedVarName())

		if pt.BelongsToUser && pt.RestrictedToUser {
			whereValues[fmt.Sprintf("%s.belongs_to_user", tableName)] = "userID"
		}

		if pt.BelongsToStruct != nil {
			whereValues[fmt.Sprintf("%s.belongs_to_%s",
				tableName, pt.Name.RouteName(),
			)] = fmt.Sprintf("%sID", pt.BelongsToStruct.UnexportedVarName())
		}
	}

	if typ.BelongsToStruct != nil {
		whereValues[fmt.Sprintf("%s.%s", fmt.Sprintf("%sTableName", puvn), fmt.Sprintf("%sTableOwnershipColumn", puvn))] = fmt.Sprintf("%sID", typ.BelongsToStruct.UnexportedVarName())
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		whereValues[fmt.Sprintf("%s.%s", fmt.Sprintf("%sTableName", puvn), fmt.Sprintf("%sUserOwnershipColumn", puvn))] = "userID"
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierListRetrievalQueryMethodConditionalClauses(proj *Project) []jen.Code {
	n := typ.Name
	puvn := n.PluralUnexportedVarName()

	whereValues := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.IDf("%sTableName", typ.Name.PluralUnexportedVarName())).MapAssign().Nil(),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		whereValues = append(
			whereValues,
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.id"),
				jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
			).MapAssign().IDf("%sID", pt.Name.UnexportedVarName()),
		)

		if pt.BelongsToUser && pt.RestrictedToUser {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
					jen.IDf("%sUserOwnershipColumn", pt.Name.PluralUnexportedVarName()),
				).MapAssign().ID("userID"),
			)
		}

		if pt.BelongsToStruct != nil {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
					jen.IDf("%sTableOwnershipColumn", pt.Name.PluralUnexportedVarName()),
				).MapAssign().IDf("%sID", pt.BelongsToStruct.UnexportedVarName()),
			)
		}
	}

	if typ.BelongsToStruct != nil {
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sUserOwnershipColumn", puvn)).MapAssign().ID("userID"))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierExistenceMethodParams(proj *Project) []jen.Code {
	return typ.buildGetSomethingParams(proj)
}

func (typ DataType) buildGetSomethingArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		params = append(params, jen.ID("userID"))
	}

	return params
}

func (typ DataType) buildArchiveSomethingArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}
	uvn := typ.Name.UnexportedVarName()

	if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID("userID"))
	}

	return params
}

func (typ DataType) BuildDBClientExistenceMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetSomethingArgs(proj)
}

func (typ DataType) BuildDBClientRetrievalMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetSomethingArgs(proj)
}

func (typ DataType) BuildDBClientArchiveMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildArchiveSomethingArgs(proj)
}

func (typ DataType) BuildDBQuerierExistenceQueryBuildingArgs(proj *Project) []jen.Code {
	params := typ.buildGetSomethingArgs(proj)

	return params[1:]
}

func (typ DataType) BuildDBQuerierRetrievalQueryBuildingArgs(proj *Project) []jen.Code {
	params := typ.buildGetSomethingArgs(proj)

	return params[1:]
}

func (typ DataType) BuildDBQuerierArchiveQueryBuildingArgs(proj *Project) []jen.Code {
	params := typ.buildArchiveSomethingArgs(proj)

	return params[1:]
}

func (typ DataType) BuildInterfaceDefinitionExistenceMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetSomethingArgs(proj)
}

func (typ DataType) BuildInterfaceDefinitionRetrievalMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetSomethingArgs(proj)
}

func (typ DataType) BuildInterfaceDefinitionArchiveMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildArchiveSomethingArgs(proj)
}

func (typ DataType) buildGetSomethingArgsWithExampleVariables(proj *Project, includeCtx bool) []jen.Code {
	params := []jen.Code{}

	if includeCtx {
		params = append(params, ctxVar())
	}

	owners := proj.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for _, pt := range owners {
		pts := pt.Name.Singular()
		params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	return params
}

func (typ DataType) BuildHTTPClientRetrievalTestCallArgs(proj *Project) []jen.Code {
	return typ.buildGetSomethingArgsWithExampleVariables(proj, true)
}

func (typ DataType) buildSingleInstanceQueryTestCallArgs(proj *Project) []jen.Code {
	params := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for _, pt := range owners {
		pts := pt.Name.Singular()
		params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return params
}

func (typ DataType) buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj *Project) []jen.Code {
	params := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			params = append(params, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			pts := pt.Name.Singular()
			params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
		}
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return params
}

func (typ DataType) BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(proj *Project) []jen.Code {
	return typ.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj)
}

func (typ DataType) BuildDBQuerierBuildSomethingExistsQueryTestExpectedArgs(proj *Project) []jen.Code {
	return typ.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj)
}

func (typ DataType) BuildDBQuerierRetrievalQueryTestCallArgs(proj *Project) []jen.Code {
	return typ.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj)
}

func (typ DataType) BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}

	if typ.BelongsToUser && typ.RestrictedToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}

	lines = append(lines, jen.ID(constants.FilterVarName).Assign().Qual(proj.FakeModelsPackage(), "BuildFleshedOutQueryFilter").Call())

	return lines
}

func (typ DataType) BuildDBQuerierCreateSomethingQueryBuilderTestPreQueryLines(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDBQuerierUpdateSomethingQueryBuilderTestPreQueryLines(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildGetSomethingLogValues(proj *Project) jen.Code {
	params := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.Litf("%s_id", pt.Name.RouteName()).Op(":").IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.Litf("%s_id", typ.Name.RouteName()).Op(":").IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		params = append(params, jen.Lit("user_id").Op(":").ID("userID"))
	}

	return jen.Map(jen.ID("string")).Interface().Valuesln(params...)
}

func (typ DataType) BuildGetListOfSomethingLogValues(proj *Project) *jen.Statement {
	params := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.Litf("%s_id", pt.Name.RouteName()).Op(":").IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if typ.RestrictedToUserAtSomeLevel(proj) {
		params = append(params, jen.Lit("user_id").Op(":").ID("userID"))
	}

	if len(params) > 0 {
		return jen.Map(jen.ID("string")).Interface().Valuesln(params...)
	}

	return nil
}

func (typ DataType) buildGetListOfSomethingParams(proj *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(proj) {
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

func (typ DataType) BuildMockDataManagerListRetrievalMethodParams(proj *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(proj, false)
}

func (typ DataType) BuildInterfaceDefinitionListRetrievalMethodParams(proj *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(proj, true)
}

func (typ DataType) BuildDBClientListRetrievalMethodParams(proj *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(proj, false)
}

func (typ DataType) BuildDBQuerierListRetrievalMethodParams(proj *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(proj, false)
}

func (typ DataType) BuildDBQuerierListRetrievalQueryBuildingMethodParams(proj *Project) []jen.Code {
	params := typ.buildGetListOfSomethingParams(proj, false)

	return params[1:]
}

const creationObjectVarName = "input"

func (typ DataType) buildCreateSomethingParams(proj *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(creationObjectVarName).Op("*").IDf("%sCreationInput", sn))
	} else {
		params = append(params, jen.ID(creationObjectVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models", "v1"), fmt.Sprintf("%sCreationInput", sn)))
	}

	return params
}

func (typ DataType) BuildMockInterfaceDefinitionCreationMethodParams(proj *Project) []jen.Code {
	return typ.buildCreateSomethingParams(proj, false)
}

func (typ DataType) BuildInterfaceDefinitionCreationMethodParams(proj *Project) []jen.Code {
	return typ.buildCreateSomethingParams(proj, true)
}

func (typ DataType) BuildDBClientCreationMethodParams(proj *Project) []jen.Code {
	return typ.buildCreateSomethingParams(proj, false)
}

func (typ DataType) BuildDBQuerierCreationMethodParams(proj *Project) []jen.Code {
	return typ.buildCreateSomethingParams(proj, false)
}

func (typ DataType) BuildDBQuerierCreationQueryBuildingMethodParams(proj *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(creationObjectVarName).Op("*").ID(sn))
	} else {
		params = append(params, jen.ID(creationObjectVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn))
	}

	return params
}

func (typ DataType) buildCreateSomethingArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar(), jen.ID(creationObjectVarName)}

	return params
}

func (typ DataType) BuildMockInterfaceDefinitionCreationMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildCreateSomethingArgs(proj)
}

func (typ DataType) BuildDBQuerierCreationMethodQueryBuildingArgs(proj *Project) []jen.Code {
	params := typ.buildCreateSomethingArgs(proj)

	return params[1:]
}

func (typ DataType) BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(proj *Project) []jen.Code {
	params := []jen.Code{}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.BelongsToUser && typ.RestrictedToUser {
		lp = append(lp, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}
	lp = append(lp, jen.ID(constants.FilterVarName))

	params = append(params, lp...)

	return params
}

func (typ DataType) BuildArgsForDBQuerierTestOfUpdateQueryBuilder(proj *Project) []jen.Code {
	params := []jen.Code{}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lp = append(lp, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}

	params = append(params, lp...)
	params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())))

	return params
}

func (typ DataType) BuildArgsForDBQuerierTestOfArchiveQueryBuilder(proj *Project) []jen.Code {
	return typ.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj)
}

func (typ DataType) BuildArgsForDBQuerierTestOfUpdateMethod(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			lp = append(lp, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}

	params = append(params, lp...)
	params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())))

	return params
}

func (typ DataType) BuildDBQuerierCreationMethodArgsToUseFromMethodTest(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			lp = append(lp, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			lp = append(lp, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}

	params = append(params, lp...)
	params = append(params, jen.ID(buildFakeVarName("Input")))

	return params
}

func (typ DataType) BuildArgsToUseForDBQuerierCreationQueryBuildingTest(proj *Project) []jen.Code {
	params := []jen.Code{}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			lp = append(lp, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			lp = append(lp, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
	}

	params = append(params, lp...)
	params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())))

	return params
}

func (typ DataType) BuildDBClientCreationMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildCreateSomethingArgs(proj)
}

func (typ DataType) buildUpdateSomethingParams(proj *Project, updatedVarName string, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	if updatedVarName == "" {
		panic("buildUpdateSomethingParams called with empty updatedVarName")
	}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(updatedVarName).Op("*").ID(sn))
	} else {
		params = append(params, jen.ID(updatedVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn))
	}

	return params
}

func (typ DataType) BuildDBClientUpdateMethodParams(proj *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(proj, updatedVarName, false)
}

func (typ DataType) BuildDBQuerierUpdateMethodParams(proj *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(proj, updatedVarName, false)
}

func (typ DataType) BuildDBQuerierUpdateQueryBuildingMethodParams(proj *Project, updatedVarName string) []jen.Code {
	params := typ.buildUpdateSomethingParams(proj, updatedVarName, false)

	return params[1:]
}

func (typ DataType) BuildInterfaceDefinitionUpdateMethodParams(proj *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(proj, updatedVarName, true)
}

func (typ DataType) BuildMockDataManagerUpdateMethodParams(proj *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(proj, updatedVarName, false)
}

func (typ DataType) buildUpdateSomethingArgsWithExampleVars(proj *Project, updatedVarName string) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.ID(pt.Name.Singular()).Dot("ID"))
	}
	if len(lp) >= 1 {
		params = append(params, jen.List(lp[:len(lp)-1]...))
	}
	params = append(params, jen.ID(updatedVarName))

	return params
}

func (typ DataType) buildUpdateSomethingArgs(proj *Project, updatedVarName string) []jen.Code {
	params := []jen.Code{ctxVar(), jen.ID(updatedVarName)}

	return params
}

func (typ DataType) BuildDBClientUpdateMethodCallArgs(proj *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingArgs(proj, updatedVarName)
}

func (typ DataType) BuildDBQuerierUpdateMethodArgs(proj *Project, updatedVarName string) []jen.Code {
	params := typ.buildUpdateSomethingArgs(proj, updatedVarName)

	return params[1:]
}

func (typ DataType) BuildMockDataManagerUpdateMethodCallArgs(proj *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingArgs(proj, updatedVarName)
}

func (typ DataType) buildGetListOfSomethingArgs(proj *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(proj) {
		params = append(params, jen.ID("userID"))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func (typ DataType) BuildDBClientListRetrievalMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetListOfSomethingArgs(proj)
}

func (typ DataType) BuildDBQuerierListRetrievalMethodArgs(proj *Project) []jen.Code {
	params := typ.buildGetListOfSomethingArgs(proj)

	return params[1:]
}

func (typ DataType) BuildMockDataManagerListRetrievalMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetListOfSomethingArgs(proj)
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWithOwnerStruct(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())

		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())
	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildCreationRequestMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDependentObjectsForDBQueriersExistenceMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)
}

func (typ DataType) BuildDependentObjectsForDBQueriersCreationMethodTest(proj *Project) []jen.Code {
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

	lines = append(lines,
		jen.IDf(buildFakeVarName("Input")).Assign().
			Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).
			Call(jen.ID(buildFakeVarName(sn))),
	)

	return lines
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}

	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())
	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	return lines
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) BuildHTTPClientRetrievalMethodTestDependentObjects(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildArchiveRequestMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildExistenceRequestMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientExistenceMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildRetrievalRequestMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientRetrievalMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientArchiveMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(proj)
}

func (typ DataType) buildDependentObjectsForHTTPClientListRetrievalTest(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}

	return lines
}

func (typ DataType) BuildDependentObjectsForHTTPClientListRetrievalTest(proj *Project) []jen.Code {
	return typ.buildDependentObjectsForHTTPClientListRetrievalTest(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildListRetrievalRequestMethodTest(proj *Project) []jen.Code {
	return typ.buildDependentObjectsForHTTPClientListRetrievalTest(proj)
}

func (typ DataType) buildVarDeclarationsOfDependentStructsForUpdateFunction(proj *Project) []jen.Code {
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

func (typ DataType) BuildDependentObjectsForHTTPClientUpdateMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsForUpdateFunction(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildUpdateRequestMethodTest(proj *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsForUpdateFunction(proj)
}

func (typ DataType) BuildDependentObjectsForHTTPClientCreationMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines,
			jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call(),
		)
	}

	lines = append(lines,
		jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
	)

	if typ.BelongsToStruct != nil {
		lines = append(lines,
			jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"),
		)
	}

	return lines
}

func (typ DataType) BuildFormatStringForHTTPClientExistenceMethodTest(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientRetrievalMethodTest(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientUpdateMethodTest(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientArchiveMethodTest(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientListMethodTest(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s", typ.Name.PluralRouteName())

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientCreateMethodTest(proj *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += typ.Name.PluralRouteName()

	return modelRoute
}

func (typ DataType) BuildFormatCallArgsForHTTPClientRetrievalMethodTest(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForHTTPClientExistenceMethodTest(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForHTTPClientListMethodTest(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForHTTPClientCreationMethodTest(proj *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForHTTPClientUpdateTest(proj *Project) (args []jen.Code) {
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

func (typ DataType) BuildArgsForHTTPClientExistenceRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
	}

	return params
}

func (typ DataType) BuildParamsForHTTPClientExistenceRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, jen.List(listParams...).Uint64())

	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64())
	}

	return params
}

func (typ DataType) BuildParamsForHTTPClientExistenceMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, jen.List(listParams...).Uint64())
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64())
	}

	return params
}

func (typ DataType) BuildArgsForHTTPClientCreateRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 {
				continue
			} else {
				listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
			}
		}
		params = append(params, listParams...)
	}

	params = append(params, jen.ID("input"))

	return params
}

func (typ DataType) BuildArgsForHTTPClientRetrievalRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}

		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
	}

	return params
}

func (typ DataType) BuildParamsForHTTPClientRetrievalRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, jen.List(listParams...).Uint64())

	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64())
	}

	return params
}

func (typ DataType) BuildParamsForHTTPClientRetrievalMethod(proj *Project, call bool) []jen.Code {
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

func (typ DataType) BuildParamsForHTTPClientCreateRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				continue
			} else {
				listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
			}
		}
		if len(listParams) > 0 {
			params = append(params, jen.List(listParams...).Uint64())
		}
	}

	params = append(params, jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))

	return params
}

func (typ DataType) BuildParamsForHTTPClientCreateMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i > len(parents)-2 {
				continue
			}
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		if len(listParams) > 0 {
			params = append(params, jen.List(listParams...).Uint64())
		}
	}

	params = append(params, jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))

	return params
}

func (typ DataType) BuildParamsForHTTPClientUpdateRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, jen.List(listParams...).Uint64())
		}
	}

	params = append(params, jen.ID(typ.Name.UnexportedVarName()).PointerTo().Qual(proj.ModelsV1Package(), typ.Name.Singular()))

	return params
}

func (typ DataType) BuildArgsForHTTPClientUpdateRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, listParams...)
		}
	}

	params = append(params, jen.ID(typ.Name.UnexportedVarName()))

	return params
}

func (typ DataType) BuildParamsForHTTPClientUpdateMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, jen.List(listParams...).Uint64())
		}
	}

	params = append(params, jen.ID(typ.Name.UnexportedVarName()).PointerTo().Qual(proj.ModelsV1Package(), typ.Name.Singular()))

	return params
}

func (typ DataType) BuildParamsForHTTPClientArchiveRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, jen.List(listParams...).Uint64())
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64())
	}

	return params
}

func (typ DataType) BuildArgsForHTTPClientArchiveRequestBuildingMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()))
	}

	return params
}

func (typ DataType) BuildParamsForHTTPClientArchiveMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, jen.List(listParams...).Uint64())
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64())
	}

	return params
}

func (typ DataType) buildParamsForMethodThatHandlesAnInstanceWithStructs(proj *Project) []jen.Code {
	owners := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(owners) > 0 {
		for _, pt := range owners {
			listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
		listParams = append(listParams, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))
	}

	return params
}

func (typ DataType) BuildArgsForHTTPClientExistenceRequestBuildingMethodTest(proj *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)
}

func (typ DataType) BuildArgsForDBQuerierExistenceMethodTest(proj *Project) []jen.Code {
	params := typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot(constants.UserOwnershipFieldName))
	}

	return params
}

func (typ DataType) BuildArgsForDBQuerierRetrievalMethodTest(proj *Project) []jen.Code {
	params := typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)

	if typ.BelongsToUser && typ.RestrictedToUser {
		params = append(params, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot(constants.UserOwnershipFieldName))
	}

	return params
}

func (typ DataType) BuildArgsForHTTPClientExistenceMethodTest(proj *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)
}

func (typ DataType) BuildArgsForHTTPClientRetrievalRequestBuilderMethodTest(proj *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)
}

func (typ DataType) BuildArgsForHTTPClientArchiveRequestBuildingMethodTest(proj *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)
}
func (typ DataType) BuildArgsForHTTPClientArchiveMethodTest(proj *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)
}

func (typ DataType) BuildArgsForHTTPClientArchiveMethodTestURLFormatCall(proj *Project) []jen.Code {
	params := typ.BuildArgsForHTTPClientArchiveMethodTest(proj)

	return params[1:]
}

func (typ DataType) BuildArgsForHTTPClientMethodTest(proj *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(proj)
}

func (typ DataType) BuildHTTPClientCreationRequestBuildingMethodArgsForTest(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				continue
			} else {
				listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
			}
		}
		params = append(params, listParams...)
	}

	params = append(params, jen.ID(buildFakeVarName("Input")))

	return params
}

func (typ DataType) BuildHTTPClientCreationMethodArgsForTest(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == len(parents)-1 && typ.BelongsToStruct != nil {
				continue
			} else {
				listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
			}
		}
		params = append(params, listParams...)
	}

	params = append(params, jen.ID(buildFakeVarName("Input")))

	return params
}

func (typ DataType) BuildArgsForHTTPClientListRequestMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		params = append(params, listParams...)
	}

	params = append(params, jen.ID(constants.FilterVarName))

	return params
}

func (typ DataType) BuildParamsForHTTPClientListRequestMethod(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		params = append(params, jen.List(listParams...).Uint64())
	}

	params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))

	return params
}

func (typ DataType) BuildParamsForHTTPClientMethodThatFetchesAList(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		params = append(params, jen.List(listParams...).Uint64())
	}

	params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"))

	return params
}

func (typ DataType) BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
		params = append(params, listParams...)
	}
	params = append(params, jen.ID(constants.FilterVarName))

	return params
}

func (typ DataType) BuildCallArgsForHTTPClientListRetrievalMethodTest(proj *Project) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
		params = append(params, listParams...)
	}
	params = append(params, jen.ID(constants.FilterVarName))

	return params
}

func (typ DataType) BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(proj *Project) []jen.Code {
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

func (typ DataType) BuildCallArgsForHTTPClientUpdateMethodTest(proj *Project) []jen.Code {
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

func (typ DataType) buildRequisiteFakeVarDecs(proj *Project, createCtx bool) []jen.Code {
	lines := []jen.Code{}
	if createCtx {
		lines = append(lines, constants.CreateCtx(), jen.Line())
	}

	if !(typ.BelongsToUser && typ.RestrictedToUser) && typ.RestrictedToUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	return lines
}

func (typ DataType) buildRequisiteFakeVarDecForModifierFuncs(proj *Project, createCtx bool) []jen.Code {
	lines := []jen.Code{}

	if createCtx {
		lines = append(lines, constants.CreateCtx(), jen.Line())
	}

	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	return lines
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientExistenceMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecs(proj, true)
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecs(proj, false)
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientCreateMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecForModifierFuncs(proj, true)
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientArchiveMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecForModifierFuncs(proj, true)
}

func (typ DataType) BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())

	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName).Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) buildRequisiteFakeVarDecsForListFunction(proj *Project) []jen.Code {
	lines := []jen.Code{}

	if !(typ.BelongsToUser && typ.RestrictedToUser) && typ.RestrictedToUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientListRetrievalMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecsForListFunction(proj)
}

func (typ DataType) BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecsForListFunction(proj)
}

func (typ DataType) buildRequisiteFakeVarCallArgsForCreation(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser && typ.RestrictedToUser {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBQuerierCreationMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForCreation(proj)
}

func (typ DataType) buildRequisiteFakeVarCallArgs(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser && typ.RestrictedToUser {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("BelongsToUser"))
	} else if typ.RestrictedToUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBClientExistenceMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgs(proj)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBClientRetrievalMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgs(proj)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser && typ.RestrictedToUser {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return lines
}

func (typ DataType) BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(proj *Project) []jen.Code {
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

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.BelongsToUser && typ.RestrictedToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	lines = append(lines, jen.ID(constants.FilterVarName))

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()
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

func (typ DataType) BuildCallArgsForDBClientCreationMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{constants.CtxVar()}

	const (
		inputVarName = "exampleInput"
	)

	lines = append(lines, jen.ID(inputVarName))

	return lines
}

func (typ DataType) BuildCallArgsForDBClientListRetrievalMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.RestrictedToUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteVarsForDBClientUpdateMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{
		constants.CreateCtx(),
		jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call(),
	}

	return lines
}
func (typ DataType) BuildCallArgsForDBClientUpdateMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{jen.ID(buildFakeVarName(typ.Name.Singular()))}

	return lines
}
