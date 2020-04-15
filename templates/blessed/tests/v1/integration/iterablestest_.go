package integration

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{constants.CtxVar()}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
		}
		listParams = listParams[:len(listParams)-1]

		if len(listParams) > 0 {
			params = append(params, listParams...)
		}
	}
	params = append(params, jen.IDf("created%s", typ.Name.Singular()))

	return params
}

func buildCreationArguments(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	creationArgs := []jen.Code{}

	if typ.BelongsToStruct != nil {
		parentTyp := proj.FindType(typ.BelongsToStruct.Singular())
		if parentTyp != nil {
			nca := buildCreationArguments(proj, varPrefix, *parentTyp)
			creationArgs = append(creationArgs, nca...)
		}
	}

	creationArgs = append(creationArgs, jen.IDf("%s%s", varPrefix, typ.Name.Singular()).Dot("ID"))

	return creationArgs
}

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(proj, ret)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()

	ret.Add(
		jen.Func().IDf("check%sEquality", sn).Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).PointerTo().Qual(proj.ModelsV1Package(), sn)).Block(
			buildEqualityCheckLines(typ)...,
		),
		jen.Line(),
	)

	ret.Add(buildBuildDummySomething(proj, typ)...)

	ret.Add(
		jen.Func().IDf("Test%s", pn).Params(jen.ID("test").PointerTo().Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be createable",
					buildTestCreating(proj, typ)...,
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be read in a list",
					buildTestListing(proj, typ)...,
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("ExistenceChecking"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to check something that does not exist",
					buildTestExistenceCheckingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, typ)...,
				),
				jen.Line(),
				utils.BuildSubTest(
					fmt.Sprintf("it should return 200 when the relevant %s exists", scn),
					buildTestExistenceCheckingShouldBeReadable(proj, typ)...,
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to read something that does not exist",
					buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, typ)...,
				),
				jen.Line(),
				utils.BuildSubTest(
					"it should be readable",
					buildTestReadingShouldBeReadable(proj, typ)...,
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"it should return an error when trying to update something that does not exist",
					buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj, typ)...,
				),
				jen.Line(),
				utils.BuildSubTest(
					"it should be updatable",
					buildTestUpdatingShouldBeUpdateable(proj, typ)...,
				),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
				utils.BuildSubTest(
					"should be able to be deleted",
					buildTestDeletingShouldBeAbleToBeDeleted(proj, typ)...,
				),
			)),
		),
	)

	return ret
}

func buildRequisiteCreationCode(proj *models.Project, typ models.DataType) []jen.Code {
	var lines []jen.Code
	sn := typ.Name.Singular()

	const (
		createdVarPrefix = "created"
	)

	creationArgs := []jen.Code{
		constants.CtxVar(),
	}
	ca := buildCreationArguments(proj, createdVarPrefix, typ)
	creationArgs = append(creationArgs, ca[:len(ca)-1]...)
	creationArgs = append(creationArgs, jen.IDf("example%sInput", sn))

	if typ.BelongsToStruct != nil {
		if parentTyp := proj.FindType(typ.BelongsToStruct.Singular()); parentTyp != nil {
			newLines := buildRequisiteCreationCode(proj, *parentTyp)
			lines = append(lines, newLines...)
		}
	}

	lines = append(lines,
		jen.Commentf("Create %s", typ.Name.SingularCommonName()),
		utils.BuildFakeVar(proj, sn),
		utils.BuildFakeVarWithCustomName(
			proj,
			fmt.Sprintf("example%sInput", sn),
			fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
			jen.ID(utils.BuildFakeVarName(sn)),
		),
		jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()).Assign().ID("todoClient").Dotf("Create%s", sn).Call(
			creationArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()),
		jen.Line(),
	)

	return lines
}

