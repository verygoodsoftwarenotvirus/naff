package httpclient

import (
	"fmt"
	"path"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildSuiteStruct(proj, typ)...)

	code.Add(buildTestClientGetSomething(proj, typ)...)
	code.Add(buildTestClientGetListOfSomething(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildTestClientSearchSomething(proj, typ)...)
	}

	code.Add(buildTestClientCreateSomething(proj, typ)...)
	code.Add(buildTestClientUpdateSomething(proj, typ)...)
	code.Add(buildTestClientArchiveSomething(proj, typ)...)

	return code
}

func buildSuiteStruct(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()

	structFields := []jen.Code{
		jen.ID("suite").Dot("Suite"),
		jen.Newline(),
		jen.ID("ctx").Qual("context", "Context"),
	}

	initFields := []jen.Code{
		jen.ID("s").Dot("ctx").Equals().Qual("context", "Background").Call(),
	}

	for _, owner := range proj.FindOwnerTypeChain(typ) {
		structFields = append(structFields, jen.IDf("example%sID", owner.Name.Singular()).Uint64())
		initFields = append(initFields, jen.ID("s").Dotf("example%sID", owner.Name.Singular()).Equals().ID("fakes").Dot("BuildFakeID").Call())
	}

	structFields = append(structFields, jen.IDf("example%s", sn).PointerTo().Qual(proj.TypesPackage(), sn))
	initFields = append(initFields, jen.ID("s").Dotf("example%s", sn).Equals().ID("fakes").Dotf("BuildFake%s", sn).Call())
	if typ.BelongsToStruct != nil {
		initFields = append(initFields, jen.ID("s").Dotf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID("s").Dotf("example%sID", typ.BelongsToStruct.Singular()))
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s", pn).Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.IDf("%sTestSuite", puvn)),
			),
		),
		jen.Newline(),
		jen.Type().IDf("%sBaseSuite", puvn).Struct(
			structFields...,
		),
		jen.Newline(),
		jen.Var().Underscore().ID("suite").Dot("SetupTestSuite").Equals().Parens(jen.PointerTo().IDf("%sBaseSuite", puvn)).Call(jen.Nil()),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sBaseSuite", puvn)).ID("SetupTest").Params().Body(
			initFields...,
		),
		jen.Newline(),
		jen.Type().IDf("%sTestSuite", puvn).Struct(
			jen.ID("suite").Dot("Suite"),
			jen.Newline(),
			jen.IDf("%sBaseSuite", puvn),
		),
		jen.Newline(),
	}

	return lines
}

func buildSomethingSpecificFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, owner := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, owner.Name.PluralRouteName(), "%s")
	}

	parts = append(parts, typ.Name.PluralRouteName(), "%s")

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildGetSomethingArgs(proj *models.Project, typ models.DataType, includeCtx, includeSelf bool, skipIndex int) (parts []jen.Code) {
	if includeCtx {
		parts = []jen.Code{jen.ID("s").Dot("ctx")}
	} else {
		parts = []jen.Code{}
	}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == skipIndex {
			parts = append(parts, jen.EmptyString())
		} else {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	if includeSelf {
		parts = append(parts, jen.ID("s").Dotf("example%s", typ.Name.Singular()).Dot("ID"))
	}

	return parts
}

func buildTestClientSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	bodyLines := []jen.Code{
		jen.Const().ID("expectedPathFormat").Equals().Lit(buildSomethingSpecificFormatString(proj, typ)),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					append([]jen.Code{
						jen.ID("true"),
						jen.Qual("net/http", "MethodHead"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
					},
						buildGetSomethingArgs(proj, typ, false, true, -1)...,
					)...,
				),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientWithStatusCodeResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
					jen.Qual("net/http", "StatusOK"),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Qual(constants.AssertionLibrary, "True").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params().Body(
					jen.ID("t").Assign().ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
						buildGetSomethingArgs(proj, typ, true, true, i)...,
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.AssertionLibrary, "False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
					append(buildGetSomethingArgs(proj, typ, true, false, -1), jen.EmptyString())...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Qual(constants.AssertionLibrary, "False").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Qual(constants.AssertionLibrary, "False").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error executing request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Qual(constants.AssertionLibrary, "False").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_%sExists", sn).Params().Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildTestClientGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Const().ID("expectedPathFormat").Equals().Lit(buildSomethingSpecificFormatString(proj, typ)),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					append([]jen.Code{
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
					},
						buildGetSomethingArgs(proj, typ, false, true, -1)...,
					)...,
				),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientWithJSONResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
					jen.ID("s").Dotf("example%s", sn),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.MustAssertPkg, "NotNil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Qual(constants.AssertionLibrary, "Equal").Call(
					jen.ID("t"),
					jen.ID("s").Dotf("example%s", sn),
					jen.ID("actual"),
				),
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params().Body(
					jen.ID("t").Assign().ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
						buildGetSomethingArgs(proj, typ, true, true, i)...,
					),
					jen.Newline(),
					jen.Qual(constants.MustAssertPkg, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		)
	}

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Litf("with invalid %s ID", typ.Name.SingularCommonName()),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					append(buildGetSomethingArgs(proj, typ, true, false, -1), jen.EmptyString())...,
				),
				jen.Newline(),
				jen.Qual(constants.MustAssertPkg, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error executing request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					append([]jen.Code{
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
					},
						buildGetSomethingArgs(proj, typ, false, true, -1)...,
					)...,
				),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_Get%s", sn).Params().Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildListOfSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%s")
	}

	parts = append(parts, typ.Name.PluralRouteName())

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildGetListOfSomethingFormatArgs(proj *models.Project, typ models.DataType, includeCtx bool) (parts []jen.Code) {
	if includeCtx {
		parts = []jen.Code{jen.ID("s").Dot("ctx")}
	} else {
		parts = []jen.Code{}
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
	}

	return parts
}

