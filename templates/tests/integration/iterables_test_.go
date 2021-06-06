package integration

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	code.Add(
		jen.Func().IDf("check%sEquality", sn).Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").Qual(proj.TypesPackage(), sn)).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Newline(),
			jen.Qual(constants.AssertionLibrary, "NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ID"),
			),
			jen.Qual(constants.AssertionLibrary, "Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
				jen.Lit("expected BucketName for "+scn+" #%d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
			),
			jen.Qual(constants.AssertionLibrary, "Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Details"),
				jen.ID("actual").Dot("Details"),
				jen.Lit("expected Details for "+scn+" #%d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot("Details"),
				jen.ID("actual").Dot("Details"),
			),
			jen.Qual(constants.AssertionLibrary, "NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedOn"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Creating", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be creatable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
						jen.List(jen.IDf("created%s", sn), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%sInput", sn),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.IDf("created%s", sn),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Commentf("assert %s equality", scn),
						jen.IDf("check%sEquality", sn).Call(
							jen.ID("t"),
							jen.IDf("example%s", sn),
							jen.IDf("created%s", sn),
						),
						jen.Newline(),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("created%s", sn).Dot("ID"),
						),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sCreationEvent", sn)))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.IDf("created%s", sn).Dot("ID"),
							jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sAssignmentKey", sn)),
						),
						jen.Newline(),
						jen.Comment("clean up"),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("created%s", sn).Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Listing", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be readable in paginated form"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Commentf("create %s", pcn),
						jen.Var().ID("expected").Index().Op("*").Qual(proj.TypesPackage(), sn),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
							jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
							jen.Newline(),
							jen.List(jen.IDf("created%s", sn), jen.IDf("%sCreationErr", uvn)).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("example%sInput", sn),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.IDf("created%s", sn),
								jen.IDf("%sCreationErr", uvn),
							),
							jen.Newline(),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.IDf("created%s", sn),
							),
						),
						jen.Newline(),
						jen.Commentf("assert %s list equality", scn),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Get%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("nil"),
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
						jen.For(jen.List(jen.ID("_"), jen.IDf("created%s", sn)).Op(":=").Range().ID("actual").Dot(pn)).Body(
							jen.Qual(constants.AssertionLibrary, "NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
									jen.ID("ctx"),
									jen.IDf("created%s", sn).Dot("ID"),
								),
							)),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Searching", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should be able to be search for %s", pcn),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Commentf("create %s", pcn),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.Var().ID("expected").Index().Op("*").Qual(proj.TypesPackage(), sn),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
							jen.IDf("example%sInput", sn).Dot("Name").Op("=").Qual("fmt", "Sprintf").Call(
								jen.Lit("%s %d"),
								jen.IDf("example%sInput", sn).Dot("Name"),
								jen.ID("i"),
							),
							jen.Newline(),
							jen.List(jen.IDf("created%s", sn), jen.IDf("%sCreationErr", uvn)).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("example%sInput", sn),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.IDf("created%s", sn),
								jen.IDf("%sCreationErr", uvn),
							),
							jen.Newline(),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.IDf("created%s", sn),
							),
						),
						jen.Newline(),
						jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(20)),
						jen.Newline(),
						jen.Commentf("assert %s list equality", scn),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Search%s", pn).Call(
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
						jen.Newline(),
						jen.Comment("clean up"),
						jen.For(jen.List(jen.ID("_"), jen.IDf("created%s", sn)).Op(":=").Range().ID("expected")).Body(
							jen.Qual(constants.AssertionLibrary, "NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
									jen.ID("ctx"),
									jen.IDf("created%s", sn).Dot("ID"),
								),
							)),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Searching_ReturnsOnly%sThatBelongToYou", pn, pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should only receive your own %s", pcn),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Commentf("create %s", pcn),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.Var().ID("expected").Index().Op("*").Qual(proj.TypesPackage(), sn),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
							jen.IDf("example%sInput", sn).Dot("Name").Op("=").Qual("fmt", "Sprintf").Call(
								jen.Lit("%s %d"),
								jen.IDf("example%sInput", sn).Dot("Name"),
								jen.ID("i"),
							),
							jen.Newline(),
							jen.List(jen.IDf("created%s", sn), jen.IDf("%sCreationErr", uvn)).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("example%sInput", sn),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.IDf("created%s", sn),
								jen.IDf("%sCreationErr", uvn),
							),
							jen.Newline(),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.IDf("created%s", sn),
							),
						),
						jen.Newline(),
						jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(20)),
						jen.Newline(),
						jen.Commentf("assert %s list equality", scn),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Search%s", pn).Call(
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
						jen.Newline(),
						jen.Comment("clean up"),
						jen.For(jen.List(jen.ID("_"), jen.IDf("created%s", sn)).Op(":=").Range().ID("expected")).Body(
							jen.Qual(constants.AssertionLibrary, "NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
									jen.ID("ctx"),
									jen.IDf("created%s", sn).Dot("ID"),
								),
							)),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_ExistenceChecking_ReturnsFalseForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should not return an error for nonexistent %s", scn),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("%sExists", sn).Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
						),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Qual(constants.AssertionLibrary, "False").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_ExistenceChecking_ReturnsTrueForValid%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("should not return an error for existent %s", scn),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Commentf("create %s", scn),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
						jen.List(jen.IDf("created%s", sn), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%sInput", sn),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.IDf("created%s", sn),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Commentf("retrieve %s", scn),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("%sExists", sn).Call(
							jen.ID("ctx"),
							jen.IDf("created%s", sn).Dot("ID"),
						),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Qual(constants.AssertionLibrary, "True").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
						jen.Newline(),
						jen.Commentf("clean up %s", scn),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("created%s", sn).Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Reading_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("it should return an error when trying to read %s that does not exist", scnwp),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Get%s", sn).Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
						),
						jen.Qual(constants.AssertionLibrary, "Error").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Reading", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should be readable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Commentf("create %s", scn),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
						jen.List(jen.IDf("created%s", sn), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%sInput", sn),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.IDf("created%s", sn),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Commentf("retrieve %s", scn),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Get%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("created%s", sn).Dot("ID"),
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
						jen.Newline(),
						jen.Commentf("clean up %s", scn),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("created%s", sn).Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Updating_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to update something that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.IDf("example%s", sn).Dot("ID").Op("=").ID("nonexistentID"),
						jen.Newline(),
						jen.Qual(constants.AssertionLibrary, "Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dotf("Update%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("example%s", sn),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("convert%sTo%sUpdateInput creates an %sUpdateInput struct from %s.", sn, sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("convert%sTo%sUpdateInput", sn, sn).Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn))).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)).Valuesln(
				jen.ID("Name").Op(":").ID("x").Dot("Name"),
				jen.ID("Details").Op(":").ID("x").Dot("Details"),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Updating", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("it should be possible to update %s", scnwp),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Commentf("create %s", scn),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
						jen.List(jen.IDf("created%s", sn), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%sInput", sn),
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
								jen.ID("ctx"),
								jen.IDf("created%s", sn),
							),
						),
						jen.Newline(),
						jen.Commentf("retrieve changed %s", scn),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Get%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("created%s", sn).Dot("ID"),
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
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("created%s", sn).Dot("ID"),
						),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(
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
						jen.Newline(),
						jen.Commentf("clean up %s", scn),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("created%s", sn).Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Archiving_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to delete something that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Qual(constants.AssertionLibrary, "Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dotf("Archive%s", sn).Call(
								jen.ID("ctx"),
								jen.ID("nonexistentID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Archiving", pn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Litf("it should be possible to delete %s", scnwp),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.Commentf("create %s", scn),
						jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
						jen.IDf("example%sInput", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
						jen.List(jen.IDf("created%s", sn), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dotf("Create%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%sInput", sn),
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
								jen.ID("ctx"),
								jen.IDf("created%s", sn).Dot("ID"),
							),
						),
						jen.Newline(),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("created%s", sn).Dot("ID"),
						),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(
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
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).IDf("Test%s_Auditing_Returns404ForNonexistent%s", pn, sn).Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to audit something that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Newline(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual(proj.InternalTracingPackage(), "StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Newline(),
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dotf("GetAuditLogFor%s", sn).Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
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
					)),
			)),
		jen.Newline(),
	)

	return code
}