func buildRequisiteCreationCodeForUpdateFunction(proj *models.Project, typ models.DataType) []jen.Code {
	var lines []jen.Code
	sn := typ.Name.Singular()

	const (
		createdVarPrefix = "created"
	)

	creationArgs := []jen.Code{
		constants.CtxVar(),
	}
	ca := buildCreationArguments(proj, createdVarPrefix, typ)
	creationArgs = append(creationArgs, ca[:len(ca)-1]...)
	creationArgs = append(creationArgs,
		jen.IDf("example%sInput", sn),
	)

	if typ.BelongsToStruct != nil {
		if parentTyp := proj.FindType(typ.BelongsToStruct.Singular()); parentTyp != nil {
			newLines := buildRequisiteCreationCode(proj, *parentTyp)
			lines = append(lines, newLines...)
		}
	}

	lines = append(lines,
		jen.Commentf("Create %s", typ.Name.SingularCommonName()),
		utils.BuildFakeVar(proj, sn),
		utils.BuildFakeVarWithCustomName(
			proj,
			fmt.Sprintf("example%sInput", sn),
			fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
			jen.ID(utils.BuildFakeVarName(sn)),
		),
		jen.List(jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()).Assign().ID("todoClient").Dotf("Create%s", sn).Call(
			creationArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, sn), jen.Err()),
		jen.Line(),
	)

	return lines
}

func buildRequisiteCleanupCode(proj *models.Project, typ models.DataType) []jen.Code {
	var lines []jen.Code
	sn := typ.Name.Singular()

	lines = append(lines,
		jen.Line(),
		jen.Commentf("Clean up %s", typ.Name.SingularCommonName()),
		utils.AssertNoError(
			jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			nil,
		),
	)

	if typ.BelongsToStruct != nil {
		if parentTyp := proj.FindType(typ.BelongsToStruct.Singular()); parentTyp != nil {
			newLines := buildRequisiteCleanupCode(proj, *parentTyp)
			lines = append(lines, newLines...)
		}
	}

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
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

func buildBuildDummySomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	blockLines := []jen.Code{
		jen.ID("t").Dot("Helper").Call(),
		jen.Line(),
	}

	stopIndex := 6 // the number of `jen.Line`s we need to skip some irrelevant bits of creation code
	cc := buildRequisiteCreationCode(proj, typ)
	if len(cc) > stopIndex {
		blockLines = append(blockLines, cc[:len(cc)-stopIndex]...)
	}

	creationArgs := []jen.Code{constants.CtxVar()}
	ca := buildCreationArguments(proj, "created", typ)
	creationArgs = append(creationArgs, ca[:len(ca)-1]...)
	creationArgs = append(creationArgs, jen.ID("exampleInput"))

	blockLines = append(blockLines,
		constants.CreateCtx(),
		utils.BuildFakeVar(proj, sn),
		utils.BuildFakeVarWithCustomName(
			proj,
			"exampleInput",
			fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn),
			jen.ID(utils.BuildFakeVarName(sn)),
		),
		jen.List(jen.ID("y"), jen.Err()).Assign().ID("todoClient").Dotf("Create%s", sn).Call(
			creationArgs...,
		),
		utils.RequireNoError(jen.Err(), nil),
		jen.Line(),
		jen.Return().ID("y"),
	)

	lines := []jen.Code{
		jen.Func().IDf("buildDummy%s", sn).Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Block(
			blockLines...,
		),
		jen.Line(),
	}

	return lines
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

func buildTestCreating(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
		jen.Line(),
	}

	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines,
		jen.Commentf("Assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.IDf("created%s", typ.Name.Singular())),
		jen.Line(),
		jen.Comment("Clean up"),
		jen.Err().Equals().ID("todoClient").Dotf("Archive%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), nil),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual")),
		utils.AssertNotZero(jen.ID("actual").Dot("ArchivedOn"), nil),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ)[3:]...)

	return lines
}

