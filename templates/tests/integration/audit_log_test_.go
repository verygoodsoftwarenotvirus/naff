package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAuditLogEntryListing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be read in a list by an admin"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogEntries").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("actual").Dot("Entries"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAuditLogEntryReading").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be read as an individual by an admin"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogEntries").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("actual").Dot("Entries")).Body(
							jen.List(jen.ID("y"), jen.ID("entryFetchErr")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogEntry").Call(
								jen.ID("ctx"),
								jen.ID("x").Dot("ID"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("y"),
								jen.ID("entryFetchErr"),
							),
						),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("actual").Dot("Entries"),
						),
					)),
			)),
		jen.Line(),
	)

	return code
}
