package integration

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	createdVarPrefix = "created"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildCheckSomethingEquality(proj, typ)...)
	code.Add(buildTestSomething(proj, typ)...)

	return code
}

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

func buildCreationArguments(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	creationArgs := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for i, ot := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else {
			creationArgs = append(creationArgs, jen.IDf("%s%s", varPrefix, ot.Name.Singular()).Dot("ID"))
		}
	}

	return creationArgs
}

func buildListArguments(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	creationArgs := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for _, ot := range owners {
		creationArgs = append(creationArgs, jen.IDf("%s%s", varPrefix, ot.Name.Singular()).Dot("ID"))

	}

	return creationArgs
}

func buildCheckSomethingEquality(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Func().IDf("check%sEquality", sn).Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).PointerTo().Qual(proj.TypesPackage(), sn)).Body(
			buildEqualityCheckLines(typ)...,
		),
		jen.Line(),
	}

	return lines
}

func buildTestSomething(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.Func().IDf("Test%s", pn).Params(jen.ID("test").PointerTo().Qual("testing", "T")).Body(
			buildTestCreating(proj, typ),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext("should be able to be read in a list", buildTestListing(proj, typ)...),
			)),
			jen.Line(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("test").Dot("Run").Call(jen.Lit("Searching"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
						utils.BuildSubTestWithoutContext(
							fmt.Sprintf("should be able to be search for %s", pcn),
							buildTestSearching(proj, typ)...,
						),
						jen.Line(),
						utils.BuildSubTestWithoutContext(
							fmt.Sprintf("should only receive your own %s", pcn),
							buildTestSearchingForOnlyYourOwnItems(proj, typ)...,
						),
					))
				}
				return jen.Null()
			}(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("ExistenceChecking"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"it should return false with no error when checking something that does not exist",
					buildTestExistenceCheckingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, typ)...,
				),
				jen.Line(),
				utils.BuildSubTestWithoutContext(
					fmt.Sprintf("it should return true with no error when the relevant %s exists", scn),
					buildTestExistenceCheckingShouldBeReadable(proj, typ)...,
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				utils.BuildSubTestWithoutContext(
					"it should return an error when trying to read something that does not exist",
					buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, typ)...,
				),
				jen.Line(),
				utils.BuildSubTestWithoutContext("it should be readable", buildTestReadingShouldBeReadable(proj, typ)...),
			)),
			jen.Line(),
			buildTestUpdating(proj, typ),
			jen.Line(),
			buildTestDeleting(proj, typ),
		),
	}

	return lines
}

func buildRequisiteCreationCode(proj *models.Project, typ models.DataType) (lines []jen.Code) {
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
			jen.List(jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", ots).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()),
			jen.Line(),
		)
	}

	sn := typ.Name.Singular()
	creationArgs := append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))

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
		jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", sn).Call(
			creationArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()),
		jen.Line(),
	)

	return lines
}

func buildRequisiteCreationCodeWithoutType(proj *models.Project, typ models.DataType) []jen.Code {
	var lines []jen.Code

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
			jen.List(jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", ots).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()),
			jen.Line(),
		)
	}

	return lines
}

