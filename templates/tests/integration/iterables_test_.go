package integration

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	createdVarPrefix = "created"
)

func buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i != len(parents)-1 {
				listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
			}
		}

		if len(listParams) > 0 {
			params = append(params, listParams...)
		}
	}
	params = append(params, jen.IDf("created%s", typ.Name.Singular()))

	return params
}

func buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructsFor404Tests(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i != len(parents)-1 {
				listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
			}
		}

		if len(listParams) > 0 {
			params = append(params, listParams...)
		}
	}
	params = append(params, jen.IDf("example%s", typ.Name.Singular()))

	return params
}

func buildCreationArguments(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	args := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for i, ot := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			args = append(args, jen.IDf("%s%s", varPrefix, ot.Name.Singular()).Dot("ID"))
		}
	}

	return args
}

func buildListArguments(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	args := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for _, ot := range owners {
		args = append(args, jen.IDf("%s%s", varPrefix, ot.Name.Singular()).Dot("ID"))
	}

	args = append(args, jen.Nil())

	return args
}

func buildRequisiteCreationCode(proj *models.Project, typ models.DataType, includeSelf bool) (lines []jen.Code) {
	for _, ot := range proj.FindOwnerTypeChain(typ) {
		ots := ot.Name.Singular()

		creationArgs := append(buildCreationArguments(proj, createdVarPrefix, ot), jen.IDf("example%sInput", ots))

		lines = append(lines,
			jen.Commentf("Create %s.", ot.Name.SingularCommonName()),
			utils.BuildFakeVar(proj, ots),
			func() jen.Code {
				if ot.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(ots)).Dotf("BelongsTo%s", ot.BelongsToStruct.Singular()).Equals().IDf("%s%s", createdVarPrefix, ot.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			utils.BuildFakeVarWithCustomName(
				proj,
				fmt.Sprintf("example%sInput", ots),
				fmt.Sprintf("BuildFake%sCreationInputFrom%s", ots, ots),
				jen.ID(utils.BuildFakeVarName(ots)),
			),
			jen.List(jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()).Assign().ID("testClients").Dot("main").Dotf("Create%s", ots).Call(
				creationArgs...,
			),
			jen.ID("requireNotNilAndNoProblems").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()),
			jen.Newline(),
		)
	}

	sn := typ.Name.Singular()
	creationArgs := append(buildCreationArguments(proj, "created", typ), jen.IDf("example%sInput", sn))

	if includeSelf {
		lines = append(lines,
			jen.Commentf("Create %s.", typ.Name.SingularCommonName()),
			utils.BuildFakeVar(proj, sn),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%s%s", createdVarPrefix, typ.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			utils.BuildFakeVarWithCustomName(
				proj,
				fmt.Sprintf("example%sInput", sn),
				fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
				jen.ID(utils.BuildFakeVarName(sn)),
			),
			jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
				creationArgs...,
			),
			jen.ID("requireNotNilAndNoProblems").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()),
			jen.Newline(),
		)
	}

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	var listParams []jen.Code
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
		}
		listParams = append(listParams, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))
	}

	return params
}

func buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnlyFor404Tests(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	var listParams []jen.Code
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
		}
		listParams = append(listParams, jen.ID("nonexistentID"))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.ID("nonexistentID"))
	}

	return params
}

func buildRequisiteCleanupCode(proj *models.Project, typ models.DataType, includeSelf bool) []jen.Code {
	var lines []jen.Code

	parentTypes := proj.FindOwnerTypeChain(typ)
	// reverse it
	for i, j := 0, len(parentTypes)-1; i < j; i, j = i+1, j-1 {
		parentTypes[i], parentTypes[j] = parentTypes[j], parentTypes[i]
	}

	if includeSelf {
		lines = append(lines,
			jen.Newline(),
			jen.Commentf("Clean up %s.", typ.Name.SingularCommonName()),
			utils.AssertNoError(
				jen.ID("testClients").Dot("main").Dotf("Archive%s",
					typ.Name.Singular()).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
				),
				nil,
			),
		)
	}

	for _, ot := range parentTypes {
		lines = append(lines,
			jen.Newline(),
			jen.Commentf("Clean up %s.", ot.Name.SingularCommonName()),
			utils.AssertNoError(
				jen.ID("testClients").Dot("main").Dotf("Archive%s", ot.Name.Singular()).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, ot)...,
				),
				nil,
			),
		)
	}

	return lines
}

