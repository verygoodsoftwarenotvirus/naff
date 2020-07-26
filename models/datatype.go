package models

import (
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"github.com/Masterminds/squirrel"
)

// DataType represents a data model
type DataType struct {
	Name             wordsmith.SuperPalabra
	Struct           *types.Struct
	BelongsToUser    bool
	BelongsToNobody  bool
	RestrictedToUser bool
	SearchEnabled    bool
	BelongsToStruct  wordsmith.SuperPalabra
	Fields           []DataField
}

// DataField represents a data model's field
type DataField struct {
	Name                  wordsmith.SuperPalabra
	Type                  string
	UnderlyingType        types.Type
	Pos                   token.Pos
	Pointer               bool
	DefaultValue          string
	ValidForCreationInput bool
	ValidForUpdateInput   bool
	BelongsToEnumeration  wordsmith.SuperPalabra
}

func buildFakeVarName(typName string) string {
	typName = strings.ToUpper(string(typName[0])) + typName[1:]
	return fmt.Sprintf("example%s", typName)
}

// ctxParam is a shorthand for a context param
func ctxParam() jen.Code {
	return ctxVar().Qual("context", "Context")
}

// ctxParam is a shorthand for a context param
func ctxVar() *jen.Statement {
	return jen.ID("ctx")
}

func (typ DataType) OwnedByAUserAtSomeLevel(proj *Project) bool {
	for _, o := range proj.FindOwnerTypeChain(typ) {
		if o.BelongsToUser {
			return true
		}
	}

	return typ.BelongsToUser
}

func (typ DataType) RestrictedToUserAtSomeLevel(proj *Project) bool {
	for _, o := range proj.FindOwnerTypeChain(typ) {
		if o.BelongsToUser && o.RestrictedToUser {
			return true
		}
	}

	return typ.BelongsToUser && typ.RestrictedToUser
}