func buildRequisiteCleanupCode(proj *models.Project, typ models.DataType, includeSelf bool) []jen.Code {
	var lines []jen.Code

	parentTypes := proj.FindOwnerTypeChain(typ)
	// reverse it
	for i, j := 0, len(parentTypes)-1; i < j; i, j = i+1, j-1 {
		parentTypes[i], parentTypes[j] = parentTypes[j], parentTypes[i]
	}

	for _, ot := range parentTypes {
		lines = append(lines,
			jen.Line(),
			jen.Commentf("Clean up %s.", ot.Name.SingularCommonName()),
			utils.AssertNoError(
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", ot.Name.Singular()).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, ot)...,
				),
				nil,
			),
		)
	}

	if includeSelf {
		lines = append(lines,
			jen.Line(),
			jen.Commentf("Clean up %s.", typ.Name.SingularCommonName()),
			utils.AssertNoError(
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s",
					typ.Name.Singular()).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
				),
				nil,
			),
		)
	}

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj *models.Project, typ models.DataType) []jen.Code {
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

func buildEqualityCheckLines(typ models.DataType) []jen.Code {
	lines := []jen.Code{
		jen.ID("t").Dot("Helper").Call(),
		jen.Line(),
		utils.AssertNotZero(jen.ID("actual").Dot("ID"), nil),
	}

	for _, field := range typ.Fields {
		sn := field.Name.Singular()
		if !field.Pointer {
			lines = append(lines, utils.AssertEqual(
				jen.ID("expected").Dot(sn),
				jen.ID("actual").Dot(sn),
				jen.Lit("expected "+sn+" for ID %d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot(sn),
				jen.ID("actual").Dot(sn),
			))
		} else {
			lines = append(lines, utils.AssertEqual(
				jen.PointerTo().ID("expected").Dot(sn),
				jen.PointerTo().ID("actual").Dot(sn),
				jen.Lit("expected "+sn+" to be %v, but it was %v "),
				jen.ID("expected").Dot(sn),
				jen.ID("actual").Dot(sn),
			))
		}
	}
	lines = append(lines, utils.AssertNotZero(jen.ID("actual").Dot("CreatedOn"), nil))

	return lines
}

func buildRequisiteCreationCodeFor404Tests(proj *models.Project, typ models.DataType, indexToStop int) (lines []jen.Code) {
	for i, ot := range proj.FindOwnerTypeChain(typ) {
		if i >= indexToStop {
			break
		}

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
			jen.List(jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", ots).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()),
			jen.Line(),
		)
	}

	return lines
}

func buildRequisiteCleanupCodeFor404s(proj *models.Project, typ models.DataType, indexToStop int) (lines []jen.Code) {
	var typesToCleanup []models.DataType

	for i, ot := range proj.FindOwnerTypeChain(typ) {
		if i >= indexToStop {
			break
		}
		typesToCleanup = append(typesToCleanup, ot)
	}

	// reverse it
	for i, j := 0, len(typesToCleanup)-1; i < j; i, j = i+1, j-1 {
		typesToCleanup[i], typesToCleanup[j] = typesToCleanup[j], typesToCleanup[i]
	}

	for _, t := range typesToCleanup {
		lines = append(lines,
			jen.Line(),
			jen.Commentf("Clean up %s.", t.Name.SingularCommonName()),
			utils.AssertNoError(
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", t.Name.Singular()).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, t)...,
				),
				nil,
			),
		)
	}

	return lines
}

func buildCreationArgumentsFor404s(proj *models.Project, varPrefix string, typ models.DataType, indexToStop int) []jen.Code {
	creationArgs := []jen.Code{constants.CtxVar()}

	owners := proj.FindOwnerTypeChain(typ)
	for i, ot := range owners {
		if i == len(owners)-1 && typ.BelongsToStruct != nil {
			continue
		} else if i >= indexToStop {
			creationArgs = append(creationArgs, jen.ID("nonexistentID"))
		} else {
			creationArgs = append(creationArgs, jen.IDf("%s%s", varPrefix, ot.Name.Singular()).Dot("ID"))
		}
	}

	return creationArgs
}

func buildSubtestsForCreation404Tests(proj *models.Project, typ models.DataType) []jen.Code {
	var subtests []jen.Code
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	for i, ot := range proj.FindOwnerTypeChain(typ) {
		lines := append(
			[]jen.Code{
				utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
			},
			buildRequisiteCreationCodeFor404Tests(proj, ot, i)...,
		)

		creationArgs := append(buildCreationArgumentsFor404s(proj, createdVarPrefix, typ, i), jen.IDf("example%sInput", sn))
		lines = append(lines,
			jen.Line(),
			jen.Commentf("Create %s.", scn),
			utils.BuildFakeVar(proj, sn),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID("nonexistentID")
				}
				return jen.Null()
			}(),
			utils.BuildFakeVarWithCustomName(
				proj,
				fmt.Sprintf("example%sInput", sn),
				fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
				jen.ID(utils.BuildFakeVarName(sn)),
			),
			jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", sn).Call(
				creationArgs...,
			),
			jen.Line(),
			utils.AssertNil(jen.IDf("%s%s", createdVarPrefix, sn), nil),
			utils.AssertError(jen.Err(), nil),
		)

		lines = append(lines, buildRequisiteCleanupCodeFor404s(proj, typ, i)...)

		subtests = append(subtests,
			jen.Line(),
			utils.BuildSubTestWithoutContext(fmt.Sprintf("should fail to create for nonexistent %s", ot.Name.SingularCommonName()),
				lines...,
			),
		)
	}

	return subtests
}