func buildGetSomethingAuditLogEntriesArgs(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{
		jen.ID("ctx"),
	}

	for _, owner := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("created%s", owner.Name.Singular()).Dot("ID"))
	}

	lines = append(lines, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))

	return lines
}

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	comparisonLines := []jen.Code{
		jen.ID("t").Dot("Helper").Call(),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "NotZero").Call(
			jen.ID("t"),
			jen.ID("actual").Dot("ID"),
		),
	}
	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			fsn := field.Name.Singular()
			comparisonLines = append(comparisonLines, jen.Qual(constants.AssertionLibrary, "Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot(fsn),
				jen.ID("actual").Dot(fsn),
				jen.Lit("expected "+fsn+" for "+scn+" #%d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot(fsn),
				jen.ID("actual").Dot(fsn),
			))
		}
	}
	comparisonLines = append(comparisonLines,
		jen.Qual(constants.AssertionLibrary, "NotZero").Call(
			jen.ID("t"),
			jen.ID("actual").Dot("CreatedOn"),
		),
	)

	code.Add(
		jen.Func().IDf("check%sEquality", sn).Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).PointerTo().Qual(proj.TypesPackage(), sn)).Body(
			comparisonLines...,
		),
		jen.Newline(),
	)

	code.Add(buildTestCreating(proj, typ)...)
	code.Add(buildTestListing(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildTestSearching(proj, typ)...)
		if typ.RestrictedToAccountMembers {
			code.Add(buildTestSearching_ReturnsOnlySomething(proj, typ)...)
		}
	}

	code.Add(buildTestExistenceChecking_ReturnsFalseForNonexistentSomething(proj, typ)...)
	code.Add(buildTestExistenceChecking_ReturnsTrueForValidSomething(proj, typ)...)
	code.Add(buildTestReading_Returns404ForNonexistentSomething(proj, typ)...)
	code.Add(buildTestReading(proj, typ)...)
	code.Add(buildTestUpdating_Returns404ForNonexistentSomething(proj, typ)...)
	code.Add(convertSomethingToSomethingUpdateInput(proj, typ)...)
	code.Add(buildTestUpdating(proj, typ)...)
	code.Add(buildTestArchiving_Returns404ForNonexistentSomething(proj, typ)...)
	code.Add(buildTestArchiving(proj, typ)...)
	code.Add(buildTestAuditing_Returns404ForNonexistentSomething(proj, typ)...)

	return code
}

