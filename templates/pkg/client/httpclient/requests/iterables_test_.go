package requests

import (
	"fmt"
	"path"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestBuilder_BuildSomethingExistsRequest(proj, typ)...)
	code.Add(buildTestBuilder_BuildGetSomethingRequest(proj, typ)...)
	code.Add(buildTestBuilder_BuildGetListOfSomethingsRequest(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildTestBuilder_BuildSearchSomethingRequest(proj, typ)...)
	}

	code.Add(buildTestBuilder_BuildCreateSomethingRequest(proj, typ)...)
	code.Add(buildTestBuilder_BuildUpdateSomethingRequest(proj, typ)...)
	code.Add(buildTestBuilder_BuildArchiveSomethingRequest(proj, typ)...)
	code.Add(buildTestBuilder_BuildGetAuditLogForSomethingRequest(proj, typ)...)

	return code
}

func buildPrerequisiteIDs(proj *models.Project, typ models.DataType, includeSelf bool) []jen.Code {
	lines := []jen.Code{}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("example%sID", dep.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
	}

	if includeSelf {
		lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	}

	return lines
}

func buildPrerequisiteIDsWithoutIndex(proj *models.Project, typ models.DataType, index int, includeSelf bool) []jen.Code {
	lines := []jen.Code{}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i != index {
			lines = append(lines, jen.IDf("example%sID", dep.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
		}
	}

	if includeSelf {
		lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	}

	return lines
}

func buildSomethingGeneralArgsWithoutIndex(proj *models.Project, typ models.DataType, index int) []jen.Code {
	parts := []jen.Code{jen.ID("helper").Dot("ctx")}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.Zero())
		} else {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	if index != -1 {
		parts = append(parts, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))
	} else {
		parts = append(parts, jen.Zero())
	}

	return parts
}

func buildSomethingSpecificFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%d")
	}

	parts = append(parts, typ.Name.PluralRouteName(), "%d")

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildSomethingGeneralArgs(proj *models.Project, typ models.DataType, includeCtx bool) (parts []jen.Code) {
	if includeCtx {
		parts = []jen.Code{jen.ID("helper").Dot("ctx")}
	} else {
		parts = []jen.Code{}
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
	}

	parts = append(parts, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))

	return parts
}

func buildTestBuilder_BuildSomethingExistsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	firstSubtest = append(firstSubtest, buildPrerequisiteIDs(proj, typ, true)...)

	firstSubtest = append(firstSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("Build%sExistsRequest", sn).Call(
			buildSomethingGeneralArgs(proj, typ, true)...,
		),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			append(
				[]jen.Code{
					jen.ID("true"),
					jen.Qual("net/http", "MethodHead"),
					jen.Lit(""),
					jen.ID("expectedPathFormat"),
				},
				buildSomethingGeneralArgs(proj, typ, false)...,
			)...,
		),
		jen.Newline(),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.ID("assertRequestQuality").Call(jen.ID("t"), jen.ID("actual"), jen.ID("spec")),
	)

	expectedPathFormat := buildSomethingSpecificFormatString(proj, typ)

	lines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Const().ID("expectedPathFormat").Equals().Lit(expectedPathFormat),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
	}

	for i, parent := range proj.FindOwnerTypeChain(typ) {
		pscn := parent.Name.SingularCommonName()

		buildArgs := buildSomethingGeneralArgsWithoutIndex(proj, typ, i)

		subtestLines := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("helper").Assign().ID("buildTestHelper").Call(),
			jen.Newline(),
		}

		subtestLines = append(subtestLines, buildPrerequisiteIDsWithoutIndex(proj, typ, i, true)...)

		subtestLines = append(subtestLines,
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("Build%sExistsRequest", sn).Call(
				buildArgs...,
			),
			jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
		)

		lines = append(lines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", pscn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					subtestLines...,
				),
			),
		)
	}

	subtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	buildArgs := buildSomethingGeneralArgsWithoutIndex(proj, typ, -1)

	subtestLines = append(subtestLines, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	subtestLines = append(subtestLines,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("Build%sExistsRequest", sn).Call(
			buildArgs...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				subtestLines...,
			),
		),
	)

	invalidBuilderSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	invalidBuilderSubtest = append(invalidBuilderSubtest, buildPrerequisiteIDs(proj, typ, true)...)

	invalidBuilderSubtest = append(invalidBuilderSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("Build%sExistsRequest", sn).Call(
			buildSomethingGeneralArgs(proj, typ, true)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				invalidBuilderSubtest...,
			),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestBuilder_Build%sExistsRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(lines...),
		jen.Newline(),
	}
}

