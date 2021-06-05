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
	Name                       wordsmith.SuperPalabra
	Struct                     *types.Struct
	BelongsToAccount           bool
	IsEnumeration              bool
	RestrictedToAccountMembers bool
	SearchEnabled              bool
	BelongsToStruct            wordsmith.SuperPalabra
	Fields                     []DataField
}

// DataField represents a data model's field
type DataField struct {
	Name                  wordsmith.SuperPalabra
	Type                  string
	UnderlyingType        types.Type
	Pos                   token.Pos
	IsPointer             bool
	DefaultValue          string
	ValidForCreationInput bool
	ValidForUpdateInput   bool
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

func (typ DataType) OwnedByAUserAtSomeLevel(p *Project) bool {
	for _, o := range p.FindOwnerTypeChain(typ) {
		if o.BelongsToAccount {
			return true
		}
	}

	return typ.BelongsToAccount
}

func (typ DataType) RestrictedToUserAtSomeLevel(p *Project) bool {
	for _, o := range p.FindOwnerTypeChain(typ) {
		if o.BelongsToAccount && o.RestrictedToAccountMembers {
			return true
		}
	}

	return typ.BelongsToAccount && typ.RestrictedToAccountMembers
}

func (typ DataType) MultipleOwnersBelongingToUser(p *Project) bool {
	var count uint

	if typ.BelongsToAccount {
		count += 1
	}

	for _, o := range p.FindOwnerTypeChain(typ) {
		if o.BelongsToAccount {
			count += 1
		}
	}

	return count > 1
}

func (typ DataType) buildGetSomethingParams(p *Project, includeAccountParam bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.RestrictedToUserAtSomeLevel(p) && includeAccountParam {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	return params
}

func (typ DataType) buildArchiveSomethingParams() []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	if typ.BelongsToStruct != nil {
		lp = append(lp, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToAccount {
		lp = append(lp, jen.ID("accountID"))
	}

	lp = append(lp, jen.ID("archivedBy"))

	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}

func (typ DataType) BuildInterfaceDefinitionExistenceMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) BuildInterfaceDefinitionRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) BuildInterfaceDefinitionArchiveMethodParams() []jen.Code {
	return typ.buildArchiveSomethingParams()
}

func (typ DataType) BuildInterfaceDefinitionAuditLogEntryRetrievalMethodParams() []jen.Code {
	return []jen.Code{
		ctxParam(),
		jen.IDf("%sID", typ.Name.UnexportedVarName()).Uint64(),
	}
}

func (typ DataType) BuildDBClientArchiveMethodParams() []jen.Code {
	return typ.buildArchiveSomethingParams()
}

func (typ DataType) BuildDBClientRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) BuildDBClientAuditLogRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, false)
}

func (typ DataType) BuildDBClientExistenceMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) BuildDBQuerierArchiveMethodParams() []jen.Code {
	return typ.buildArchiveSomethingParams()
}

func (typ DataType) BuildDBQuerierArchiveQueryMethodParams() []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	if typ.BelongsToStruct != nil {
		lp = append(lp, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.BelongsToAccount {
		lp = append(lp, jen.ID("accountID"))
	}

	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}

func (typ DataType) BuildDBQuerierRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) BuildDBQuerierRetrievalQueryMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) BuildDBQuerierExistenceQueryMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) ModifyQueryBuildingStatementWithJoinClauses(p *Project, qbStmt *jen.Statement) *jen.Statement {
	if typ.BelongsToStruct != nil {
		qbStmt = qbStmt.Dotln("Join").Call(
			jen.IDf("%sOn%sJoinClause", typ.BelongsToStruct.PluralUnexportedVarName(), typ.Name.Plural()),
		)
	}

	owners := p.FindOwnerTypeChain(typ)
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
}