func buildTestCreating(proj *models.Project, typ models.DataType) jen.Code {

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
		jen.Line(),
	}

	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines,
		jen.Commentf("Assert %s equality.", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.IDf("created%s", typ.Name.Singular())),
		jen.Line(),
		jen.Comment("Clean up."),
		jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), nil),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual")),
		utils.AssertNotNil(jen.ID("actual").Dot("ArchivedOn"), nil),
		utils.AssertNotZero(jen.ID("actual").Dot("ArchivedOn"), nil),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	testBlock := append(
		[]jen.Code{utils.BuildSubTestWithoutContext("should be createable", lines...)},
		buildSubtestsForCreation404Tests(
			proj,
			typ,
		)...,
	)

	return jen.ID("test").Dot("Run").Call(
		jen.Lit("Creating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			testBlock...,
		),
	)
}

func buildTestListing(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCodeWithoutType(proj, typ)...)

	listArgs := append(buildListArguments(proj, createdVarPrefix, typ), jen.Nil())
	creationArgs := append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))

	lines = append(lines,
		jen.Commentf("Create %s.", pcn),
		jen.Var().ID("expected").Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Body(
			jen.Commentf("Create %s.", scn),
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
			jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", sn).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)),
			jen.Line(),
			utils.AppendItemsToList(jen.ID("expected"), jen.IDf("%s%s", createdVarPrefix, sn)),
		),
		jen.Line(),
		jen.Commentf("Assert %s list equality.", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Get%s", pn).Call(
			listArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Qual(constants.AssertPkg, "True").Callln(
			jen.ID("t"),
			jen.Len(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(pn)),
			jen.Lit("expected %d to be <= %d"), jen.Len(jen.ID("expected")),
			jen.Len(jen.ID("actual").Dot(pn)),
		),
		jen.Line(),
		jen.Comment("Clean up."),
		jen.For(jen.List(jen.Underscore(), jen.IDf("%s%s", createdVarPrefix, sn)).Assign().Range().ID("actual").Dot(pn)).Body(
			jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			utils.AssertNoError(jen.Err(), nil),
		),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	return lines
}

func buildTestSearching(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCodeWithoutType(proj, typ)...)

	creationArgs := append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))

	var firstStringField wordsmith.SuperPalabra
	for _, field := range typ.Fields {
		if field.Type == "string" {
			firstStringField = field.Name
			break
		}
	}
	searchArgs := []jen.Code{
		constants.CtxVar(),
		jen.ID(utils.BuildFakeVarName(sn)).Dot(firstStringField.Singular()),
		jen.ID(utils.BuildFakeVarName("Limit")),
	}

	lines = append(lines,
		jen.Commentf("Create %s.", pcn),
		utils.BuildFakeVar(proj, sn),
		jen.Var().ID("expected").Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Body(
			jen.Commentf("Create %s.", scn),
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
			jen.IDf("example%sInput", sn).Dot(firstStringField.Singular()).Equals().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s %d"),
				jen.IDf("example%sInput", sn).Dot(firstStringField.Singular()),
				jen.ID("i"),
			),
			jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", sn).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)),
			jen.Line(),
			utils.AppendItemsToList(jen.ID("expected"), jen.IDf("%s%s", createdVarPrefix, sn)),
		),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName("Limit")).Assign().Uint8().Call(jen.Lit(20)),
		jen.Line(),
		jen.Commentf("Assert %s list equality.", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Search%s", pn).Call(
			searchArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Qual(constants.AssertPkg, "True").Callln(
			jen.ID("t"),
			jen.Len(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual")),
			jen.Lit("expected results length %d to be <= %d"), jen.Len(jen.ID("expected")),
			jen.Len(jen.ID("actual")),
		),
		jen.Line(),
		jen.Comment("Clean up."),
		jen.For(jen.List(jen.Underscore(), jen.IDf("%s%s", createdVarPrefix, sn)).Assign().Range().ID("expected")).Body(
			jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			utils.AssertNoError(jen.Err(), nil),
		),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	return lines
}

func buildTestSearchingForOnlyYourOwnItems(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCodeWithoutType(proj, typ)...)

	var firstStringField wordsmith.SuperPalabra
	for _, field := range typ.Fields {
		if field.Type == "string" {
			firstStringField = field.Name
			break
		}
	}
	searchArgs := []jen.Code{
		constants.CtxVar(),
		jen.ID("query"),
		jen.ID(utils.BuildFakeVarName("Limit")),
	}

	lines = append(lines,
		jen.Comment("create user and oauth2 client A."),
		jen.List(jen.ID("userA"), jen.Err()).Assign().Qual(proj.TestUtilPackage(), "CreateObligatoryUser").Call(
			jen.ID("urlToUse"),
			jen.ID("debug"),
		),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.List(jen.ID("ca"), jen.Err()).Assign().Qual(proj.TestUtilPackage(), "CreateObligatoryClient").Call(
			jen.ID("urlToUse"),
			jen.ID("userA"),
		),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.List(jen.ID("clientA"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
			constants.CtxVar(),
			jen.ID("ca").Dot("ClientID"),
			jen.ID("ca").Dot("ClientSecret"),
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
			jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("buildHTTPClient").Call(),
			jen.ID("ca").Dot("Scopes"),
			jen.True(),
		),
		jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientA"), jen.Err()),
		jen.Line(),
		jen.Commentf("Create %s for user A.", pcn),
		jen.IDf("example%sA", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		jen.Var().ID("createdForA").Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Body(
			jen.Commentf("Create %s.", scn),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%s%s", createdVarPrefix, typ.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			utils.BuildFakeVarWithCustomName(
				proj,
				fmt.Sprintf("example%sInputA", sn),
				fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
				jen.IDf("example%sA", sn),
			),
			jen.IDf("example%sInputA", sn).Dot(firstStringField.Singular()).Equals().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s %d"),
				jen.IDf("example%sInputA", sn).Dot(firstStringField.Singular()),
				jen.ID("i"),
			),
			jen.Line(),
			jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)).Assign().ID("clientA").Dotf("Create%s", sn).Call(
				append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInputA", sn))...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)),
			jen.Line(),
			utils.AppendItemsToList(jen.ID("createdForA"), jen.IDf("%s%s", createdVarPrefix, sn)),
		),
		jen.Line(),
		jen.ID(utils.BuildFakeVarName("Limit")).Assign().Uint8().Call(jen.Lit(20)),
		jen.ID("query").Assign().IDf("example%sA", sn).Dot(firstStringField.Singular()),
		jen.Line(),
		jen.Comment("create user and oauth2 client B."),
		jen.List(jen.ID("userB"), jen.Err()).Assign().Qual(proj.TestUtilPackage(), "CreateObligatoryUser").Call(
			jen.ID("urlToUse"),
			jen.ID("debug"),
		),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.List(jen.ID("cb"), jen.Err()).Assign().Qual(proj.TestUtilPackage(), "CreateObligatoryClient").Call(
			jen.ID("urlToUse"),
			jen.ID("userB"),
		),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.List(jen.ID("clientB"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
			constants.CtxVar(),
			jen.ID("cb").Dot("ClientID"),
			jen.ID("cb").Dot("ClientSecret"),
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
			jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
			jen.ID("buildHTTPClient").Call(),
			jen.ID("cb").Dot("Scopes"),
			jen.True(),
		),
		jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientB"), jen.Err()),
		jen.Line(),
		jen.Commentf("Create %s for user B.", pcn),
		jen.IDf("example%sB", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		jen.IDf("example%sB", sn).Dot(firstStringField.Singular()).Equals().ID("reverse").Call(jen.IDf("example%sA", sn).Dot(firstStringField.Singular())),
		jen.Var().ID("createdForB").Index().PointerTo().Qual(proj.TypesPackage(), sn),
		jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Body(
			jen.Commentf("Create %s.", scn),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%s%s", createdVarPrefix, typ.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			utils.BuildFakeVarWithCustomName(
				proj,
				fmt.Sprintf("example%sInputB", sn),
				fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
				jen.IDf("example%sB", sn),
			),
			jen.IDf("example%sInputB", sn).Dot(firstStringField.Singular()).Equals().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s %d"),
				jen.IDf("example%sInputB", sn).Dot(firstStringField.Singular()),
				jen.ID("i"),
			),
			jen.Line(),
			jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)).Assign().ID("clientB").Dotf("Create%s", sn).Call(
				append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInputB", sn))...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.IDf("%sCreationErr", uvn)),
			jen.Line(),
			utils.AppendItemsToList(jen.ID("createdForB"), jen.IDf("%s%s", createdVarPrefix, sn)),
		),
		jen.Line(),
		jen.ID("expected").Assign().ID("createdForA"),
		jen.Line(),
		jen.Commentf("Assert %s list equality.", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("clientA").Dotf("Search%s", pn).Call(
			searchArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Qual(constants.AssertPkg, "True").Callln(
			jen.ID("t"),
			jen.Len(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual")),
			jen.Lit("expected results length %d to be <= %d"), jen.Len(jen.ID("expected")),
			jen.Len(jen.ID("actual")),
		),
		jen.Line(),
		jen.Comment("Clean up."),
		jen.For(jen.List(jen.Underscore(), jen.IDf("%s%s", createdVarPrefix, sn)).Assign().Range().ID("createdForA")).Body(
			jen.Err().Equals().ID("clientA").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			utils.AssertNoError(jen.Err(), nil),
		),
		jen.Line(),
		jen.For(jen.List(jen.Underscore(), jen.IDf("%s%s", createdVarPrefix, sn)).Assign().Range().ID("createdForB")).Body(
			jen.Err().Equals().ID("clientB").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			utils.AssertNoError(jen.Err(), nil),
		),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	return lines
}

func buildTestExistenceCheckingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCodeWithoutType(proj, typ)...)

	args := buildParamsForCheckingATypeThatDoesNotExistButIncludesPredecessorID(proj, typ)

	lines = append(lines,
		jen.Commentf("Attempt to fetch nonexistent %s.", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("%sExists", sn).Call(
			args...,
		),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertFalse(jen.ID("actual"), nil),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	return lines
}

func buildTestExistenceCheckingShouldBeReadable(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}

	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines,
		jen.Line(),
		jen.Commentf("Fetch %s.", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("%sExists", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertTrue(jen.ID("actual"), nil),
		jen.Line(),
		jen.Commentf("Clean up %s.", typ.Name.SingularCommonName()),
		utils.AssertNoError(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", typ.Name.Singular()).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			nil,
		),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	return lines
}

func buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCodeWithoutType(proj, typ)...)

	args := buildParamsForCheckingATypeThatDoesNotExistButIncludesPredecessorID(proj, typ)

	lines = append(lines,
		jen.Commentf("Attempt to fetch nonexistent %s.", scn),
		jen.List(jen.Underscore(), jen.Err()).Op(func() string {
			if typ.BelongsToStruct == nil {
				return ":="
			} else {
				return "="
			}
		}()).IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Get%s", sn).Call(
			args...,
		),
		utils.AssertError(jen.Err(), nil),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	return lines
}

func buildTestReadingShouldBeReadable(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines,
		jen.Line(),
		jen.Commentf("Fetch %s.", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Line(),
		jen.Commentf("Assert %s equality.", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual")),
		jen.Line(),
		jen.Commentf("Clean up %s.", typ.Name.SingularCommonName()),
		utils.AssertNoError(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", typ.Name.Singular()).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			nil,
		),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ, false)...)

	return lines
}

func buildParamsForCheckingATypeThatDoesNotExistButIncludesPredecessorID(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	params := []jen.Code{constants.CtxVar()}

	for _, pt := range parents {
		params = append(params, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
	}

	params = append(params, jen.ID("nonexistentID"))

	return params
}

func buildParamsForCheckingATypeThatDoesNotExistAndIncludesItsOwnerVar(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if !(i == len(parents)-1 && typ.BelongsToStruct != nil) {
				listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
			}
		}

		params = append(params, listParams...)
	}

	params = append(params, jen.ID(utils.BuildFakeVarName(sn)))

	return params
}