func buildListOfSomethingArgsWithoutIndex(proj *models.Project, typ models.DataType, index int) []jen.Code {
	parts := []jen.Code{jen.ID("s").Dot("ctx")}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.EmptyString())
		} else {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	parts = append(parts, jen.ID("filter"))

	return parts
}

func buildTestClientGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()

	expectedPath := buildListOfSomethingFormatString(proj, typ)

	specArgs := append([]jen.Code{
		jen.ID("true"),
		jen.Qual("net/http", "MethodGet"),
		jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
		jen.ID("expectedPath"),
	},
		buildGetListOfSomethingFormatArgs(proj, typ, false)...,
	)

	lines := []jen.Code{
		jen.Const().ID("expectedPath").Equals().Lit(expectedPath),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Parens(jen.PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Call(jen.Nil()),
				jen.Newline(),
				jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					specArgs...,
				),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientWithJSONResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
					jen.IDf("example%sList", sn),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					buildListOfSomethingArgsWithoutIndex(proj, typ, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.MustAssertPkg, "NotNil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Qual(constants.AssertionLibrary, "Equal").Call(
					jen.ID("t"),
					jen.IDf("example%sList", sn),
					jen.ID("actual"),
				),
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params().Body(
					jen.ID("t").Assign().ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("filter").Assign().Parens(jen.PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Call(jen.Nil()),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.Underscore()).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
						buildListOfSomethingArgsWithoutIndex(proj, typ, i)...,
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		)
	}

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Parens(jen.PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Call(jen.Nil()),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					buildListOfSomethingArgsWithoutIndex(proj, typ, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error executing request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Parens(jen.PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Call(jen.Nil()),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					specArgs...,
				),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					buildListOfSomethingArgsWithoutIndex(proj, typ, -1)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_Get%s", pn).Params().Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildSearchSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%s")
	}

	parts = append(parts, typ.Name.PluralRouteName(), "search")

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildSearchSomethingFormatArgs(proj *models.Project, typ models.DataType, index int) []jen.Code {
	parts := []jen.Code{}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.EmptyString())
		} else {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	return parts
}

func buildSearchSomethingArgs(proj *models.Project, typ models.DataType, index int, withEmptyQuery bool) []jen.Code {
	parts := []jen.Code{
		jen.ID("s").Dot("ctx"),
	}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.EmptyString())
		} else {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	if withEmptyQuery {
		parts = append(parts, jen.EmptyString())
	} else {
		parts = append(parts, jen.ID("exampleQuery"))
	}

	parts = append(parts,
		jen.Zero(),
	)

	return parts
}

func buildTestClientSearchSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()

	firstSubtest := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		jen.Newline(),
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
		jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientWithJSONResponse").Call(
			jen.ID("t"),
			jen.ID("spec"),
			jen.IDf("example%sList", sn).Dot(pn),
		),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Search%s", pn).Call(
			buildSearchSomethingArgs(proj, typ, -1, false)...,
		),
		jen.Newline(),
		jen.Qual(constants.MustAssertPkg, "NotNil").Call(
			jen.ID("t"),
			jen.ID("actual"),
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(
			jen.ID("t"),
			jen.IDf("example%sList", sn).Dot(pn),
			jen.ID("actual"),
		),
	}

	lines := []jen.Code{
		jen.Const().ID("expectedPath").Equals().Lit(buildSearchSomethingFormatString(proj, typ)),
		jen.Newline(),
		jen.ID("exampleQuery").Assign().Lit("whatever"),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(jen.Lit("standard"), jen.Func().Params().Body(
			firstSubtest...,
		)),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		subtest := []jen.Code{
			jen.ID("t").Assign().ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Search%s", pn).Call(
				buildSearchSomethingArgs(proj, typ, i, false)...,
			),
			jen.Newline(),
			jen.Qual(constants.AssertionLibrary, "Nil").Call(
				jen.ID("t"),
				jen.ID("actual"),
			),
			jen.Qual(constants.AssertionLibrary, "Error").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		}

		lines = append(lines,
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()), jen.Func().Params().Body(
				subtest...,
			)),
		)
	}

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with empty query"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(
					jen.ID("t"),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Search%s", pn).Call(
					buildSearchSomethingArgs(proj, typ, -1, true)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Search%s", pn).Call(
					buildSearchSomethingArgs(proj, typ, -1, false)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with bad response from server"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
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
				jen.ID("c").Assign().ID("buildTestClientWithInvalidResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Search%s", pn).Call(
					buildSearchSomethingArgs(proj, typ, -1, false)...,
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_Search%s", pn).Params().Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildCreateSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%s")
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
		parts = []jen.Code{jen.ID("s").Dot("ctx")}
	} else {
		parts = []jen.Code{}
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, dep := range owners {
		parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))

	}

	return parts
}

func buildCreateSomethingArgsWithoutIndex(proj *models.Project, typ models.DataType, index int, includeExampleInput bool) []jen.Code {
	parts := []jen.Code{jen.ID("s").Dot("ctx")}
	owners := proj.FindOwnerTypeChain(typ)

	for i, dep := range owners {
		if i == index {
			parts = append(parts, jen.EmptyString())
		} else if i != len(owners)-1 {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	if includeExampleInput {
		parts = append(parts, jen.ID("exampleInput"))
	}

	return parts
}

func buildTestClientCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	puvn := typ.Name.PluralUnexportedVarName()

	expectedPath := buildCreateSomethingFormatString(proj, typ)
	specArgs := append([]jen.Code{
		jen.ID("false"),
		jen.Qual("net/http", "MethodPost"),
		jen.Lit(""),
		jen.ID("expectedPath"),
	},
		buildCreateSomethingFormatArgs(proj, typ, false)...,
	)

	lines := []jen.Code{
		jen.Const().ID("expectedPath").Equals().Lit(expectedPath),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sCreationRequestInput", sn).Call(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleInput").Dot("BelongsToAccount").Equals().EmptyString()),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.ID("exampleInput").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID("s").Dotf("example%sID", typ.BelongsToStruct.Singular())
					}
					return jen.Null()
				}(),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					specArgs...,
				),
				jen.List(jen.ID("c"), jen.Underscore()).Assign().ID("buildTestClientWithJSONResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
					jen.AddressOf().Qual(proj.TypesPackage(), "PreWriteResponse").Values(jen.ID("ID").MapAssign().ID("s").Dotf("example%s", sn).Dot("ID")),
				),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
					buildCreateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
				),
				jen.Qual(constants.MustAssertPkg, "NotEmpty").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Equal").Call(
					jen.ID("t"),
					jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					jen.ID("actual"),
				),
			),
		),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i, owner := range owners {
		if i != len(owners)-1 {
			subtestLines := []jen.Code{
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sCreationRequestInput", sn).Call(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleInput").Dot("BelongsToAccount").Equals().EmptyString()),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.ID("exampleInput").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().ID("s").Dotf("example%sID", typ.BelongsToStruct.Singular())
					}
					return jen.Null()
				}(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
					buildCreateSomethingArgsWithoutIndex(proj, typ, i, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.ID("err")),
			}

			lines = append(lines,
				jen.Newline(),
				jen.ID("s").Dot("Run").Call(jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()), jen.Func().Params().Body(
					subtestLines...,
				)),
			)
		}
	}

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with nil input"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
					append(
						buildCreateSomethingArgsWithoutIndex(proj, typ, -1, false),
						jen.Nil(),
					)...,
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with invalid input"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.ID("exampleInput").Assign().AddressOf().ID("types").Dotf("%sCreationRequestInput", sn).Values(),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
					buildCreateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sCreationRequestInputFrom%s", sn, sn).Call(jen.ID("s").Dotf("example%s", sn)),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
					buildCreateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error executing request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sCreationRequestInputFrom%s", sn, sn).Call(jen.ID("s").Dotf("example%s", sn)),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
					buildCreateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_Create%s", sn).Params().Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildUpdateSomethingSpecArgs(proj *models.Project, typ models.DataType) (parts []jen.Code) {
	parts = []jen.Code{}

	owners := proj.FindOwnerTypeChain(typ)
	for _, dep := range owners {
		parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
	}

	parts = append(parts, jen.ID("s").Dotf("example%s", typ.Name.Singular()).Dot("ID"))

	return parts
}