func buildTestBuilder_BuildGetSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	happyPathSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	happyPathSubtest = append(happyPathSubtest, buildPrerequisiteIDs(proj, typ, true)...)

	happyPathSubtest = append(happyPathSubtest,
		jen.Newline(),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			append(
				[]jen.Code{
					jen.ID("true"),
					jen.Qual("net/http", "MethodGet"),
					jen.Lit(""),
					jen.ID("expectedPathFormat"),
				},
				buildSomethingGeneralArgs(proj, typ, false)...,
			)...,
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGet%sRequest", sn).Call(
			buildSomethingGeneralArgs(proj, typ, true)...,
		),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.ID("assertRequestQuality").Call(jen.ID("t"), jen.ID("actual"), jen.ID("spec")),
	)

	lines := []jen.Code{
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				happyPathSubtest...,
			),
		),
	}

	for i, parent := range proj.FindOwnerTypeChain(typ) {
		pscn := parent.Name.SingularCommonName()

		subtestLines := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("helper").Assign().ID("buildTestHelper").Call(),
			jen.Newline(),
		}

		subtestLines = append(subtestLines, buildPrerequisiteIDsWithoutIndex(proj, typ, i, true)...)

		subtestLines = append(subtestLines,
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGet%sRequest", sn).Call(
				buildSomethingGeneralArgsWithoutIndex(proj, typ, i)...,
			),
			jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
		)

		lines = append(lines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", pscn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					subtestLines...,
				),
			),
		)
	}

	subtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	subtestLines = append(subtestLines, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	subtestLines = append(subtestLines,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGet%sRequest", sn).Call(
			buildSomethingGeneralArgsWithoutIndex(proj, typ, -1)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", typ.Name.SingularCommonName()),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				subtestLines...,
			),
		),
	)

	invalidBuilderSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	invalidBuilderSubtest = append(invalidBuilderSubtest, buildPrerequisiteIDs(proj, typ, true)...)

	invalidBuilderSubtest = append(invalidBuilderSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGet%sRequest", sn).Call(
			buildSomethingGeneralArgs(proj, typ, true)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				invalidBuilderSubtest...,
			),
		),
	)

	expectedPathFormat := buildSomethingSpecificFormatString(proj, typ)

	return []jen.Code{
		jen.Func().IDf("TestBuilder_BuildGet%sRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			append([]jen.Code{
				jen.ID("T").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Const().ID("expectedPathFormat").Equals().Lit(expectedPathFormat),
				jen.Newline(),
			},
				lines...,
			)...,
		),
		jen.Newline(),
	}
}

func buildListOfSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%d")
	}

	parts = append(parts, typ.Name.PluralRouteName())

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildPrerequisiteIDsWithoutIndexForList(proj *models.Project, typ models.DataType, index int) []jen.Code {
	lines := []jen.Code{}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i != index {
			lines = append(lines, jen.IDf("example%sID", dep.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
		}
	}

	return lines
}

func buildGetListOfSomethingFormatArgs(proj *models.Project, typ models.DataType, includeCtx bool) (parts []jen.Code) {
	if includeCtx {
		parts = []jen.Code{jen.ID("helper").Dot("ctx")}
	} else {
		parts = []jen.Code{}
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
	}

	return parts
}

func buildListOfSomethingArgsWithoutIndex(proj *models.Project, typ models.DataType, index int) []jen.Code {
	parts := []jen.Code{jen.ID("helper").Dot("ctx")}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.Zero())
		} else {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	parts = append(parts, jen.ID("filter"))

	return parts
}