func buildRequisiteCreationCodeForUpdate404Tests(proj *models.Project, typ models.DataType, nonexistentArgIndex int) (lines []jen.Code) {
	pkguvn := proj.Name.UnexportedVarName()

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
			jen.List(jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()).Assign().IDf("%sClient", pkguvn).Dotf("Create%s", ots).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()),
			jen.Line(),
		)
	}

	var isTheRealDeal bool
	updateArgs := buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)
	revisedUpdateArgs := []jen.Code{updateArgs[0]}
	for i := 1; i < len(updateArgs); i++ {
		if i-1 == nonexistentArgIndex && i == len(updateArgs)-1 {
			isTheRealDeal = true
			revisedUpdateArgs = append(revisedUpdateArgs, updateArgs[i])
		} else if i-1 == nonexistentArgIndex && i != len(updateArgs)-1 {
			revisedUpdateArgs = append(revisedUpdateArgs, jen.ID("nonexistentID"))
		} else {
			revisedUpdateArgs = append(revisedUpdateArgs, updateArgs[i])
		}
	}

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	creationArgs := append(buildCreationArguments(proj, createdVarPrefix, typ), jen.IDf("example%sInput", sn))

	lines = append(lines,
		jen.Commentf("Create %s.", scn),
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
		jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()).Assign().IDf("%sClient", pkguvn).Dotf("Create%s", sn).Call(
			creationArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()),
		jen.Line(),
		jen.Commentf("Change %s.", scn),
		jen.IDf("%s%s", createdVarPrefix, sn).Dot("Update").Call(
			jen.ID(utils.BuildFakeVarName(sn)).Dot("ToUpdateInput").Call(),
		),
		func() jen.Code {
			if typ.BelongsToStruct != nil && isTheRealDeal {
				return jen.IDf("%s%s", createdVarPrefix, sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID("nonexistentID")
			}
			return jen.Null()
		}(),
		jen.Err().Equals().IDf("%sClient", pkguvn).Dotf("Update%s", sn).Call(
			revisedUpdateArgs...,
		),
		utils.AssertError(jen.Err(), nil),
		jen.Line(),
	)

	return lines
}