func (typ DataType) ModifyQueryBuilderWithJoinClauses(p *Project, qb squirrel.SelectBuilder) squirrel.SelectBuilder {
	if typ.BelongsToStruct != nil {
		qb = qb.Join(typ.buildJoinClause(typ.BelongsToStruct.PluralRouteName(), typ.Name.PluralRouteName(), typ.BelongsToStruct.RouteName()))
	}

	owners := p.FindOwnerTypeChain(typ)
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

func (typ DataType) buildDBQuerierSingleInstanceQueryMethodConditionalClauses(p *Project) []jen.Code {
	n := typ.Name
	uvn := n.UnexportedVarName()
	puvn := n.PluralUnexportedVarName()

	whereValues := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.ID("idColumn")).MapAssign().IDf("%sID", uvn),
	}
	for _, pt := range p.FindOwnerTypeChain(typ) {
		whereValues = append(
			whereValues,
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
				jen.ID("idColumn"),
			).MapAssign().IDf("%sID", pt.Name.UnexportedVarName()),
		)

		if pt.BelongsToAccount && pt.RestrictedToAccountMembers {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
					jen.IDf("%sAccountOwnershipColumn", pt.Name.PluralUnexportedVarName()),
				).MapAssign().ID("accountID"),
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
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sAccountOwnershipColumn", puvn)).MapAssign().ID("accountID"))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierExistenceQueryMethodConditionalClauses(p *Project) []jen.Code {
	return typ.buildDBQuerierSingleInstanceQueryMethodConditionalClauses(p)
}

type Coder interface {
	Code() jen.Code
}

var _ Coder = codeWrapper{}

type codeWrapper struct {
	code jen.Code
}

// NewCodeWrapper creates a new codeWrapper from some code
func NewCodeWrapper(c jen.Code) Coder {
	return codeWrapper{code: c}
}

func (c codeWrapper) Code() jen.Code {
	return c.code
}