func buildTestListing(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	listArgs := []jen.Code{constants.CtxVar()}
	ca := buildCreationArguments(proj, "created", typ)
	listArgs = append(listArgs, ca[:len(ca)-1]...)
	listArgs = append(listArgs, jen.Nil())

	stopIndex := 6 // the number of `jen.Line`s we need to skip some irrelevant bits of creation code
	cc := buildRequisiteCreationCode(proj, typ)
	if len(cc) > stopIndex {
		lines = append(lines, cc[:len(cc)-stopIndex]...)
	}
	cc = append(cc, jen.ID("expected").Equals().ID("append").Call(jen.ID("expected"), jen.IDf("created%s", typ.Name.Singular())))

	lines = append(lines,
		jen.Commentf("Create %s", pcn),
		jen.Var().ID("expected").Index().PointerTo().Qual(proj.ModelsV1Package(), sn),
		jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Block(
			jen.ID("expected").Equals().Append(jen.ID("expected"), jen.IDf("buildDummy%s", sn).Call(jen.ID("t"))),
		),
		jen.Line(),
		jen.Commentf("Assert %s list equality", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dotf("Get%s", pn).Call(
			listArgs...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Qual(utils.AssertPkg, "True").Callln(
			jen.ID("t"),
			jen.Len(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(pn)),
			jen.Lit("expected %d to be <= %d"), jen.Len(jen.ID("expected")),
			jen.Len(jen.ID("actual").Dot(pn)),
		),
		jen.Line(),
		jen.Comment("Clean up"),
		jen.For(jen.List(jen.Underscore(), jen.IDf("created%s", sn)).Assign().Range().ID("actual").Dot(pn)).Block(
			jen.Err().Equals().ID("todoClient").Dotf("Archive%s", sn).Call(
				buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
			),
			utils.AssertNoError(jen.Err(), nil),
		),
	)

	ccsi := 3 // cleanupCodeStopIndex: the number of `jen.Line`s we need to skip some irrelevant bits of cleanup code
	dc := buildRequisiteCleanupCode(proj, typ)
	if len(dc) > ccsi {
		dc = dc[ccsi:]
	} else if len(dc) == ccsi {
		dc = []jen.Code{}
	}
	lines = append(lines, dc...)

	return lines
}

func buildTestExistenceCheckingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	args := buildParamsForCheckingATypeThatDoesNotExistButIncludesPredecessorID(proj, typ)

	stopIndex := 6 // the number of `jen.Line`s we need to skip some irrelevant bits of creation code
	cc := buildRequisiteCreationCode(proj, typ)
	if len(cc) > stopIndex {
		lines = append(lines, cc[:len(cc)-stopIndex]...)
	}

	lines = append(lines,
		jen.Commentf("Attempt to fetch nonexistent %s", scn),
		jen.List(jen.ID("actual"), jen.Err()).Op(func() string {
			if typ.BelongsToStruct == nil {
				return ":="
			} else {
				return "="
			}
		}()).ID("todoClient").Dotf("%sExists", sn).Call(
			args...,
		),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertFalse(jen.ID("actual"), nil),
	)

	ccsi := 3 // cleanupCodeStopIndex: the number of `jen.Line`s we need to skip some irrelevant bits of cleanup code
	dc := buildRequisiteCleanupCode(proj, typ)
	if len(dc) > ccsi {
		dc = dc[ccsi:]
	} else if len(dc) == ccsi {
		dc = []jen.Code{}
	}

	lines = append(lines, dc...)

	return lines
}

func buildTestExistenceCheckingShouldBeReadable(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines,
		jen.Line(),
		jen.Commentf("Fetch %s", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dotf("%sExists", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertTrue(jen.ID("actual"), nil),
		jen.Line(),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ)...)

	return lines
}

func buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	args := buildParamsForCheckingATypeThatDoesNotExistButIncludesPredecessorID(proj, typ)

	stopIndex := 6 // the number of `jen.Line`s we need to skip some irrelevant bits of creation code
	cc := buildRequisiteCreationCode(proj, typ)
	if len(cc) > stopIndex {
		lines = append(lines, cc[:len(cc)-stopIndex]...)
	}

	lines = append(lines,
		jen.Commentf("Attempt to fetch nonexistent %s", scn),
		jen.List(jen.Underscore(), jen.Err()).Op(func() string {
			if typ.BelongsToStruct == nil {
				return ":="
			} else {
				return "="
			}
		}()).ID("todoClient").Dotf("Get%s", sn).Call(
			args...,
		),
		utils.AssertError(jen.Err(), nil),
	)

	ccsi := 3 // cleanupCodeStopIndex: the number of `jen.Line`s we need to skip some irrelevant bits of cleanup code
	dc := buildRequisiteCleanupCode(proj, typ)
	if len(dc) > ccsi {
		dc = dc[ccsi:]
	} else if len(dc) == ccsi {
		dc = []jen.Code{}
	}

	lines = append(lines, dc...)

	return lines
}