func buildRequisiteCleanupCodeForUpdate404s(proj *models.Project, typ models.DataType) (lines []jen.Code) {
	var typesToCleanup []models.DataType

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		typesToCleanup = append(typesToCleanup, ot)
	}
	typesToCleanup = append(typesToCleanup, typ)

	// reverse it
	for i, j := 0, len(typesToCleanup)-1; i < j; i, j = i+1, j-1 {
		typesToCleanup[i], typesToCleanup[j] = typesToCleanup[j], typesToCleanup[i]
	}

	for _, t := range typesToCleanup {
		lines = append(lines,
			jen.Line(),
			jen.Commentf("Clean up %s.", t.Name.SingularCommonName()),
			utils.AssertNoError(
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", t.Name.Singular()).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, t)...,
				),
				nil,
			),
		)
	}

	return lines
}

func buildSubtestsForUpdate404Tests(proj *models.Project, typ models.DataType) []jen.Code {
	var subtests []jen.Code

	for i, ot := range proj.FindOwnerTypeChain(typ) {
		lines := append(
			[]jen.Code{
				utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
			},
			buildRequisiteCreationCodeForUpdate404Tests(proj, typ, i)...,
		)

		updateArgs := buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)
		revisedUpdateArgs := []jen.Code{updateArgs[0]}
		for j := 1; j < len(updateArgs); j++ {
			if j-1 == i {
				revisedUpdateArgs = append(revisedUpdateArgs, jen.ID("nonexistentID"))
			} else {
				revisedUpdateArgs = append(revisedUpdateArgs, updateArgs[j])
			}
		}

		cleanupCode := buildRequisiteCleanupCodeForUpdate404s(proj, typ)
		if len(cleanupCode) > 0 {
			lines = append(lines, jen.Line())
		}
		lines = append(lines, cleanupCode...)

		subtests = append(subtests,
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				fmt.Sprintf("it should return an error when trying to update something that belongs to %s that does not exist", ot.Name.SingularCommonNameWithPrefix()),
				lines...,
			),
		)
	}

	return subtests
}