func buildTestCreating(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(
			jen.ID("t"),
			jen.IDf("example%s", sn),
			jen.IDf("created%s", sn),
		),
		jen.Newline(),
		jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Assign().ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
			buildGetSomethingAuditLogEntriesArgs(proj, typ)...,
		),
		jen.Qual(constants.MustAssertPkg, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.ID("expectedAuditLogEntries").Assign().Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sCreationEvent", sn)))),
		jen.ID("validateAuditLogEntries").Call(
			jen.ID("t"),
			jen.ID("expectedAuditLogEntries"),
			jen.ID("auditLogEntries"),
			jen.IDf("created%s", sn).Dot("ID"),
			jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sAssignmentKey", sn)),
		),
		jen.Newline(),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, true)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Creating", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be creatable"),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestListing(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("create %s", pcn),
		jen.Var().ID("expected").Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.For(jen.ID("i").Assign().Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
			jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			jen.IDf("example%sInput", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
			jen.Newline(),
			jen.List(jen.IDf("created%s", sn), jen.IDf("%sCreationErr", uvn)).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
				append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))...,
			),
			jen.ID("requireNotNilAndNoProblems").Call(
				jen.ID("t"),
				jen.IDf("created%s", sn),
				jen.IDf("%sCreationErr", uvn),
			),
			jen.Newline(),
			jen.ID("expected").Equals().ID("append").Call(
				jen.ID("expected"),
				jen.IDf("created%s", sn),
			),
		),
		jen.Newline(),
		jen.Commentf("assert %s list equality", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Get%s", pn).Call(
			buildListArguments(proj, createdVarPrefix, typ)...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.ID("actual"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "True").Callln(
			jen.ID("t"),
			jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(pn)),
			jen.Lit("expected %d to be <= %d"),
			jen.ID("len").Call(jen.ID("expected")),
			jen.ID("len").Call(jen.ID("actual").Dot(pn)),
		),
		jen.Newline(),
		jen.Comment("clean up"),
		jen.For(jen.List(jen.Underscore(), jen.IDf("created%s", sn)).Assign().Range().ID("actual").Dot(pn)).Body(
			jen.Qual(constants.AssertionLibrary, "NoError").Call(
				jen.ID("t"),
				jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
				),
			),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Listing", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be readable in paginated form"),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestSearching(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	var firstStringField wordsmith.SuperPalabra
	for _, field := range typ.Fields {
		if field.Type == "string" {
			firstStringField = field.Name
			break
		}
	}

	searchArgs := []jen.Code{
		constants.CtxVar(),
	}
	for _, owner := range proj.FindOwnerTypeChain(typ) {
		searchArgs = append(searchArgs, jen.IDf("created%s", owner.Name.Singular()).Dot("ID"))
	}

	searchArgs = append(searchArgs,
		jen.ID(utils.BuildFakeVarName(sn)).Dot(firstStringField.Singular()),
		jen.ID(utils.BuildFakeVarName("Limit")),
	)

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Commentf("create %s", pcn),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.Var().ID("expected").Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.For(jen.ID("i").Assign().Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
			jen.IDf("example%sInput", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
			jen.IDf("example%sInput", sn).Dot(firstStringField.Singular()).Equals().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s %d"),
				jen.IDf("example%sInput", sn).Dot(firstStringField.Singular()),
				jen.ID("i"),
			),
			jen.Newline(),
			jen.List(jen.IDf("created%s", sn), jen.IDf("%sCreationErr", uvn)).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
				append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))...,
			),
			jen.ID("requireNotNilAndNoProblems").Call(
				jen.ID("t"),
				jen.IDf("created%s", sn),
				jen.IDf("%sCreationErr", uvn),
			),
			jen.Newline(),
			jen.ID("expected").Equals().ID("append").Call(
				jen.ID("expected"),
				jen.IDf("created%s", sn),
			),
		),
		jen.Newline(),
		jen.ID("exampleLimit").Assign().ID("uint8").Call(jen.Lit(20)),
		jen.Newline(),
		jen.Commentf("assert %s list equality", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Search%s", pn).Call(
			searchArgs...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.ID("actual"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "True").Callln(
			jen.ID("t"),
			jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual")),
			jen.Lit("expected results length %d to be <= %d"),
			jen.ID("len").Call(jen.ID("expected")),
			jen.ID("len").Call(jen.ID("actual")),
		),
		jen.Newline(),
		jen.Comment("clean up"),
		jen.For(jen.List(jen.Underscore(), jen.IDf("created%s", sn)).Assign().Range().ID("expected")).Body(
			jen.Qual(constants.AssertionLibrary, "NoError").Call(
				jen.ID("t"),
				jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
				),
			),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Searching", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should be able to be search for %s", pcn),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestSearching_ReturnsOnlySomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Commentf("create %s", pcn),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.Var().ID("expected").Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.For(jen.ID("i").Assign().Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
			jen.IDf("example%sInput", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
			jen.IDf("example%sInput", sn).Dot("Name").Equals().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s %d"),
				jen.IDf("example%sInput", sn).Dot("Name"),
				jen.ID("i"),
			),
			jen.Newline(),
			jen.List(jen.IDf("created%s", sn), jen.IDf("%sCreationErr", uvn)).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
				append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))...,
			),
			jen.ID("requireNotNilAndNoProblems").Call(
				jen.ID("t"),
				jen.IDf("created%s", sn),
				jen.IDf("%sCreationErr", uvn),
			),
			jen.Newline(),
			jen.ID("expected").Equals().ID("append").Call(
				jen.ID("expected"),
				jen.IDf("created%s", sn),
			),
		),
		jen.Newline(),
		jen.ID("exampleLimit").Assign().ID("uint8").Call(jen.Lit(20)),
		jen.Newline(),
		jen.Commentf("assert %s list equality", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Search%s", pn).Call(
			jen.ID("ctx"),
			jen.IDf("example%s", sn).Dot("Name"),
			jen.ID("exampleLimit"),
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.ID("actual"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "True").Callln(
			jen.ID("t"),
			jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual")),
			jen.Lit("expected results length %d to be <= %d"),
			jen.ID("len").Call(jen.ID("expected")),
			jen.ID("len").Call(jen.ID("actual")),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Comment("clean up"),
		jen.For(jen.List(jen.Underscore(), jen.IDf("created%s", sn)).Assign().Range().ID("expected")).Body(
			jen.Qual(constants.AssertionLibrary, "NoError").Call(
				jen.ID("t"),
				jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
				),
			),
		),
	)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Searching_ReturnsOnly%sThatBelongToYou", pn, pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should only receive your own %s", pcn),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestExistenceChecking_ReturnsFalseForNonexistentSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("%sExists", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnlyFor404Tests(proj, typ)...,
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "False").Call(
			jen.ID("t"),
			jen.ID("actual"),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_ExistenceChecking_ReturnsFalseForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should not return an error for nonexistent %s", scn),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestExistenceChecking_ReturnsTrueForValidSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("create %s", scn),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.IDf("example%sInput", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.List(jen.IDf("created%s", sn), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
			append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.IDf("created%s", sn),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Commentf("retrieve %s", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("%sExists", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "True").Call(
			jen.ID("t"),
			jen.ID("actual"),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("clean up %s", scn),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
			),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_ExistenceChecking_ReturnsTrueForValid%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should not return an error for existent %s", scn),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestReading_Returns404ForNonexistentSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	op := jen.Assign()
	if len(proj.FindOwnerTypeChain(typ)) > 0 {
		op = jen.Equals()
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.List(jen.Underscore(), jen.ID("err")).Add(op).ID("testClients").Dot("main").Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnlyFor404Tests(proj, typ)...,
		),
		jen.Qual(constants.AssertionLibrary, "Error").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Reading_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("it should return an error when trying to read %s that does not exist", scnwp),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestReading(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("create %s", scn),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.IDf("example%sInput", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.List(jen.IDf("created%s", sn), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
			append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.IDf("created%s", sn),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Commentf("retrieve %s", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.ID("actual"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Commentf("assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(
			jen.ID("t"),
			jen.IDf("example%s", sn),
			jen.ID("actual"),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("clean up %s", scn),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
			),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Reading", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should be readable"),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestUpdating_Returns404ForNonexistentSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		jen.IDf("example%s", sn).Dot("ID").Equals().ID("nonexistentID"),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Error").Call(
			jen.ID("t"),
			jen.ID("testClients").Dot("main").Dotf("Update%s", sn).Call(
				buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructsFor404Tests(proj, typ)...,
			),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Updating_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to update something that does not exist"),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func convertSomethingToSomethingUpdateInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	updateInputLines := []jen.Code{}
	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			fsn := field.Name.Singular()
			updateInputLines = append(updateInputLines, jen.ID(fsn).Op(":").ID("x").Dot(fsn))
		}
	}

	lines := []jen.Code{
		jen.Commentf("convert%sTo%sUpdateInput creates an %sUpdateInput struct from %s.", sn, sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("convert%sTo%sUpdateInput", sn, sn).Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn))).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)).Valuesln(
				updateInputLines...,
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestUpdating(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("create %s", scn),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.IDf("example%sInput", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.List(jen.IDf("created%s", sn), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
			append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.IDf("created%s", sn),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Commentf("change %s", scn),
		jen.IDf("created%s", sn).Dot("Update").Call(jen.IDf("convert%sTo%sUpdateInput", sn, sn).Call(jen.IDf("example%s", sn))),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("testClients").Dot("main").Dotf("Update%s", sn).Call(
				buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)...,
			),
		),
		jen.Newline(),
		jen.Commentf("retrieve changed %s", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.ID("actual"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Commentf("assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(
			jen.ID("t"),
			jen.IDf("example%s", sn),
			jen.ID("actual"),
		),
		jen.Qual(constants.AssertionLibrary, "NotNil").Call(
			jen.ID("t"),
			jen.ID("actual").Dot("LastUpdatedOn"),
		),
		jen.Newline(),
		jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Assign().ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
			buildGetSomethingAuditLogEntriesArgs(proj, typ)...,
		),
		jen.Qual(constants.MustAssertPkg, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.ID("expectedAuditLogEntries").Assign().Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(
			jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sCreationEvent", sn))),
			jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sUpdateEvent", sn))),
		),
		jen.ID("validateAuditLogEntries").Call(
			jen.ID("t"),
			jen.ID("expectedAuditLogEntries"),
			jen.ID("auditLogEntries"),
			jen.IDf("created%s", sn).Dot("ID"),
			jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sAssignmentKey", sn)),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("clean up %s", scn),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
			),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Updating", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("it should be possible to update %s", scnwp),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestArchiving_Returns404ForNonexistentSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Error").Call(
			jen.ID("t"),
			jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnlyFor404Tests(proj, typ)...,
			),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Archiving_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to delete something that does not exist"),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestArchiving(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("create %s", scn),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.IDf("example%sInput", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.List(jen.IDf("created%s", sn), jen.ID("err")).Assign().ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
			append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))...,
		),
		jen.ID("requireNotNilAndNoProblems").Call(
			jen.ID("t"),
			jen.IDf("created%s", sn),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Commentf("clean up %s", scn),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructIDsOnly(proj, typ)...,
			),
		),
		jen.Newline(),
		jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Assign().ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
			buildGetSomethingAuditLogEntriesArgs(proj, typ)...,
		),
		jen.Qual(constants.MustAssertPkg, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.ID("expectedAuditLogEntries").Assign().Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(
			jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sCreationEvent", sn))),
			jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sArchiveEvent", sn))),
		),
		jen.ID("validateAuditLogEntries").Call(
			jen.ID("t"),
			jen.ID("expectedAuditLogEntries"),
			jen.ID("auditLogEntries"),
			jen.IDf("created%s", sn).Dot("ID"),
			jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sAssignmentKey", sn)),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Archiving", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("it should be possible to delete %s", scnwp),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}

func buildTestAuditing_Returns404ForNonexistentSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	bodyLines := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("t").Dot("Name").Call(),
		),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildRequisiteCreationCode(proj, typ, false)...)

	callArgs := []jen.Code{
		constants.CtxVar(),
	}
	for _, owner := range proj.FindOwnerTypeChain(typ) {
		callArgs = append(callArgs, jen.IDf("created%s", owner.Name.Singular()).Dot("ID"))
	}
	callArgs = append(callArgs, jen.ID("nonexistentID"))

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.List(jen.ID("x"), jen.ID("err")).Assign().ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
			callArgs...,
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "Empty").Call(
			jen.ID("t"),
			jen.ID("x"),
		),
	)

	bodyLines = append(bodyLines, buildRequisiteCleanupCode(proj, typ, false)...)

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("TestSuite")).IDf("Test%s_Auditing_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to audit something that does not exist"),
				jen.Func().Params(jen.ID("testClients").PointerTo().ID("testClientWrapper")).Params(jen.Func().Params()).Body(jen.Return().Func().Params().Body(
					bodyLines...,
				)),
			)),
		jen.Newline(),
	}

	return lines
}
