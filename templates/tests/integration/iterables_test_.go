package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()

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
				jen.Lit("expected BucketName for item #%d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
			),
			jen.Qual(constants.AssertionLibrary, "Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Details"),
				jen.ID("actual").Dot("Details"),
				jen.Lit("expected Details for item #%d to be %v, but it was %v "),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Creating").Params().Body(
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
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Comment("assert item equality"),
						jen.ID("checkItemEquality").Call(
							jen.ID("t"),
							jen.ID("exampleItem"),
							jen.ID("createdItem"),
						),
						jen.Newline(),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), "ItemCreationEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdItem").Dot("ID"),
							jen.Qual(proj.InternalAuditPackage(), "ItemAssignmentKey"),
						),
						jen.Newline(),
						jen.Comment("clean up"),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Listing").Params().Body(
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
						jen.Comment("create items"),
						jen.Var().ID("expected").Index().Op("*").Qual(proj.TypesPackage(), sn),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
							jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
							jen.Newline(),
							jen.List(jen.ID("createdItem"), jen.ID("itemCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItemInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdItem"),
								jen.ID("itemCreationErr"),
							),
							jen.Newline(),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdItem"),
							),
						),
						jen.Newline(),
						jen.Comment("assert item list equality"),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItems").Call(
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
							jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Items")),
							jen.Lit("expected %d to be <= %d"),
							jen.ID("len").Call(jen.ID("expected")),
							jen.ID("len").Call(jen.ID("actual").Dot("Items")),
						),
						jen.Newline(),
						jen.Comment("clean up"),
						jen.For(jen.List(jen.ID("_"), jen.ID("createdItem")).Op(":=").Range().ID("actual").Dot("Items")).Body(
							jen.Qual(constants.AssertionLibrary, "NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
									jen.ID("ctx"),
									jen.ID("createdItem").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Searching").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be search for items"),
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
						jen.Comment("create items"),
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.Var().ID("expected").Index().Op("*").Qual(proj.TypesPackage(), sn),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
							jen.ID("exampleItemInput").Dot("Name").Op("=").Qual("fmt", "Sprintf").Call(
								jen.Lit("%s %d"),
								jen.ID("exampleItemInput").Dot("Name"),
								jen.ID("i"),
							),
							jen.Newline(),
							jen.List(jen.ID("createdItem"), jen.ID("itemCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItemInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdItem"),
								jen.ID("itemCreationErr"),
							),
							jen.Newline(),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdItem"),
							),
						),
						jen.Newline(),
						jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(20)),
						jen.Newline(),
						jen.Comment("assert item list equality"),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("SearchItems").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem").Dot("Name"),
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
						jen.For(jen.List(jen.ID("_"), jen.ID("createdItem")).Op(":=").Range().ID("expected")).Body(
							jen.Qual(constants.AssertionLibrary, "NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
									jen.ID("ctx"),
									jen.ID("createdItem").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Searching_ReturnsOnlyItemsThatBelongToYou").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should only receive your own items"),
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
						jen.Comment("create items"),
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.Var().ID("expected").Index().Op("*").Qual(proj.TypesPackage(), sn),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
							jen.ID("exampleItemInput").Dot("Name").Op("=").Qual("fmt", "Sprintf").Call(
								jen.Lit("%s %d"),
								jen.ID("exampleItemInput").Dot("Name"),
								jen.ID("i"),
							),
							jen.Newline(),
							jen.List(jen.ID("createdItem"), jen.ID("itemCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItemInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdItem"),
								jen.ID("itemCreationErr"),
							),
							jen.Newline(),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdItem"),
							),
						),
						jen.Newline(),
						jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(20)),
						jen.Newline(),
						jen.Comment("assert item list equality"),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("SearchItems").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem").Dot("Name"),
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
						jen.For(jen.List(jen.ID("_"), jen.ID("createdItem")).Op(":=").Range().ID("expected")).Body(
							jen.Qual(constants.AssertionLibrary, "NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
									jen.ID("ctx"),
									jen.ID("createdItem").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_ExistenceChecking_ReturnsFalseForNonexistentItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not return an error for nonexistent item"),
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
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("ItemExists").Call(
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_ExistenceChecking_ReturnsTrueForValidItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not return an error for existent item"),
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
						jen.Comment("create item"),
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Comment("retrieve item"),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("ItemExists").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
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
						jen.Comment("clean up item"),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Reading_Returns404ForNonexistentItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to read an item that does not exist"),
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
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItem").Call(
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Reading").Params().Body(
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
						jen.Comment("create item"),
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Comment("retrieve item"),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Comment("assert item equality"),
						jen.ID("checkItemEquality").Call(
							jen.ID("t"),
							jen.ID("exampleItem"),
							jen.ID("actual"),
						),
						jen.Newline(),
						jen.Comment("clean up item"),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Updating_Returns404ForNonexistentItem").Params().Body(
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
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.ID("exampleItem").Dot("ID").Op("=").ID("nonexistentID"),
						jen.Newline(),
						jen.Qual(constants.AssertionLibrary, "Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItem"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("convertItemToItemUpdateInput creates an ItemUpdateInput struct from an item."),
		jen.Newline(),
		jen.Func().ID("convertItemToItemUpdateInput").Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Op("*").Qual(proj.TypesPackage(), "ItemUpdateInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("x").Dot("Name"), jen.ID("Details").Op(":").ID("x").Dot("Details"))),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Updating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should be possible to update an item"),
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
						jen.Comment("create item"),
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Comment("change item"),
						jen.ID("createdItem").Dot("Update").Call(jen.ID("convertItemToItemUpdateInput").Call(jen.ID("exampleItem"))),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem"),
							),
						),
						jen.Newline(),
						jen.Comment("retrieve changed item"),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Comment("assert item equality"),
						jen.ID("checkItemEquality").Call(
							jen.ID("t"),
							jen.ID("exampleItem"),
							jen.ID("actual"),
						),
						jen.Qual(constants.AssertionLibrary, "NotNil").Call(
							jen.ID("t"),
							jen.ID("actual").Dot("LastUpdatedOn"),
						),
						jen.Newline(),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(
							jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), "ItemCreationEvent")),
							jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), "ItemUpdateEvent")),
						),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdItem").Dot("ID"),
							jen.Qual(proj.InternalAuditPackage(), "ItemAssignmentKey"),
						),
						jen.Newline(),
						jen.Comment("clean up item"),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Archiving_Returns404ForNonexistentItem").Params().Body(
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
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("nonexistentID"),
							),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Archiving").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should be possible to delete an item"),
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
						jen.Comment("create item"),
						jen.ID("exampleItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.Comment("clean up item"),
						jen.Qual(constants.AssertionLibrary, "NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
						jen.Newline(),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Valuesln(
							jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), "ItemCreationEvent")),
							jen.Values(jen.ID("EventType").Op(":").Qual(proj.InternalAuditPackage(), "ItemArchiveEvent")),
						),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdItem").Dot("ID"),
							jen.Qual(proj.InternalAuditPackage(), "ItemAssignmentKey"),
						),
					)),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Auditing_Returns404ForNonexistentItem").Params().Body(
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
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
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