func buildTestUpdating(proj *models.Project, typ models.DataType) jen.Code {
	subtests := []jen.Code{
		utils.BuildSubTestWithoutContext(
			"it should return an error when trying to update something that does not exist",
			buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj, typ)...,
		),
		jen.Line(),
		utils.BuildSubTestWithoutContext("it should be updatable", buildTestUpdatingShouldBeUpdatable(proj, typ)...),
	}

	subtests = append(subtests, buildSubtestsForUpdate404Tests(proj, typ)...)

	return jen.ID("test").Dot("Run").Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
		subtests...,
	))
}

func buildTestUpdatingShouldBeUpdatable(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines, jen.Line(),
		jen.Commentf("Change %s.", scn),
		jen.List(jen.IDf("created%s", sn).Dot("Update").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ToUpdateInput").Call())),
		jen.Err().Equals().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Update%s", sn).Call(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), nil),
		jen.Line(),
		jen.Commentf("Fetch %s.", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Line(),
		jen.Commentf("Assert %s equality.", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual")),
		utils.AssertNotNil(jen.ID("actual").Dot("LastUpdatedOn"), nil),
		jen.Line(),
		jen.Commentf("Clean up %s.", typ.Name.SingularCommonName()),
		utils.AssertNoError(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", typ.Name.Singular()).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			nil,
		),
	)

	cleanupCode := buildRequisiteCleanupCode(proj, typ, false)
	if len(cleanupCode) > 0 {
		lines = append(lines, jen.Line())
	}
	lines = append(lines, cleanupCode...)

	return lines
}

func buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
	}
	lines = append(lines, buildRequisiteCreationCodeWithoutType(proj, typ)...)

	args := buildParamsForCheckingATypeThatDoesNotExistAndIncludesItsOwnerVar(proj, typ)

	lines = append(lines,
		utils.BuildFakeVar(proj, sn),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("created%s", typ.BelongsToStruct.Singular()).Dot("ID")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName(sn)).Dot("ID").Equals().ID("nonexistentID"),
		jen.Line(),
		utils.AssertError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Update%s", sn).Call(args...), nil),
	)

	cleanupCode := buildRequisiteCleanupCode(proj, typ, false)
	if len(cleanupCode) > 0 {
		lines = append(lines, jen.Line())
	}
	lines = append(lines, cleanupCode...)

	return lines
}

