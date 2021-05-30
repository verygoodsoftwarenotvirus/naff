package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("checkItemEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("types").Dot("Item")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
				jen.Lit("expected BucketName for item #%d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Details"),
				jen.ID("actual").Dot("Details"),
				jen.Lit("expected Details for item #%d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot("Details"),
				jen.ID("actual").Dot("Details"),
			),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedOn"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Creating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be creatable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.ID("checkItemEquality").Call(
							jen.ID("t"),
							jen.ID("exampleItem"),
							jen.ID("createdItem"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("ItemCreationEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdItem").Dot("ID"),
							jen.ID("audit").Dot("ItemAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Listing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be readable in paginated form"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Var().Defs(
							jen.ID("expected").Index().Op("*").ID("types").Dot("Item"),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
							jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
							jen.List(jen.ID("createdItem"), jen.ID("itemCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItemInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdItem"),
								jen.ID("itemCreationErr"),
							),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdItem"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItems").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("True").Call(
							jen.ID("t"),
							jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Items")),
							jen.Lit("expected %d to be <= %d"),
							jen.ID("len").Call(jen.ID("expected")),
							jen.ID("len").Call(jen.ID("actual").Dot("Items")),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("createdItem")).Op(":=").Range().ID("actual").Dot("Items")).Body(
							jen.ID("assert").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
									jen.ID("ctx"),
									jen.ID("createdItem").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Searching").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be search for items"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.Var().Defs(
							jen.ID("expected").Index().Op("*").ID("types").Dot("Item"),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
							jen.ID("exampleItemInput").Dot("Name").Op("=").Qual("fmt", "Sprintf").Call(
								jen.Lit("%s %d"),
								jen.ID("exampleItemInput").Dot("Name"),
								jen.ID("i"),
							),
							jen.List(jen.ID("createdItem"), jen.ID("itemCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItemInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdItem"),
								jen.ID("itemCreationErr"),
							),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdItem"),
							),
						),
						jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(20)),
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
						jen.ID("assert").Dot("True").Call(
							jen.ID("t"),
							jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual")),
							jen.Lit("expected results length %d to be <= %d"),
							jen.ID("len").Call(jen.ID("expected")),
							jen.ID("len").Call(jen.ID("actual")),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("createdItem")).Op(":=").Range().ID("expected")).Body(
							jen.ID("assert").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
									jen.ID("ctx"),
									jen.ID("createdItem").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Searching_ReturnsOnlyItemsThatBelongToYou").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should only receive your own items"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.Var().Defs(
							jen.ID("expected").Index().Op("*").ID("types").Dot("Item"),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
							jen.ID("exampleItemInput").Dot("Name").Op("=").Qual("fmt", "Sprintf").Call(
								jen.Lit("%s %d"),
								jen.ID("exampleItemInput").Dot("Name"),
								jen.ID("i"),
							),
							jen.List(jen.ID("createdItem"), jen.ID("itemCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItemInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdItem"),
								jen.ID("itemCreationErr"),
							),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdItem"),
							),
						),
						jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(20)),
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
						jen.ID("assert").Dot("True").Call(
							jen.ID("t"),
							jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual")),
							jen.Lit("expected results length %d to be <= %d"),
							jen.ID("len").Call(jen.ID("expected")),
							jen.ID("len").Call(jen.ID("actual")),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("createdItem")).Op(":=").Range().ID("expected")).Body(
							jen.ID("assert").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
									jen.ID("ctx"),
									jen.ID("createdItem").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_ExistenceChecking_ReturnsFalseForNonexistentItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not return an error for nonexistent item"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("ItemExists").Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("False").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_ExistenceChecking_ReturnsTrueForValidItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not return an error for existent item"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("ItemExists").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("True").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Reading_Returns404ForNonexistentItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to read an item that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItem").Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
						),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Reading").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should be readable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkItemEquality").Call(
							jen.ID("t"),
							jen.ID("exampleItem"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Updating_Returns404ForNonexistentItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to update something that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.ID("exampleItem").Dot("ID").Op("=").ID("nonexistentID"),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateItem").Call(
								jen.ID("ctx"),
								jen.ID("exampleItem"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Updating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should be possible to update an item"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.ID("createdItem").Dot("Update").Call(jen.ID("converters").Dot("ConvertItemToItemUpdateInput").Call(jen.ID("exampleItem"))),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkItemEquality").Call(
							jen.ID("t"),
							jen.ID("exampleItem"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("actual").Dot("LastUpdatedOn"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("ItemCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("ItemUpdateEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdItem").Dot("ID"),
							jen.ID("audit").Dot("ItemAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Archiving_Returns404ForNonexistentItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to delete something that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("nonexistentID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Archiving").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should be possible to delete an item"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
						jen.ID("exampleItemInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
						jen.List(jen.ID("createdItem"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdItem"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveItem").Call(
								jen.ID("ctx"),
								jen.ID("createdItem").Dot("ID"),
							),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
							jen.ID("ctx"),
							jen.ID("createdItem").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("ItemCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("ItemArchiveEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdItem").Dot("ID"),
							jen.ID("audit").Dot("ItemAssignmentKey"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestItems_Auditing_Returns404ForNonexistentItem").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to audit something that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForItem").Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("Empty").Call(
							jen.ID("t"),
							jen.ID("x"),
						),
					)),
			)),
		jen.Line(),
	)

	return code
}