func (typ DataType) buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(p *Project) squirrel.Eq {
	n := typ.Name
	sn := n.Singular()
	tableName := typ.Name.PluralRouteName()

	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): NewCodeWrapper(jen.ID(buildFakeVarName(sn)).Dot("ID")),
	}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pTableName := pt.Name.PluralRouteName()

		whereValues[fmt.Sprintf("%s.id", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.Name.UnexportedVarName())).Dot("ID"))
		whereValues[fmt.Sprintf("%s.archived_on", pTableName)] = nil

		if pt.BelongsToAccount && pt.RestrictedToAccountMembers {
			whereValues[fmt.Sprintf("%s.belongs_to_account", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName("Account")).Dot("ID"))
		}

		if pt.BelongsToStruct != nil {
			whereValues[fmt.Sprintf("%s.belongs_to_%s", pTableName, pt.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}

	whereValues[fmt.Sprintf("%s.archived_on", tableName)] = nil

	if typ.BelongsToStruct != nil {
		whereValues[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		whereValues[fmt.Sprintf("%s.belongs_to_account", tableName)] = NewCodeWrapper(jen.ID(buildFakeVarName(sn)).Dot(constants.AccountOwnershipFieldName))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(p *Project) squirrel.Eq {
	return typ.buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(p)
}

func (typ DataType) BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(p *Project) squirrel.Eq {
	return typ.buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(p)
}

func (typ DataType) BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(p *Project) squirrel.Eq {
	tableName := typ.Name.PluralRouteName()

	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.archived_on", tableName): nil,
	}
	for _, pt := range p.FindOwnerTypeChain(typ) {
		pTableName := pt.Name.PluralRouteName()

		whereValues[fmt.Sprintf("%s.id", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.Name.UnexportedVarName())).Dot("ID"))
		whereValues[fmt.Sprintf("%s.archived_on", pTableName)] = nil

		if pt.BelongsToAccount && pt.RestrictedToAccountMembers {
			whereValues[fmt.Sprintf("%s.belongs_to_account", pTableName)] = NewCodeWrapper(jen.ID(buildFakeVarName("Account")).Dot("ID"))
		}

		if pt.BelongsToStruct != nil {
			whereValues[fmt.Sprintf("%s.belongs_to_%s", pTableName, pt.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}

	whereValues[fmt.Sprintf("%s.archived_on", tableName)] = nil

	if typ.BelongsToStruct != nil && !typ.IsEnumeration {
		whereValues[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = NewCodeWrapper(jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers && !typ.IsEnumeration {
		whereValues[fmt.Sprintf("%s.belongs_to_account", tableName)] = NewCodeWrapper(jen.ID(buildFakeVarName("Account")).Dot("ID"))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierRetrievalQueryMethodConditionalClauses(p *Project) []jen.Code {
	return typ.buildDBQuerierSingleInstanceQueryMethodConditionalClauses(p)
}

func (typ DataType) BuildDBQuerierListRetrievalQueryMethodConditionalClauses(p *Project) []jen.Code {
	n := typ.Name
	puvn := n.PluralUnexportedVarName()

	whereValues := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", typ.Name.PluralUnexportedVarName()), jen.ID("archivedOnColumn")).MapAssign().Nil(),
	}

	for _, pt := range p.FindOwnerTypeChain(typ) {
		whereValues = append(
			whereValues,
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
				jen.ID("idColumn"),
			).MapAssign().IDf("%sID", pt.Name.UnexportedVarName()),
		)

		if pt.BelongsToAccount && pt.RestrictedToAccountMembers {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.IDf("%sTableName", pt.Name.PluralUnexportedVarName()),
					jen.IDf("%sAccountOwnershipColumn", pt.Name.PluralUnexportedVarName()),
				).MapAssign().ID("accountID"),
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
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sAccountOwnershipColumn", puvn)).MapAssign().ID("accountID"))
	}

	return whereValues
}

func (typ DataType) BuildDBQuerierExistenceMethodParams(p *Project) []jen.Code {
	return typ.buildGetSomethingParams(p, true)
}

func (typ DataType) buildGetSomethingArgs(p *Project) []jen.Code {
	params := []jen.Code{ctxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.ID("accountID"))
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

	if typ.BelongsToAccount {
		params = append(params, jen.ID("accountID"), jen.ID("archivedBy"))
	}

	return params
}

func (typ DataType) BuildDBClientExistenceMethodCallArgs(p *Project) []jen.Code {
	return typ.buildGetSomethingArgs(p)
}

func (typ DataType) BuildDBClientRetrievalMethodCallArgs(p *Project) []jen.Code {
	return typ.buildGetSomethingArgs(p)
}

func (typ DataType) BuildDBClientArchiveMethodCallArgs() []jen.Code {
	return typ.buildArchiveSomethingArgs()
}

func (typ DataType) BuildDBQuerierExistenceQueryBuildingArgs(p *Project) []jen.Code {
	params := typ.buildGetSomethingArgs(p)

	return params[1:]
}

func (typ DataType) BuildDBQuerierRetrievalQueryBuildingArgs(p *Project) []jen.Code {
	params := typ.buildGetSomethingArgs(p)

	return params[1:]
}

func (typ DataType) BuildDBQuerierArchiveQueryBuildingArgs() []jen.Code {
	params := typ.buildArchiveSomethingArgs()

	return params[1:]
}

func (typ DataType) BuildInterfaceDefinitionExistenceMethodCallArgs(p *Project) []jen.Code {
	return typ.buildGetSomethingArgs(p)
}

func (typ DataType) BuildInterfaceDefinitionRetrievalMethodCallArgs(p *Project) []jen.Code {
	return typ.buildGetSomethingArgs(p)
}

func (typ DataType) BuildInterfaceDefinitionArchiveMethodCallArgs() []jen.Code {
	return typ.buildArchiveSomethingArgs()
}

func (typ DataType) buildGetSomethingArgsWithExampleVariables(p *Project, includeCtx bool) []jen.Code {
	params := []jen.Code{}

	if includeCtx {
		params = append(params, ctxVar())
	}

	owners := p.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for _, pt := range owners {
		pts := pt.Name.Singular()
		params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	return params
}

func (typ DataType) BuildHTTPClientRetrievalTestCallArgs(p *Project) []jen.Code {
	return typ.buildGetSomethingArgsWithExampleVariables(p, true)
}

func (typ DataType) buildSingleInstanceQueryTestCallArgs(p *Project) []jen.Code {
	params := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for _, pt := range owners {
		pts := pt.Name.Singular()
		params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		params = append(params, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return params
}

func (typ DataType) buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(p *Project) []jen.Code {
	owners := p.FindOwnerTypeChain(typ)
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

	if typ.RestrictedToUserAtSomeLevel(p) {
		args = append(args, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return args
}

func (typ DataType) BuildArgsForDBQuerierExistenceMethodTest(p *Project) []jen.Code {
	params := typ.buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(p)

	return params
}

func (typ DataType) BuildArgsForDBQuerierRetrievalMethodTest(p *Project) []jen.Code {
	params := typ.buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(p)

	return params
}

func (typ DataType) BuildArgsForServiceRouteExistenceCheck(p *Project) []jen.Code {
	owners := p.FindOwnerTypeChain(typ)
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

	if typ.RestrictedToUserAtSomeLevel(p) {
		args = append(args, jen.ID("accountID"))
	}

	return args
}

func (typ DataType) buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p *Project) []jen.Code {
	params := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for _, pt := range owners {
		pts := pt.Name.Singular()
		params = append(params, jen.ID(buildFakeVarName(pts)).Dot("ID"))
	}
	params = append(params, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return params
}

func (typ DataType) BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(p *Project) []jen.Code {
	return typ.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p)
}

func (typ DataType) BuildDBQuerierRetrievalQueryTestCallArgs(p *Project) []jen.Code {
	return typ.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p)
}

func (typ DataType) BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(p *Project) []jen.Code {
	lines := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	if typ.RestrictedToUserAtSomeLevel(p) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeUser").Call())
	}

	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())

		if pt.BelongsToAccount && typ.RestrictedToAccountMembers {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}

	lines = append(lines, jen.ID(constants.FilterVarName).Assign().Qual(p.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call())

	return lines
}

func (typ DataType) BuildDBQuerierCreateSomethingQueryBuilderTestPreQueryLines(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDBQuerierUpdateSomethingQueryBuilderTestPreQueryLines(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildGetSomethingLogValues(p *Project) jen.Code {
	params := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.Litf("%s_id", pt.Name.RouteName()).Op(":").IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.Litf("%s_id", typ.Name.RouteName()).Op(":").IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.Lit("user_id").Op(":").ID("accountID"))
	}

	return jen.Map(jen.ID("string")).Interface().Valuesln(params...)
}

func (typ DataType) BuildGetListOfSomethingLogValues(p *Project) *jen.Statement {
	params := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.Litf("%s_id", pt.Name.RouteName()).Op(":").IDf("%sID", pt.Name.UnexportedVarName()))
	}

	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.Lit("user_id").Op(":").ID("accountID"))
	}

	if len(params) > 0 {
		return jen.Map(jen.ID("string")).Interface().Valuesln(params...)
	}

	return nil
}

func (typ DataType) BuildGetListOfSomethingFromIDsParams(p *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
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

func (typ DataType) BuildGetListOfSomethingFromIDsQueryBuilderParams(p *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	params = append(params,
		jen.ID("limit").Uint8(),
		jen.ID("ids").Index().Uint64(),
		jen.ID("forAdmin").Bool(),
	)

	return params
}

func (typ DataType) BuildGetListOfSomethingFromIDsArgs(p *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.ID("accountID"))
	}

	params = append(params,
		jen.ID("limit"),
		jen.ID("ids"),
	)

	return params
}

func (typ DataType) BuildGetListOfSomethingFromIDsArgsForTest(p *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("example%s", pt.Name.Singular()).Dot("ID"))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.ID("exampleUser").Dot("ID"))
	}

	params = append(params,
		jen.ID("defaultLimit"),
		jen.ID("exampleIDs"),
	)

	return params
}

func (typ DataType) buildGetListOfSomethingParams(p *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	if !isModelsPackage {
		params = append(params, jen.ID("filter").Op("*").Qual(p.TypesPackage(), "QueryFilter"))
	} else {
		params = append(params, jen.ID("filter").Op("*").ID("QueryFilter"))
	}

	return params
}

func (typ DataType) BuildMockDataManagerListRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(p, false)
}

func (typ DataType) BuildInterfaceDefinitionListRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(p, true)
}