func buildRequisiteCreationCodeFor404DeletionTests(proj *models.Project, typ models.DataType, nonexistentArgIndex int) (lines []jen.Code) {
	pkguvn := proj.Name.UnexportedVarName()

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
			jen.List(jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()).Assign().IDf("%sClient", pkguvn).Dotf("Create%s", ots).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, ots), jen.Err()),
			jen.Line(),
		)
	}

	return lines
}

func buildRequisiteCleanupCodeFor404DeletionTests(proj *models.Project, typ models.DataType, indexToStop int) (lines []jen.Code) {
	var typesToCleanup []models.DataType

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		typesToCleanup = append(typesToCleanup, ot)
	}
	typesToCleanup = append(typesToCleanup, typ)

	// reverse it
	for i, j := 0, len(typesToCleanup)-1; i < j; i, j = i+1, j-1 {
		typesToCleanup[i], typesToCleanup[j] = typesToCleanup[j], typesToCleanup[i]
	}

	for _, t := range typesToCleanup {
		lines = append(lines,
			jen.Line(),
			jen.Commentf("Clean up %s.", t.Name.SingularCommonName()),
			utils.AssertNoError(
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", t.Name.Singular()).Call(
					buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, t)...,
				),
				nil,
			),
		)
	}

	return lines
}