func buildUpdateSomethingArgsWithoutIndex(proj *models.Project, typ models.DataType, index int, includeSelf bool) (parts []jen.Code) {
	parts = []jen.Code{jen.ID("s").Dot("ctx")}

	owners := proj.FindOwnerTypeChain(typ)
	for i, dep := range owners {
		if i == index {
			parts = append(parts, jen.EmptyString())
		} else if i != len(owners)-1 {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	if includeSelf {
		parts = append(parts, jen.ID("s").Dotf("example%s", typ.Name.Singular()))
	}

	return parts
}

func buildTestClientUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	puvn := typ.Name.PluralUnexportedVarName()

	expectedPathFormat := buildSomethingSpecificFormatString(proj, typ)

	specArgs := append(
		[]jen.Code{
			jen.ID("false"),
			jen.Qual("net/http", "MethodPut"),
			jen.Lit(""),
			jen.ID("expectedPathFormat"),
		},
		buildUpdateSomethingSpecArgs(proj, typ)...)

	lines := []jen.Code{
		jen.Const().ID("expectedPathFormat").Equals().Lit(expectedPathFormat),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					specArgs...,
				),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientWithJSONResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
					jen.ID("s").Dotf("example%s", sn),
				),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Update%s", sn).Call(
					buildUpdateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for i, owner := range owners {
		if i != len(owners)-1 {
			lines = append(lines,
				jen.Newline(),
				jen.ID("s").Dot("Run").Call(
					jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
					jen.Func().Params().Body(
						jen.ID("t").Assign().ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
						jen.Newline(),
						jen.ID("err").Assign().ID("c").Dotf("Update%s", sn).Call(
							buildUpdateSomethingArgsWithoutIndex(proj, typ, i, true)...,
						),
						jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.ID("err")),
					),
				),
			)
		}
	}

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with nil input"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Update%s", sn).Call(
					append(buildUpdateSomethingArgsWithoutIndex(proj, typ, -1, false), jen.Nil())...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.ID("err")),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Update%s", sn).Call(
					buildUpdateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.ID("err")),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error executing request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Update%s", sn).Call(
					buildUpdateSomethingArgsWithoutIndex(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.ID("err")),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_Update%s", sn).Params().Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildTestClientArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Const().ID("expectedPathFormat").Equals().Lit(buildSomethingSpecificFormatString(proj, typ)),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("spec").Assign().ID("newRequestSpec").Call(
					append(
						[]jen.Code{
							jen.ID("true"),
							jen.Qual("net/http", "MethodDelete"),
							jen.Lit(""),
							jen.ID("expectedPathFormat"),
						},
						buildGetSomethingArgs(proj, typ, false, true, -1)...,
					)...,
				),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientWithStatusCodeResponse").Call(
					jen.ID("t"),
					jen.ID("spec"),
					jen.Qual("net/http", "StatusOK"),
				),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Archive%s", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params().Body(
					jen.ID("t").Assign().ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Assign().ID("c").Dotf("Archive%s", sn).Call(
						buildGetSomethingArgs(proj, typ, true, true, i)...,
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		)
	}

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Archive%s", sn).Call(
					append(buildGetSomethingArgs(proj, typ, true, false, -1), jen.EmptyString())...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Archive%s", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error executing request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("err").Assign().ID("c").Dotf("Archive%s", sn).Call(
					buildGetSomethingArgs(proj, typ, true, true, -1)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_Archive%s", sn).Params().Body(
			lines...,
		),
		jen.Newline(),
	}
}