func buildTestBuilder_BuildGetListOfSomethingsRequest(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	expectedPathFormat := buildListOfSomethingFormatString(proj, typ)

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.Newline(),
		jen.ID("filter").Assign().Parens(jen.PointerTo().ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			append(
				[]jen.Code{
					jen.ID("true"),
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
					jen.ID("expectedPathFormat"),
				},
				buildGetListOfSomethingFormatArgs(proj, typ, false)...,
			)...,
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGet%sRequest", pn).Call(
			buildListOfSomethingArgsWithoutIndex(proj, typ, -1)...,
		),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.ID("assertRequestQuality").Call(jen.ID("t"), jen.ID("actual"), jen.ID("spec")),
	)

	lines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Const().ID("expectedPathFormat").Equals().Lit(expectedPathFormat),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtestLines...,
			),
		),
	}

	for i, parent := range proj.FindOwnerTypeChain(typ) {
		pscn := parent.Name.SingularCommonName()

		subtestLines := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("helper").Assign().ID("buildTestHelper").Call(),
			jen.Newline(),
		}

		subtestLines = append(subtestLines, buildPrerequisiteIDsWithoutIndexForList(proj, typ, i)...)

		subtestLines = append(subtestLines,
			jen.Newline(),
			jen.ID("filter").Assign().Parens(jen.PointerTo().ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGet%sRequest", pn).Call(
				buildListOfSomethingArgsWithoutIndex(proj, typ, i)...,
			),
			jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
		)

		lines = append(lines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", pscn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					subtestLines...,
				),
			),
		)
	}

	invalidBuilderSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	invalidBuilderSubtest = append(invalidBuilderSubtest, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	invalidBuilderSubtest = append(invalidBuilderSubtest,
		jen.Newline(),
		jen.ID("filter").Assign().Parens(jen.PointerTo().ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGet%sRequest", pn).Call(
			buildListOfSomethingArgsWithoutIndex(proj, typ, -1)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				invalidBuilderSubtest...,
			),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestBuilder_BuildGet%sRequest", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildSearchSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%d")
	}

	parts = append(parts, typ.Name.PluralRouteName(), "search")

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildSearchSomethingFormatArgs(proj *models.Project, typ models.DataType, index int) []jen.Code {
	parts := []jen.Code{}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.Zero())
		} else {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	return parts
}

func buildSearchSomethingCallArgs(proj *models.Project, typ models.DataType, index int) []jen.Code {
	parts := []jen.Code{}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.Zero())
		} else {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	parts = append(parts, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))

	return parts
}

func buildSearchSomethingRequestBuilderArgs(proj *models.Project, typ models.DataType, index int) []jen.Code {
	parts := []jen.Code{
		jen.ID("helper").Dot("ctx"),
	}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.Zero())
		} else {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	parts = append(parts,
		jen.ID("exampleQuery"),
		jen.ID("limit"),
	)

	return parts
}