func buildParamsForDeletionWithNonexistentID(proj *models.Project, typ models.DataType, indexToNotExist int) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	var listParams []jen.Code
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for i, pt := range parents {
			if i == indexToNotExist {
				listParams = append(listParams, jen.ID("nonexistentID"))
			} else {
				listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
			}
		}
		listParams = append(listParams, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))

		params = append(params, listParams...)
	} else {
		params = append(params, jen.IDf("created%s", typ.Name.Singular()).Dot("ID"))
	}

	return params
}

func buildSubtestsForDeletion404Tests(proj *models.Project, typ models.DataType) []jen.Code {
	var subtests []jen.Code
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	for i, ot := range proj.FindOwnerTypeChain(typ) {
		lines := append(
			[]jen.Code{
				utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
			},
			buildRequisiteCreationCodeFor404DeletionTests(proj, typ, i)...,
		)

		creationArgs := append(buildCreationArguments(proj, "created", typ), jen.IDf("example%sInput", sn))
		archiveArgs := buildParamsForDeletionWithNonexistentID(proj, typ, i)

		lines = append(lines,
			jen.Line(),
			jen.Commentf("Create %s.", scn),
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
			jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()).Assign().IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Create%s", sn).Call(
				creationArgs...,
			),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()),
			jen.Line(),
			utils.AssertError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", sn).Call(
				archiveArgs...,
			), nil),
		)

		lines = append(lines, buildRequisiteCleanupCodeFor404DeletionTests(proj, typ, i)...)

		subtests = append(subtests,
			jen.Line(),
			utils.BuildSubTestWithoutContext(fmt.Sprintf("returns error when trying to archive post belonging to nonexistent %s", ot.Name.SingularCommonName()),
				lines...,
			),
		)
	}

	return subtests
}

func buildTestDeleting(proj *models.Project, typ models.DataType) jen.Code {
	subtests := []jen.Code{
		utils.BuildSubTestWithoutContext("it should return an error when trying to delete something that does not exist", buildTestDeletingShouldFailForNonexistent(proj, typ)...),
		jen.Line(),
		utils.BuildSubTestWithoutContext("should be able to be deleted", buildTestDeletingShouldBeAbleToBeDeleted(proj, typ)...),
	}

	notFoundSubtests := buildSubtestsForDeletion404Tests(proj, typ)
	if len(notFoundSubtests) > 0 {
		subtests = append(subtests, jen.Line())
	}
	subtests = append(subtests, notFoundSubtests...)

	return jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
		subtests...,
	))
}

func buildTestDeletingShouldBeAbleToBeDeleted(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
		jen.Line(),
	}

	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines,
		jen.Line(),
		jen.Commentf("Clean up %s.", typ.Name.SingularCommonName()),
		utils.AssertNoError(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", typ.Name.Singular()).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			nil,
		),
	)

	cleanupCode := buildRequisiteCleanupCode(proj, typ, false)
	if len(cleanupCode) > 0 {
		lines = append(lines, jen.Line())
	}
	lines = append(lines, cleanupCode...)

	return lines
}

func buildTestDeletingShouldFailForNonexistent(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{
		utils.StartSpanWithInlineCtx(proj, true, jen.ID("t").Dot("Name").Call()),
		jen.Line(),
	}

	archiveArgs := buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)
	archiveArgs[len(archiveArgs)-1] = jen.ID("nonexistentID")

	lines = append(lines, buildRequisiteCreationCodeWithoutType(proj, typ)...)

	lines = append(lines,
		jen.Line(),
		utils.AssertError(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", typ.Name.Singular()).Call(
				archiveArgs...,
			),
			nil,
		),
	)

	cleanupCode := buildRequisiteCleanupCode(proj, typ, false)
	if len(cleanupCode) > 0 {
		lines = append(lines, jen.Line())
	}
	lines = append(lines, cleanupCode...)

	return lines
}