func (typ DataType) BuildDBClientListRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(p, false)
}

func (typ DataType) BuildDBQuerierListRetrievalMethodParams(p *Project) []jen.Code {
	return typ.buildGetListOfSomethingParams(p, false)
}

func (typ DataType) BuildDBQuerierListRetrievalQueryBuildingMethodParams(p *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	params = append(params, jen.ID("forAdmin").Bool(), jen.ID("filter").Op("*").Qual(p.TypesPackage(), "QueryFilter"))

	return params
}

const creationObjectVarName = "input"

func (typ DataType) buildCreateSomethingParams(p *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(creationObjectVarName).Op("*").IDf("%sCreationInput", sn))
	} else {
		params = append(params, jen.ID(creationObjectVarName).Op("*").Qual(p.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)))
	}

	params = append(params, jen.ID("createdByUser").Uint64())

	return params
}

func (typ DataType) BuildMockInterfaceDefinitionCreationMethodParams(p *Project) []jen.Code {
	return typ.buildCreateSomethingParams(p, false)
}

func (typ DataType) BuildInterfaceDefinitionCreationMethodParams(p *Project) []jen.Code {
	return typ.buildCreateSomethingParams(p, true)
}

func (typ DataType) BuildDBClientCreationMethodParams(p *Project) []jen.Code {
	return typ.buildCreateSomethingParams(p, false)
}