func buildTestBuilder_BuildSearchSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	firstSubtest = append(firstSubtest, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	firstSubtest = append(firstSubtest,
		jen.Newline(),
		jen.ID("limit").Assign().ID("types").Dot("DefaultQueryFilter").Call().Dot("Limit"),
		jen.ID("exampleQuery").Assign().Lit("whatever"),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			append([]jen.Code{
				jen.ID("true"),
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("limit=20&q=whatever"),
				jen.ID("expectedPath"),
			},
				buildSearchSomethingFormatArgs(proj, typ, -1)...,
			)...,
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildSearch%sRequest", pn).Call(
			buildSearchSomethingRequestBuilderArgs(proj, typ, -1)...,
		),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.ID("assertRequestQuality").Call(
			jen.ID("t"),
			jen.ID("actual"),
			jen.ID("spec"),
		),
	)

	bodyLines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Const().ID("expectedPath").Equals().Lit(buildSearchSomethingFormatString(proj, typ)),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		subtest := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("helper").Assign().ID("buildTestHelper").Call(),
			jen.Newline(),
		}

		subtest = append(subtest, buildPrerequisiteIDsWithoutIndex(proj, typ, i, false)...)

		subtest = append(subtest,
			jen.Newline(),
			jen.ID("limit").Assign().ID("types").Dot("DefaultQueryFilter").Call().Dot("Limit"),
			jen.ID("exampleQuery").Assign().Lit("whatever"),
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildSearch%sRequest", pn).Call(
				buildSearchSomethingRequestBuilderArgs(proj, typ, i)...,
			),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		)

		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					subtest...,
				),
			),
		)
	}

	secondSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	secondSubtest = append(secondSubtest, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	secondSubtest = append(secondSubtest,
		jen.ID("limit").Assign().ID("types").Dot("DefaultQueryFilter").Call().Dot("Limit"),
		jen.ID("exampleQuery").Assign().Lit("whatever"),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildSearch%sRequest", pn).Call(
			buildSearchSomethingRequestBuilderArgs(proj, typ, -1)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				secondSubtest...,
			),
		),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestBuilder_BuildSearch%sRequest", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildCreateSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%d")
	}

	parts = append(parts, typ.Name.PluralRouteName())

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildPrerequisiteIDsWithoutIndexOrSelf(proj *models.Project, typ models.DataType, index int) []jen.Code {
	owners := proj.FindOwnerTypeChain(typ)
	lines := []jen.Code{}

	for i, dep := range owners {
		if i != index && i != len(owners)-1 {
			lines = append(lines, jen.IDf("example%sID", dep.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
		}
	}

	return lines
}

func buildCreateSomethingFormatArgs(proj *models.Project, typ models.DataType, includeCtx bool) (parts []jen.Code) {
	if includeCtx {
		parts = []jen.Code{jen.ID("helper").Dot("ctx")}
	} else {
		parts = []jen.Code{}
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i, dep := range owners {
		if i != len(owners)-1 {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		} else {
			parts = append(parts, jen.ID("exampleInput").Dotf("BelongsTo%s", dep.Name.Singular()))
		}
	}

	return parts
}

func buildCreateSomethingArgsWithoutIndex(proj *models.Project, typ models.DataType, index int, includeExampleInput bool) []jen.Code {
	parts := []jen.Code{jen.ID("helper").Dot("ctx")}
	owners := proj.FindOwnerTypeChain(typ)

	for i, dep := range owners {
		if i == index {
			parts = append(parts, jen.Zero())
		} else if i != len(owners)-1 {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	if includeExampleInput {
		parts = append(parts, jen.ID("exampleInput"))
	}

	return parts
}

func buildTestBuilder_BuildCreateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	specArgs := append([]jen.Code{
		jen.ID("false"),
		jen.Qual("net/http", "MethodPost"),
		jen.Lit(""),
		jen.ID("expectedPath"),
	},
		buildCreateSomethingFormatArgs(proj, typ, false)...,
	)

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	firstSubtest = append(firstSubtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, -1)...)

	firstSubtest = append(firstSubtest,
		jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sCreationInput", sn).Call(),
		jen.Newline(),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			specArgs...,
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
			buildCreateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
		),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.ID("assertRequestQuality").Call(jen.ID("t"), jen.ID("actual"), jen.ID("spec")),
	)

	expectedPath := buildCreateSomethingFormatString(proj, typ)

	lines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Const().ID("expectedPath").Equals().Lit(expectedPath),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i, owner := range owners {
		if i != len(owners)-1 {
			subtest := []jen.Code{
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("helper").Assign().ID("buildTestHelper").Call(),
				jen.Newline(),
			}

			subtest = append(subtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, i)...)

			subtest = append(subtest,
				jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sCreationInput", sn).Call(),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
					buildCreateSomethingArgsWithoutIndex(proj, typ, i, true)...,
				),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)

			lines = append(lines,
				jen.Newline(),
				jen.ID("T").Dot("Run").Call(
					jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
					jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
						subtest...,
					),
				),
			)
		}
	}

	secondSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	secondSubtest = append(secondSubtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, -1)...)

	secondSubtest = append(secondSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
			append(buildCreateSomethingArgsWithoutIndex(proj, typ, -1, false), jen.Nil())...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	thirdSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	thirdSubtest = append(thirdSubtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, -1)...)

	thirdSubtest = append(thirdSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
			append(
				buildCreateSomethingArgsWithoutIndex(proj, typ, -1, false),
				jen.AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values(),
			)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with nil input"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				secondSubtest...,
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid input"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				thirdSubtest...,
			),
		),
	)

	invalidBuilderSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	invalidBuilderSubtest = append(invalidBuilderSubtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, -1)...)

	invalidBuilderSubtest = append(invalidBuilderSubtest,
		jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sCreationInput", sn).Call(),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
			buildCreateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				invalidBuilderSubtest...,
			),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestBuilder_BuildCreate%sRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildUpdateSomethingSpecArgs(proj *models.Project, typ models.DataType) (parts []jen.Code) {
	parts = []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for i, dep := range owners {
		if i != len(owners)-1 {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		} else {
			parts = append(parts, jen.IDf("example%s", typ.Name.Singular()).Dotf("BelongsTo%s", dep.Name.Singular()))
		}
	}

	parts = append(parts, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))

	return parts
}