func (typ DataType) MultipleOwnersBelongingToUser(proj *Project) bool {
	var count uint

	if typ.BelongsToUser {
		count += 1
	}

	for _, o := range proj.FindOwnerTypeChain(typ) {
		if o.BelongsToUser {
			count += 1
		}
	}

	return count > 1
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

	if typ.BelongsToUser {
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

func (typ DataType) buildJoinClause(fromTableName, onTableName, fkColumnName string) string {
	if typ.BelongsToStruct != nil {
		return fmt.Sprintf("%s ON %s.%s=%s.id", fromTableName, onTableName, fmt.Sprintf("belongs_to_%s", fkColumnName), fromTableName)
	}

	panic("buildJoinClause called on struct that doesn't belong to anything!")
	return ""
}

func (typ DataType) ModifyQueryBuilderWithJoinClauses(proj *Project, qb squirrel.SelectBuilder) squirrel.SelectBuilder {
	if typ.BelongsToStruct != nil {
		qb = qb.Join(typ.buildJoinClause(typ.BelongsToStruct.PluralRouteName(), typ.Name.PluralRouteName(), typ.BelongsToStruct.RouteName()))
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i := len(owners) - 1; i >= 0; i-- {
		pt := owners[i]

		if pt.BelongsToStruct != nil {
			qb = qb.Join(
				typ.buildJoinClause(pt.BelongsToStruct.PluralUnexportedVarName(), pt.Name.PluralRouteName(), pt.BelongsToStruct.RouteName()),
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
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.ID("idColumn")).MapAssign().IDf("%sID", uvn),
	}
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		whereValues = append(
			whereValues,
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
				jen.ID("idColumn"),
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

type Coder interface {
	Code() jen.Code
}

var _ Coder = codeWrapper{}

type codeWrapper struct {
	repr jen.Code
}

// NewCodeWrapper creates a new codeWrapper from some code
func NewCodeWrapper(c jen.Code) Coder {
	return codeWrapper{repr: c}
}

func (c codeWrapper) Code() jen.Code {
	return c.repr
}

func (typ DataType) buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(proj *Project) squirrel.Eq {
	n := typ.Name
	sn := n.Singular()
	tableName := typ.Name.PluralRouteName()

	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): NewCodeWrapper(jen.ID(buildFakeVarName(sn)).Dot("ID")),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pTableName := pt.Name.PluralRouteName()
		whereValues[fmt.Sprintf("%s.id", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.Name.UnexportedVarName())).Dot("ID"))

		if pt.BelongsToUser && pt.RestrictedToUser {
			whereValues[fmt.Sprintf("%s.belongs_to_user", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName("User")).Dot("ID"))
		}

		if pt.BelongsToStruct != nil {
			whereValues[fmt.Sprintf("%s.belongs_to_%s", pTableName, pt.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}

	if typ.BelongsToStruct != nil {
		whereValues[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		whereValues[fmt.Sprintf("%s.belongs_to_user", tableName)] = NewCodeWrapper(jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(proj *Project) squirrel.Eq {
	return typ.buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(proj)
}

func (typ DataType) BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(proj *Project) squirrel.Eq {
	return typ.buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(proj)
}

func (typ DataType) BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(proj *Project) squirrel.Eq {
	tableName := typ.Name.PluralRouteName()

	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.archived_on", tableName): nil,
	}
	for _, pt := range proj.FindOwnerTypeChain(typ) {
		pTableName := pt.Name.PluralRouteName()
		whereValues[fmt.Sprintf("%s.id", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.Name.UnexportedVarName())).Dot("ID"))

		if pt.BelongsToUser && pt.RestrictedToUser {
			whereValues[fmt.Sprintf("%s.belongs_to_user", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName("User")).Dot("ID"))
		}

		if pt.BelongsToStruct != nil {
			whereValues[fmt.Sprintf("%s.belongs_to_%s", pTableName, pt.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}

	if typ.BelongsToStruct != nil {
		whereValues[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		whereValues[fmt.Sprintf("%s.belongs_to_user", tableName)] = NewCodeWrapper(jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierRetrievalQueryMethodConditionalClauses(proj *Project) []jen.Code {
	return typ.buildDBQuerierSingleInstanceQueryMethodConditionalClauses(proj)
}

func (typ DataType) BuildDBQuerierListRetrievalQueryMethodConditionalClauses(proj *Project) []jen.Code {
	n := typ.Name
	puvn := n.PluralUnexportedVarName()

	whereValues := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", typ.Name.PluralUnexportedVarName()), jen.ID("archivedOnColumn")).MapAssign().Nil(),
	}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		whereValues = append(
			whereValues,
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
				jen.ID("idColumn"),
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

func (typ DataType) buildArchiveSomethingArgs() []jen.Code {
	params := []jen.Code{ctxVar()}
	uvn := typ.Name.UnexportedVarName()

	if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.BelongsToUser {
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

func (typ DataType) BuildDBClientArchiveMethodCallArgs() []jen.Code {
	return typ.buildArchiveSomethingArgs()
}

func (typ DataType) BuildDBQuerierExistenceQueryBuildingArgs(proj *Project) []jen.Code {
	params := typ.buildGetSomethingArgs(proj)

	return params[1:]
}

func (typ DataType) BuildDBQuerierRetrievalQueryBuildingArgs(proj *Project) []jen.Code {
	params := typ.buildGetSomethingArgs(proj)

	return params[1:]
}

func (typ DataType) BuildDBQuerierArchiveQueryBuildingArgs() []jen.Code {
	params := typ.buildArchiveSomethingArgs()

	return params[1:]
}

func (typ DataType) BuildInterfaceDefinitionExistenceMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetSomethingArgs(proj)
}

func (typ DataType) BuildInterfaceDefinitionRetrievalMethodCallArgs(proj *Project) []jen.Code {
	return typ.buildGetSomethingArgs(proj)
}

func (typ DataType) BuildInterfaceDefinitionArchiveMethodCallArgs() []jen.Code {
	return typ.buildArchiveSomethingArgs()
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

func (typ DataType) buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(proj *Project) []jen.Code {
	owners := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	args := []jen.Code{constants.CtxVar()}

	sn := typ.Name.Singular()
	if len(owners) > 0 {
		for _, pt := range owners {
			listParams = append(listParams, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
		}
		listParams = append(listParams, jen.ID(buildFakeVarName(sn)).Dot("ID"))

		if len(listParams) > 0 {
			args = append(args, listParams...)
		}
	} else {
		args = append(args, jen.ID(buildFakeVarName(sn)).Dot("ID"))
	}

	if typ.RestrictedToUserAtSomeLevel(proj) {
		args = append(args, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return args
}

func (typ DataType) BuildArgsForDBQuerierExistenceMethodTest(proj *Project) []jen.Code {
	params := typ.buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(proj)

	return params
}

func (typ DataType) BuildArgsForDBQuerierRetrievalMethodTest(proj *Project) []jen.Code {
	params := typ.buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(proj)

	return params
}

func (typ DataType) BuildArgsForServiceRouteExistenceCheck(proj *Project) []jen.Code {
	owners := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	args := []jen.Code{constants.CtxVar()}

	uvn := typ.Name.UnexportedVarName()
	if len(owners) > 0 {
		for _, pt := range owners {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", uvn))

		if len(listParams) > 0 {
			args = append(args, listParams...)
		}
	} else {
		args = append(args, jen.IDf("%sID", uvn))
	}

	if typ.RestrictedToUserAtSomeLevel(proj) {
		args = append(args, jen.ID("userID"))
	}

	return args
}

func (typ DataType) buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj *Project) []jen.Code {
	params := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for _, pt := range owners {
		pts := pt.Name.Singular()
		params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		params = append(params, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return params
}

func (typ DataType) BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(proj *Project) []jen.Code {
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
	if typ.RestrictedToUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}

	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())

		if pt.BelongsToUser && typ.RestrictedToUser {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			btssn := pt.BelongsToStruct.Singular()
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", btssn).Equals().ID(buildFakeVarName(btssn)).Dot("ID"))
		}
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

func (typ DataType) buildCreateSomethingArgs() []jen.Code {
	params := []jen.Code{ctxVar(), jen.ID(creationObjectVarName)}

	return params
}

func (typ DataType) BuildMockInterfaceDefinitionCreationMethodCallArgs() []jen.Code {
	return typ.buildCreateSomethingArgs()
}

func (typ DataType) BuildDBQuerierCreationMethodQueryBuildingArgs() []jen.Code {
	params := typ.buildCreateSomethingArgs()

	return params[1:]
}

func (typ DataType) BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(proj *Project) []jen.Code {
	params := []jen.Code{}

	lp := []jen.Code{}
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.RestrictedToUserAtSomeLevel(proj) {
		lp = append(lp, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}
	lp = append(lp, jen.ID(constants.FilterVarName))

	params = append(params, lp...)

	return params
}

func (typ DataType) BuildArgsForDBQuerierTestOfUpdateQueryBuilder() []jen.Code {
	params := []jen.Code{jen.ID(buildFakeVarName(typ.Name.Singular()))}

	return params
}

func (typ DataType) BuildArgsForDBQuerierTestOfArchiveQueryBuilder() []jen.Code {
	args := []jen.Code{}

	if typ.BelongsToStruct != nil {
		args = append(args, jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	args = append(args, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	if typ.BelongsToUser {
		args = append(args, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return args
}

func (typ DataType) BuildArgsForDBQuerierTestOfUpdateMethod() []jen.Code {
	params := []jen.Code{ctxVar(), jen.ID(buildFakeVarName(typ.Name.Singular()))}

	return params
}

func (typ DataType) BuildDBQuerierCreationMethodArgsToUseFromMethodTest() []jen.Code {
	params := []jen.Code{ctxVar(), jen.ID(buildFakeVarName("Input"))}

	return params
}

func (typ DataType) BuildArgsToUseForDBQuerierCreationQueryBuildingTest() []jen.Code {
	params := []jen.Code{jen.ID(buildFakeVarName(typ.Name.Singular()))}

	return params
}

func (typ DataType) BuildDBClientCreationMethodCallArgs() []jen.Code {
	return typ.buildCreateSomethingArgs()
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

func (typ DataType) buildUpdateSomethingArgs(updatedVarName string) []jen.Code {
	params := []jen.Code{ctxVar(), jen.ID(updatedVarName)}

	return params
}

func (typ DataType) BuildDBClientUpdateMethodCallArgs(updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingArgs(updatedVarName)
}

func (typ DataType) BuildDBQuerierUpdateMethodArgs(updatedVarName string) []jen.Code {
	params := typ.buildUpdateSomethingArgs(updatedVarName)

	return params[1:]
}

func (typ DataType) BuildMockDataManagerUpdateMethodCallArgs(updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingArgs(updatedVarName)
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
	if typ.OwnedByAUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}

	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())

		if pt.BelongsToUser {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
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
	lines := typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(proj)

	sn := typ.Name.Singular()
	lines = append(lines, jen.ID(buildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.ID(buildFakeVarName(sn))))

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

func (typ DataType) BuildFormatStringForHTTPClientSearchMethodTest() (path string) {
	modelRoute := "/" + filepath.Join("api", "v1", typ.Name.PluralRouteName(), "search")

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
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		args = append(args, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	args = append(args, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	return args
}

func (typ DataType) BuildFormatCallArgsForHTTPClientListMethodTest(proj *Project) (args []jen.Code) {
	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		args = append(args, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	return args
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

	if typ.OwnedByAUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
		if pt.BelongsToUser {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}

	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}
	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	return lines
}

func (typ DataType) buildRequisiteFakeVarDecForModifierFuncs(proj *Project, createCtx bool) []jen.Code {
	lines := []jen.Code{}

	if createCtx {
		lines = append(lines, constants.CreateCtx(), jen.Line())
	}
	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}
	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientExistenceMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecs(proj, true)
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecs(proj, false)
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientCreateMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{constants.CreateCtx(), jen.Line()}

	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientArchiveMethodTest(proj *Project) []jen.Code {
	var lines []jen.Code

	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}

	return append([]jen.Code{
		constants.CreateCtx(),
		jen.Line(),
		jen.Var().ID("expected").Error(),
		jen.Line(),
	}, lines...)
}

func (typ DataType) BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(proj *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()
	owners := proj.FindOwnerTypeChain(typ)

	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
		if pt.BelongsToUser {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
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

func (typ DataType) BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj *Project, includeFilter bool) []jen.Code {
	lines := []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)

	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())

		if pt.BelongsToUser && typ.RestrictedToUser {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			btssn := pt.BelongsToStruct.Singular()
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", btssn).Equals().ID(buildFakeVarName(btssn)).Dot("ID"))
		}
	}

	if includeFilter {
		lines = append(lines, jen.ID(constants.FilterVarName).Assign().Qual(proj.FakeModelsPackage(), "BuildFleshedOutQueryFilter").Call())
	}

	return lines
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

func (typ DataType) buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(proj *Project) []jen.Code {

	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if (typ.BelongsToUser && typ.RestrictedToUser) || typ.RestrictedToUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(proj)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceReadHandlerTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(proj)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceCreateHandlerTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(proj)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceUpdateHandlerTest(proj *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(proj)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest() []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser {
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

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest() []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser {
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

	if typ.RestrictedToUserAtSomeLevel(proj) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	lines = append(lines, jen.ID(constants.FilterVarName))

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest() []jen.Code {
	lines := []jen.Code{constants.CtxVar()}

	sn := typ.Name.Singular()

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToUser {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildCallArgsForDBClientCreationMethodTest() []jen.Code {
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
		jen.Var().ID("expected").Error(),
		jen.Line(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(buildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call()
			}
			return jen.Null()
		}(),
		jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToUser").Equals().ID(buildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.Line(),
	}

	return lines
}
func (typ DataType) BuildCallArgsForDBClientUpdateMethodTest() []jen.Code {
	lines := []jen.Code{jen.ID(buildFakeVarName(typ.Name.Singular()))}

	return lines
}