func (typ DataType) BuildDBQuerierCreationMethodParams(p *Project) []jen.Code {
	return typ.buildCreateSomethingParams(p, false)
}

func (typ DataType) BuildDBQuerierCreationQueryBuildingMethodParams(p *Project, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params, jen.ID(creationObjectVarName).Op("*").IDf("%sCreationInput", sn))
	} else {
		params = append(params, jen.ID(creationObjectVarName).Op("*").Qual(p.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)))
	}

	return params
}

func (typ DataType) buildCreateSomethingArgs() []jen.Code {
	params := []jen.Code{ctxVar(), jen.ID(creationObjectVarName), jen.ID("createdByUser")}

	return params
}

func (typ DataType) BuildMockInterfaceDefinitionCreationMethodCallArgs() []jen.Code {
	return typ.buildCreateSomethingArgs()
}

func (typ DataType) BuildDBQuerierCreationMethodQueryBuildingArgs() []jen.Code {
	params := typ.buildCreateSomethingArgs()

	return params[1:]
}

func (typ DataType) BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(p *Project) []jen.Code {
	params := []jen.Code{}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.RestrictedToUserAtSomeLevel(p) {
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

	if typ.BelongsToAccount {
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

func (typ DataType) buildUpdateSomethingParams(p *Project, updatedVarName string, isModelsPackage bool) []jen.Code {
	params := []jen.Code{ctxParam()}

	if updatedVarName == "" {
		panic("buildUpdateSomethingParams called with empty updatedVarName")
	}

	sn := typ.Name.Singular()
	if isModelsPackage {
		params = append(params,
			jen.ID(updatedVarName).Op("*").ID(sn),
			jen.ID("changedByUser").Uint64(), jen.ID("changes").Index().PointerTo().ID("FieldChangeSummary"),
		)
	} else {
		params = append(params,
			jen.ID(updatedVarName).Op("*").Qual(p.TypesPackage(), sn),
			jen.ID("changedByUser").Uint64(), jen.ID("changes").Index().PointerTo().Qual(p.TypesPackage(), "FieldChangeSummary"),
		)
	}

	return params
}

func (typ DataType) BuildDBClientUpdateMethodParams(p *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(p, updatedVarName, false)
}

func (typ DataType) BuildDBQuerierUpdateMethodParams(p *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(p, updatedVarName, false)
}

func (typ DataType) BuildDBQuerierUpdateQueryBuildingMethodParams(p *Project) []jen.Code {
	params := []jen.Code{ctxParam()}

	sn := typ.Name.Singular()

	params = append(params,
		jen.ID("input").Op("*").Qual(p.TypesPackage(), sn),
	)

	return params
}

func (typ DataType) BuildInterfaceDefinitionUpdateMethodParams(p *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(p, updatedVarName, true)
}

func (typ DataType) BuildMockDataManagerUpdateMethodParams(p *Project, updatedVarName string) []jen.Code {
	return typ.buildUpdateSomethingParams(p, updatedVarName, false)
}

func (typ DataType) buildUpdateSomethingArgsWithExampleVars(p *Project, updatedVarName string) []jen.Code {
	params := []jen.Code{ctxVar()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
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
	params := []jen.Code{ctxVar(), jen.ID(updatedVarName), jen.ID("changedByUser"), jen.ID("changes")}

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

func (typ DataType) buildGetListOfSomethingArgs(p *Project) []jen.Code {
	params := []jen.Code{ctxVar()}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.ID("accountID"))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func (typ DataType) BuildDBClientListRetrievalMethodCallArgs(p *Project) []jen.Code {
	return typ.buildGetListOfSomethingArgs(p)
}

func (typ DataType) BuildDBQuerierListRetrievalMethodArgs(p *Project) []jen.Code {
	params := typ.buildGetListOfSomethingArgs(p)

	return params[1:]
}

func (typ DataType) BuildMockDataManagerListRetrievalMethodCallArgs(p *Project) []jen.Code {
	return typ.buildGetListOfSomethingArgs(p)
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWithOwnerStruct(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	if typ.OwnedByAUserAtSomeLevel(p) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeUser").Call())
	}

	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())

		if pt.BelongsToAccount {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildCreationRequestMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDependentObjectsForDBQueriersExistenceMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)
}

func (typ DataType) BuildDependentObjectsForDBQueriersCreationMethodTest(p *Project) []jen.Code {
	lines := typ.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)

	sn := typ.Name.Singular()
	lines = append(lines, jen.ID(buildFakeVarName("Input")).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.ID(buildFakeVarName(sn))))

	return lines
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}

	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())
	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	return lines
}