func buildAuditSomethingFormatString(proj *models.Project, typ models.DataType) string {
	parts := []string{"api", "v1"}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		parts = append(parts, dep.Name.PluralRouteName(), "%s")
	}

	parts = append(parts, typ.Name.PluralRouteName(), "%s", "audit")

	return fmt.Sprintf("/%s", path.Join(parts...))
}

func buildAuditSomethingFormatArgs(proj *models.Project, typ models.DataType, index int, includeSelf bool) []jen.Code {
	parts := []jen.Code{}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.EmptyString())
		} else {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	if includeSelf {
		parts = append(parts, jen.ID("s").Dotf("example%s", typ.Name.Singular()).Dot("ID"))
	}

	return parts
}

func buildAuditSomethingArgs(proj *models.Project, typ models.DataType, index int, includeSelf bool) []jen.Code {
	parts := []jen.Code{
		jen.ID("s").Dot("ctx"),
	}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if i == index {
			parts = append(parts, jen.EmptyString())
		} else {
			parts = append(parts, jen.ID("s").Dotf("example%sID", dep.Name.Singular()))
		}
	}

	if includeSelf {
		parts = append(parts, jen.ID("s").Dotf("example%s", typ.Name.Singular()).Dot("ID"))
	} else {
		parts = append(parts, jen.EmptyString())
	}

	return parts
}

func buildTestClientGetAuditLogForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	firstSubtest := []jen.Code{
		jen.ID("t").Assign().ID("s").Dot("T").Call(),
		jen.Newline(),
		jen.ID("spec").Assign().ID("newRequestSpec").Call(
			append([]jen.Code{
				jen.ID("true"),
				jen.ID("expectedMethod"),
				jen.Lit(""),
				jen.ID("expectedPath"),
			},
				buildAuditSomethingFormatArgs(proj, typ, -1, true)...,
			)...,
		),
		jen.ID("exampleAuditLogEntryList").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryList").Call().Dot("Entries"),
		jen.Newline(),
		jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientWithJSONResponse").Call(
			jen.ID("t"),
			jen.ID("spec"),
			jen.ID("exampleAuditLogEntryList"),
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogFor%s", sn).Call(
			buildAuditSomethingArgs(proj, typ, -1, true)...,
		),
		jen.Qual(constants.MustAssertPkg, "NotNil").Call(
			jen.ID("t"),
			jen.ID("actual"),
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(
			jen.ID("t"),
			jen.ID("exampleAuditLogEntryList"),
			jen.ID("actual"),
		),
	}

	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("expectedPath").Equals().Lit(buildAuditSomethingFormatString(proj, typ)),
			jen.ID("expectedMethod").Equals().Qual("net/http", "MethodGet"),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(jen.Lit("standard"), jen.Func().Params().Body(
			firstSubtest...,
		)),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		subtest := []jen.Code{
			jen.ID("t").Assign().ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(
				jen.ID("t"),
			),
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogFor%s", sn).Call(
				buildAuditSomethingArgs(proj, typ, i, true)...,
			),
			jen.Qual(constants.AssertionLibrary, "Nil").Call(
				jen.ID("t"),
				jen.ID("actual"),
			),
			jen.Qual(constants.AssertionLibrary, "Error").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		}

		lines = append(lines,
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()), jen.Func().Params().Body(
				subtest...,
			)),
		)
	}

	lines = append(lines,
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildSimpleTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogFor%s", sn).Call(
					buildAuditSomethingArgs(proj, typ, -1, false)...,
				),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error building request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("c").Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogFor%s", sn).Call(
					buildAuditSomethingArgs(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
		jen.Newline(),
		jen.ID("s").Dot("Run").Call(
			jen.Lit("with error executing request"),
			jen.Func().Params().Body(
				jen.ID("t").Assign().ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogFor%s", sn).Call(
					buildAuditSomethingArgs(proj, typ, -1, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().IDf("%sTestSuite", puvn)).IDf("TestClient_GetAuditLogFor%s", sn).Params().Body(
			lines...,
		),
		jen.Newline(),
	}
}