func buildUpdateSomethingArgsWithoutIndex(proj *models.Project, typ models.DataType, index int, includeSelf bool) (parts []jen.Code) {
	parts = []jen.Code{jen.ID("helper").Dot("ctx")}

	owners := proj.FindOwnerTypeChain(typ)
	for i, dep := range owners {
		if i == index {
			parts = append(parts, jen.Zero())
		} else if i != len(owners)-1 {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	if includeSelf {
		parts = append(parts, jen.IDf("example%s", typ.Name.Singular()))
	}

	return parts
}

func buildTestBuilder_BuildUpdateSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	specArgs := append([]jen.Code{
		jen.ID("false"),
		jen.Qual("net/http", "MethodPut"),
		jen.Lit(""),
		jen.ID("expectedPathFormat"),
	},
		buildUpdateSomethingSpecArgs(proj, typ)...,
	)

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	firstSubtest = append(firstSubtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, -1)...)

	firstSubtest = append(firstSubtest,
		jen.IDf("example%s", sn).Assign().ID("fakes").Dotf("BuildFake%s", sn).Call(),
		jen.Newline(),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			specArgs...,
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildUpdate%sRequest", sn).Call(
			buildUpdateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
		),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.ID("assertRequestQuality").Call(jen.ID("t"), jen.ID("actual"), jen.ID("spec")),
	)

	expectedPathFormat := buildSomethingSpecificFormatString(proj, typ)
	lines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Const().ID("expectedPathFormat").Equals().Lit(expectedPathFormat),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i, owner := range owners {
		if i != len(owners)-1 {
			subtest := []jen.Code{
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("helper").Assign().ID("buildTestHelper").Call(),
				jen.Newline(),
			}

			subtest = append(subtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, i)...)

			subtest = append(subtest,
				jen.IDf("example%s", sn).Assign().ID("fakes").Dotf("BuildFake%s", sn).Call(),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildUpdate%sRequest", sn).Call(
					buildUpdateSomethingArgsWithoutIndex(proj, typ, i, true)...,
				),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)

			lines = append(lines,
				jen.Newline(),
				jen.ID("T").Dot("Run").Call(
					jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
					jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
						subtest...,
					),
				),
			)
		}
	}

	nilSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	nilSubtest = append(nilSubtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, -1)...)

	nilSubtest = append(nilSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildUpdate%sRequest", sn).Call(
			append(buildUpdateSomethingArgsWithoutIndex(proj, typ, -1, false), jen.Nil())...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)
	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with nil input"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				nilSubtest...,
			),
		),
	)

	invalidBuilderSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	invalidBuilderSubtest = append(invalidBuilderSubtest, buildPrerequisiteIDsWithoutIndexOrSelf(proj, typ, -1)...)

	invalidBuilderSubtest = append(invalidBuilderSubtest,
		jen.IDf("example%s", sn).Assign().ID("fakes").Dotf("BuildFake%s", sn).Call(),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildUpdate%sRequest", sn).Call(
			buildUpdateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				invalidBuilderSubtest...,
			),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestBuilder_BuildUpdate%sRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildTestBuilder_BuildArchiveSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	firstSubtest = append(firstSubtest, buildPrerequisiteIDs(proj, typ, true)...)

	firstSubtest = append(firstSubtest,
		jen.Newline(),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			append(
				[]jen.Code{
					jen.ID("true"),
					jen.Qual("net/http", "MethodDelete"),
					jen.Lit(""),
					jen.ID("expectedPathFormat"),
				},
				buildSomethingGeneralArgs(proj, typ, false)...,
			)...,
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildArchive%sRequest", sn).Call(
			buildSomethingGeneralArgs(proj, typ, true)...,
		),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.ID("assertRequestQuality").Call(jen.ID("t"), jen.ID("actual"), jen.ID("spec")),
	)

	expectedPathFormat := buildSomethingSpecificFormatString(proj, typ)

	lines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Const().ID("expectedPathFormat").Equals().Lit(expectedPathFormat),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		subtest := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("helper").Assign().ID("buildTestHelper").Call(),
			jen.Newline(),
		}

		subtest = append(subtest, buildPrerequisiteIDsWithoutIndex(proj, typ, i, true)...)

		subtest = append(subtest,
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildArchive%sRequest", sn).Call(
				buildSomethingGeneralArgsWithoutIndex(proj, typ, i)...,
			),
			jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
		)

		lines = append(lines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					subtest...,
				),
			),
		)
	}

	subtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	subtest = append(subtest, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	subtest = append(subtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildArchive%sRequest", sn).Call(
			buildSomethingGeneralArgsWithoutIndex(proj, typ, -1)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", typ.Name.SingularCommonName()),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				subtest...,
			),
		),
	)

	invalidRequestBuilderSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	invalidRequestBuilderSubtest = append(invalidRequestBuilderSubtest, buildPrerequisiteIDs(proj, typ, true)...)

	invalidRequestBuilderSubtest = append(invalidRequestBuilderSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildArchive%sRequest", sn).Call(
			buildSomethingGeneralArgs(proj, typ, true)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				invalidRequestBuilderSubtest...,
			),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestBuilder_BuildArchive%sRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			lines...,
		),
	}
}

func buildAuditSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%d")
	}

	parts = append(parts, typ.Name.PluralRouteName(), "%d", "audit")

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildAuditSomethingRequestBuilderArgs(proj *models.Project, typ models.DataType, index int, includeSelf bool) []jen.Code {
	parts := []jen.Code{
		jen.ID("helper").Dot("ctx"),
	}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.Zero())
		} else {
			parts = append(parts, jen.IDf("example%sID", dep.Name.Singular()))
		}
	}

	if includeSelf {
		parts = append(parts, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))
	} else {
		parts = append(parts, jen.Zero())
	}

	return parts
}

func buildTestBuilder_BuildGetAuditLogForSomethingRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	firstSubtest = append(firstSubtest, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, true)...)

	firstSubtest = append(firstSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
			buildAuditSomethingRequestBuilderArgs(proj, typ, -1, true)...,
		),
		jen.ID("require").Dot("NotNil").Call(
			jen.ID("t"),
			jen.ID("actual"),
		),
		jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			append([]jen.Code{
				jen.ID("true"),
				jen.Qual("net/http", "MethodGet"),
				jen.Lit(""),
				jen.ID("expectedPath"),
			},
				buildSearchSomethingCallArgs(proj, typ, -1)...,
			)...,
		),
		jen.ID("assertRequestQuality").Call(jen.ID("t"), jen.ID("actual"), jen.ID("spec")),
	)

	lines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Const().ID("expectedPath").Equals().Lit(buildAuditSomethingFormatString(proj, typ)),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		subtest := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("helper").Assign().ID("buildTestHelper").Call(),
			jen.Newline(),
		}

		subtest = append(subtest, buildPrerequisiteIDsWithoutIndex(proj, typ, i, true)...)

		subtest = append(subtest,
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
				buildAuditSomethingRequestBuilderArgs(proj, typ, i, true)...,
			),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Nil").Call(
				jen.ID("t"),
				jen.ID("actual"),
			),
		)

		lines = append(lines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					subtest...,
				),
			),
		)
	}

	secondSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.Newline(),
	}

	secondSubtest = append(secondSubtest, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, false)...)

	secondSubtest = append(secondSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
			buildAuditSomethingRequestBuilderArgs(proj, typ, -1, false)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", typ.Name.SingularCommonName()),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				secondSubtest...,
			),
		),
	)

	invalidBuilderSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("helper").Assign().ID("buildTestHelper").Call(),
		jen.ID("helper").Dot("builder").Equals().ID("buildTestRequestBuilderWithInvalidURL").Call(),
		jen.Newline(),
	}

	invalidBuilderSubtest = append(invalidBuilderSubtest, buildPrerequisiteIDsWithoutIndex(proj, typ, -1, true)...)

	invalidBuilderSubtest = append(invalidBuilderSubtest,
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("helper").Dot("builder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
			buildAuditSomethingRequestBuilderArgs(proj, typ, -1, true)...,
		),
		jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid request builder"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				invalidBuilderSubtest...,
			),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestBuilder_BuildGetAuditLogFor%sRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			lines...,
		),
		jen.Newline(),
	}
}