func (typ DataType) buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) BuildHTTPClientRetrievalMethodTestDependentObjects(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildArchiveRequestMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildExistenceRequestMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientExistenceMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildRetrievalRequestMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientRetrievalMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientArchiveMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(p)
}

func (typ DataType) buildDependentObjectsForHTTPClientListRetrievalTest(p *Project) []jen.Code {
	lines := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
	}

	return lines
}

func (typ DataType) BuildDependentObjectsForHTTPClientListRetrievalTest(p *Project) []jen.Code {
	return typ.buildDependentObjectsForHTTPClientListRetrievalTest(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildListRetrievalRequestMethodTest(p *Project) []jen.Code {
	return typ.buildDependentObjectsForHTTPClientListRetrievalTest(p)
}

func (typ DataType) buildVarDeclarationsOfDependentStructsForUpdateFunction(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()
	owners := p.FindOwnerTypeChain(typ)

	for i, pt := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			pts := pt.Name.Singular()
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call())

	return lines
}

func (typ DataType) BuildDependentObjectsForHTTPClientUpdateMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsForUpdateFunction(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientBuildUpdateRequestMethodTest(p *Project) []jen.Code {
	return typ.buildVarDeclarationsOfDependentStructsForUpdateFunction(p)
}

func (typ DataType) BuildDependentObjectsForHTTPClientCreationMethodTest(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines,
			jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call(),
		)
	}

	lines = append(lines,
		jen.ID(buildFakeVarName(sn)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
	)

	if typ.BelongsToStruct != nil {
		lines = append(lines,
			jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"),
		)
	}

	return lines
}

func (typ DataType) BuildFormatStringForHTTPClientExistenceMethodTest(p *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientRetrievalMethodTest(p *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientUpdateMethodTest(p *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientArchiveMethodTest(p *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += fmt.Sprintf("%s/", typ.Name.PluralRouteName()) + "%d"

	return modelRoute
}

func (typ DataType) BuildFormatStringForHTTPClientListMethodTest(p *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildFormatStringForHTTPClientCreateMethodTest(p *Project) (path string) {
	modelRoute := "/api/v1/"
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		modelRoute += fmt.Sprintf("%s/", pt.Name.PluralRouteName()) + "%d/"
	}
	modelRoute += typ.Name.PluralRouteName()

	return modelRoute
}

func (typ DataType) BuildFormatCallArgsForHTTPClientRetrievalMethodTest(p *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	callArgs = append(callArgs, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForHTTPClientExistenceMethodTest(p *Project) (args []jen.Code) {
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		args = append(args, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	args = append(args, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("ID"))

	return args
}

func (typ DataType) BuildFormatCallArgsForHTTPClientListMethodTest(p *Project) (args []jen.Code) {
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		args = append(args, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	return args
}

func (typ DataType) BuildFormatCallArgsForHTTPClientCreationMethodTest(p *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		callArgs = append(callArgs, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	return callArgs
}

func (typ DataType) BuildFormatCallArgsForHTTPClientUpdateTest(p *Project) (args []jen.Code) {
	callArgs := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)

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

func (typ DataType) BuildArgsForHTTPClientExistenceRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientExistenceRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientExistenceMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildArgsForHTTPClientCreateRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildArgsForHTTPClientRetrievalRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientRetrievalRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientRetrievalMethod(p *Project, call bool) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientCreateRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

	params = append(params, jen.ID("input").PointerTo().Qual(p.TypesPackage(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))

	return params
}

func (typ DataType) BuildParamsForHTTPClientCreateMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

	params = append(params, jen.ID("input").PointerTo().Qual(p.TypesPackage(), fmt.Sprintf("%sCreationInput", typ.Name.Singular())))

	return params
}

func (typ DataType) BuildParamsForHTTPClientUpdateRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

	params = append(params, jen.ID(typ.Name.UnexportedVarName()).PointerTo().Qual(p.TypesPackage(), typ.Name.Singular()))

	return params
}

func (typ DataType) BuildArgsForHTTPClientUpdateRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientUpdateMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

	params = append(params, jen.ID(typ.Name.UnexportedVarName()).PointerTo().Qual(p.TypesPackage(), typ.Name.Singular()))

	return params
}

func (typ DataType) BuildParamsForHTTPClientArchiveRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildArgsForHTTPClientArchiveRequestBuildingMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientArchiveMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) buildParamsForMethodThatHandlesAnInstanceWithStructs(p *Project) []jen.Code {
	owners := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildArgsForHTTPClientExistenceRequestBuildingMethodTest(p *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(p)
}

func (typ DataType) BuildArgsForHTTPClientExistenceMethodTest(p *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(p)
}

func (typ DataType) BuildArgsForHTTPClientRetrievalRequestBuilderMethodTest(p *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(p)
}

func (typ DataType) BuildArgsForHTTPClientArchiveRequestBuildingMethodTest(p *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(p)
}
func (typ DataType) BuildArgsForHTTPClientArchiveMethodTest(p *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(p)
}

func (typ DataType) BuildArgsForHTTPClientArchiveMethodTestURLFormatCall(p *Project) []jen.Code {
	params := typ.BuildArgsForHTTPClientArchiveMethodTest(p)

	return params[1:]
}

func (typ DataType) BuildArgsForHTTPClientMethodTest(p *Project) []jen.Code {
	return typ.buildParamsForMethodThatHandlesAnInstanceWithStructs(p)
}

func (typ DataType) BuildHTTPClientCreationRequestBuildingMethodArgsForTest(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildHTTPClientCreationMethodArgsForTest(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildArgsForHTTPClientListRequestMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildParamsForHTTPClientListRequestMethod(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		params = append(params, jen.List(listParams...).Uint64())
	}

	params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(p.TypesPackage(), "QueryFilter"))

	return params
}

func (typ DataType) BuildParamsForHTTPClientMethodThatFetchesAList(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxParam()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		params = append(params, jen.List(listParams...).Uint64())
	}

	params = append(params, jen.ID(constants.FilterVarName).PointerTo().Qual(p.TypesPackage(), "QueryFilter"))

	return params
}

func (typ DataType) BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildCallArgsForHTTPClientListRetrievalMethodTest(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) BuildCallArgsForHTTPClientUpdateMethodTest(p *Project) []jen.Code {
	parents := p.FindOwnerTypeChain(typ)
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

func (typ DataType) buildRequisiteFakeVarDecs(p *Project, createCtx bool) []jen.Code {
	lines := []jen.Code{}
	if createCtx {
		lines = append(lines, constants.CreateCtx(), jen.Newline())
	}

	if typ.OwnedByAUserAtSomeLevel(p) {
		lines = append(lines, jen.ID(buildFakeVarName("Account")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeAccount").Call())
	}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
		if pt.BelongsToAccount {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("Account")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
	}

	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("Account")).Dot("ID"))
	}
	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	return lines
}

func (typ DataType) buildRequisiteFakeVarDecForModifierFuncs(p *Project, createCtx bool) []jen.Code {
	lines := []jen.Code{}

	if createCtx {
		lines = append(lines, constants.CreateCtx(), jen.Newline())
	}
	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName("Account")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeAccount").Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("Account")).Dot("ID"))
	}
	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientExistenceMethodTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecs(p, true)
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecs(p, false)
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientCreateMethodTest(p *Project) []jen.Code {
	lines := []jen.Code{constants.CreateCtx(), jen.Newline()}

	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeUser").Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientArchiveMethodTest(p *Project) []jen.Code {
	var lines []jen.Code

	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeUser").Call())
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}

	return append([]jen.Code{
		constants.CreateCtx(),
		jen.Newline(),
		jen.Var().ID("expected").Error(),
		jen.Newline(),
	}, lines...)
}

func (typ DataType) BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()
	owners := p.FindOwnerTypeChain(typ)

	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
		if pt.BelongsToStruct != nil {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dotf("BelongsTo%s", pt.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(pt.BelongsToStruct.Singular())).Dot("ID"))
		}
		if pt.BelongsToAccount {
			lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
	}
	lines = append(lines, jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName).Equals().ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) buildRequisiteFakeVarDecsForListFunction(p *Project) []jen.Code {
	lines := []jen.Code{}

	if !(typ.BelongsToAccount && typ.RestrictedToAccountMembers) && typ.RestrictedToUserAtSomeLevel(p) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeUser").Call())
	}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarsForDBClientListRetrievalMethodTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarDecsForListFunction(p)
}

func (typ DataType) BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(p *Project, includeFilter bool) []jen.Code {
	lines := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)

	if typ.RestrictedToUserAtSomeLevel(p) {
		lines = append(lines,
			jen.ID(buildFakeVarName("User")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeUser").Call(),
		)
	}

	for _, pt := range owners {
		pts := pt.Name.Singular()
		lines = append(lines, jen.ID(buildFakeVarName(pts)).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", pts)).Call())

		if pt.BelongsToAccount && typ.RestrictedToAccountMembers {
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID"))
		}
		if pt.BelongsToStruct != nil {
			btssn := pt.BelongsToStruct.Singular()
			lines = append(lines, jen.ID(buildFakeVarName(pts)).Dotf("BelongsTo%s", btssn).Equals().ID(buildFakeVarName(btssn)).Dot("ID"))
		}
	}

	if includeFilter {
		lines = append(lines, jen.ID(constants.FilterVarName).Assign().Qual(p.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call())
	}

	return lines
}

func (typ DataType) buildRequisiteFakeVarCallArgsForCreation(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return lines
}

func (typ DataType) buildRequisiteFakeVarCallArgs(p *Project) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("BelongsToAccount"))
	} else if typ.RestrictedToUserAtSomeLevel(p) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p *Project) []jen.Code {

	lines := []jen.Code{}
	sn := typ.Name.Singular()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if (typ.BelongsToAccount && typ.RestrictedToAccountMembers) || typ.RestrictedToUserAtSomeLevel(p) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceReadHandlerTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceCreateHandlerTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceUpdateHandlerTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest() []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBClientExistenceMethodTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgs(p)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBClientRetrievalMethodTest(p *Project) []jen.Code {
	return typ.buildRequisiteFakeVarCallArgs(p)
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest() []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToAccount {
		lines = append(lines, jen.ID(buildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}

	return lines
}

func (typ DataType) BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(p *Project) []jen.Code {
	lines := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(p *Project) []jen.Code {
	lines := []jen.Code{constants.CtxVar()}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.RestrictedToUserAtSomeLevel(p) {
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

	if typ.BelongsToAccount {
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

func (typ DataType) BuildCallArgsForDBClientListRetrievalMethodTest(p *Project) []jen.Code {
	lines := []jen.Code{}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID(buildFakeVarName(pt.Name.Singular())).Dot("ID"))
	}

	if typ.RestrictedToUserAtSomeLevel(p) {
		lines = append(lines, jen.ID(buildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func (typ DataType) BuildRequisiteVarsForDBClientUpdateMethodTest(p *Project) []jen.Code {
	lines := []jen.Code{
		constants.CreateCtx(),
		jen.Var().ID("expected").Error(),
		jen.Newline(),
		func() jen.Code {
			if typ.BelongsToAccount {
				return jen.ID(buildFakeVarName("User")).Assign().Qual(p.FakeTypesPackage(), "BuildFakeUser").Call()
			}
			return jen.Null()
		}(),
		jen.ID(buildFakeVarName(typ.Name.Singular())).Assign().Qual(p.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call(),
		func() jen.Code {
			if typ.BelongsToAccount {
				return jen.ID(buildFakeVarName(typ.Name.Singular())).Dot("BelongsToAccount").Equals().ID(buildFakeVarName("User")).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.Newline(),
	}

	return lines
}
func (typ DataType) BuildCallArgsForDBClientUpdateMethodTest() []jen.Code {
	lines := []jen.Code{jen.ID(buildFakeVarName(typ.Name.Singular()))}

	return lines
}