func buildTestReadingShouldBeReadable(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)

	lines = append(lines,
		jen.Line(),
		jen.Commentf("Fetch %s", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Line(),
		jen.Commentf("Assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual")),
		jen.Line(),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ)...)

	return lines
}

func buildParamsForCheckingATypeThatDoesNotExist(proj *models.Project, typ models.DataType) []jen.Code {
	parents := proj.FindOwnerTypeChain(typ)
	params := []jen.Code{constants.CtxVar()}

	for i, pt := range parents {
		if i == len(parents)-1 {
			params = append(params, jen.ID("nonexistentID"))
		} else {
			params = append(params, jen.IDf("created%s", pt.Name.Singular()).Dot("ID"))
		}
	}

	if len(params) == 0 {
		params = append(params, jen.ID("nonexistentID"))
	}

	return params
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

func buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
		jen.Line(),
	}

	args := buildParamsForCheckingATypeThatDoesNotExistAndIncludesItsOwnerVar(proj, typ)

	lines = append(lines,
		utils.BuildFakeVar(proj, sn),
		jen.ID(utils.BuildFakeVarName(sn)).Dot("ID").Equals().ID("nonexistentID"),
		jen.Line(),
		utils.AssertError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dotf("Update%s", sn).Call(args...), nil),
	)

	ccsi := 3 // cleanupCodeStopIndex: the number of `jen.Line`s we need to skip some irrelevant bits of cleanup code
	dc := buildRequisiteCleanupCode(proj, typ)
	if len(dc) > ccsi {
		dc = dc[ccsi:]
	} else if len(dc) == ccsi {
		dc = []jen.Code{}
	}

	lines = append(lines, dc...)

	return lines
}

func buildTestUpdatingShouldBeUpdateable(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
		jen.Line(),
	}

	creationCode := buildRequisiteCreationCodeForUpdateFunction(proj, typ)
	//stopIndex := 5 // the number of `jen.Line`s we need to skip some irrelevant bits of creation code

	//if len(creationCode) > stopIndex {
	//	precursorCode := creationCode[:len(creationCode)-stopIndex]
	//	postcursorCode := creationCode[len(creationCode)-stopIndex:]
	//
	//	lines = append(lines, precursorCode...)
	//	lines = append(lines, postcursorCode...)
	//} else {
	//}
	lines = append(lines, creationCode...)

	lines = append(lines, jen.Line(),
		jen.Commentf("Change %s", scn),
		jen.List(jen.IDf("created%s", sn).Dot("Update").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ToUpdateInput").Call())),
		jen.Err().Equals().ID("todoClient").Dotf("Update%s", sn).Call(
			buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)...,
		),
		utils.AssertNoError(jen.Err(), nil),
		jen.Line(),
		jen.Commentf("Fetch %s", scn),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("todoClient").Dotf("Get%s", sn).Call(
			buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)...,
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
		jen.Line(),
		jen.Commentf("Assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual")),
		utils.AssertNotNil(jen.ID("actual").Dot("UpdatedOn"), nil),
		jen.Line(),
	)

	lines = append(lines, buildRequisiteCleanupCode(proj, typ)...)

	return lines
}

func buildTestDeletingShouldBeAbleToBeDeleted(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{
		utils.StartSpanWithVar(proj, true, jen.ID("t").Dot("Name").Call()),
		jen.Line(),
	}

	lines = append(lines, buildRequisiteCreationCode(proj, typ)...)
	lines = append(lines, buildRequisiteCleanupCode(proj, typ)...)

	return lines
}
